package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/robfig/cron/v3"
)

var (
	ErrEmptyScheduler = errors.New("unable to start a empty scheduler.")
	ErrEmptySchedule  = errors.New("unable to add a job with a empty schedule.")
)

var (
	SchedulerJobs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "chadburn_scheduler_jobs",
		Help: "Active job count registered on the scheduler.",
	})
	JobRegisterErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chadburn_scheduler_register_errors_total",
		Help: "Total number of failed scheduler registrations.",
	})
	RunsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chadburn_run_total",
			Help: "Total number of completed job runs.",
		},
		[]string{"job_name"},
	)
	RunErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chadburn_run_errors_total",
			Help: "Total number of completed job runs that resulted in an error.",
		},
		[]string{"job_name"},
	)
	RunLatest = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "chadburn_run_latest_timestamp",
			Help: "Last time a job run completed.",
		},
		[]string{"job_name"},
	)
	RunDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "chadburn_run_duration_seconds",
			Help: "Duration of all runs.",
		},
		[]string{"job_name"},
	)
)

type Scheduler struct {
	Jobs   []Job
	Logger Logger

	middlewareContainer
	cron      *cron.Cron
	wg        sync.WaitGroup
	isRunning bool
}

func NewScheduler(l Logger) *Scheduler {
	cronUtils := NewCronUtils(l)
	return &Scheduler{
		Logger: l,
		cron:   cron.New(cron.WithLogger(cronUtils), cron.WithChain(cron.Recover(cronUtils))),
	}
}

func (s *Scheduler) AddJob(j Job) error {
	if j.GetSchedule() == "" {
		JobRegisterErrorsTotal.Inc()
		return ErrEmptySchedule
	}

	id, err := s.cron.AddJob(j.GetSchedule(), &jobWrapper{s, j})
	if err != nil {
		JobRegisterErrorsTotal.Inc()
		return err
	}
	j.SetCronJobID(int(id)) // Cast to int in order to avoid pushing cron external to common
	j.Use(s.Middlewares()...)
	SchedulerJobs.Inc()
	s.Logger.Noticef("New job registered %q - %q - %q - ID: %v", j.GetName(), j.GetCommand(), j.GetSchedule(), id)
	return nil
}

func (s *Scheduler) RemoveJob(j Job) error {
	s.Logger.Noticef("Job deregistered (will not fire again) %q - %q - %q - ID: %v", j.GetName(), j.GetCommand(), j.GetSchedule(), j.GetCronJobID())
	s.cron.Remove(cron.EntryID(j.GetCronJobID()))
	SchedulerJobs.Dec()
	return nil
}

func (s *Scheduler) Start() error {
	s.Logger.Debugf("Starting scheduler")
	s.isRunning = true
	s.cron.Start()
	return nil
}

func (s *Scheduler) Stop() error {
	s.wg.Wait()
	s.cron.Stop()
	s.isRunning = false

	return nil
}

func (s *Scheduler) IsRunning() bool {
	return s.isRunning
}

type jobWrapper struct {
	s *Scheduler
	j Job
}

func (w *jobWrapper) Run() {
	ctx := &Context{
		Scheduler:   w.s,
		Job:         w.j,
		Logger:      w.s.Logger,
		Execution:   NewExecution(),
		middlewares: w.s.Middlewares(),
	}

	w.start(ctx)
	err := ctx.Run()
	w.stop(ctx, err)
}

func (w *jobWrapper) start(ctx *Context) {
	ctx.Log("Started - " + ctx.Job.GetCommand())
}

func (w *jobWrapper) stop(ctx *Context, err error) {
	ctx.Stop(err)

	errText := "none"
	if ctx.Execution.Error() != nil {
		errText = ctx.Execution.Error().Error()
	}

	output := ctx.Execution.OutputStream.Bytes()

	if len(output) > 0 {
		ctx.Log("Output: " + string(output))
	}

	msg := fmt.Sprintf(
		"Finished in %q, failed: %t, skipped: %t, error: %s",
		ctx.Execution.Duration(), ctx.Execution.Failed(), ctx.Execution.Skipped(), errText,
	)

	RunsTotal.WithLabelValues(ctx.Job.GetName()).Inc()
	if ctx.Execution.Failed() {
		RunErrorsTotal.WithLabelValues(ctx.Job.GetName()).Inc()
	}
	RunLatest.WithLabelValues(ctx.Job.GetName()).SetToCurrentTime()
	RunDuration.WithLabelValues(ctx.Job.GetName()).Observe(ctx.Execution.Duration().Seconds())

	ctx.Log(msg)
}
