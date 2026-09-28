package core

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"
)

type SuiteScheduler struct{}

var _ = Suite(&SuiteScheduler{})

func (s *SuiteScheduler) TestAddJob(c *C) {
	job := &TestJob{}
	job.Schedule = "@hourly"

	logger := &schedulerTestLogger{}
	sc := NewScheduler(logger)
	err := sc.AddJob(job)
	c.Assert(err, IsNil)
	c.Assert(logger.errors, HasLen, 0)

	e := sc.cron.Entries()
	c.Assert(e, HasLen, 1)
	c.Assert(e[0].Job.(*jobWrapper).j, DeepEquals, job)
}

func (s *SuiteScheduler) TestAddJobInvalidSchedule(c *C) {
	for _, schedule := range []string{"", "@every 1d", "0 0 4 * * *"} {
		logger := &schedulerTestLogger{}
		sc := NewScheduler(logger)
		job := &TestJob{}
		job.Name = "optimize"
		job.Schedule = schedule

		err := sc.AddJob(job)
		c.Assert(err, NotNil)
		if schedule == "" {
			c.Check(err, Equals, ErrEmptySchedule)
		}
		c.Check(logger.errors, DeepEquals, []string{
			fmt.Sprintf("Failed to register job %q with schedule %q: %v", "optimize", schedule, err),
		})
		c.Check(sc.cron.Entries(), HasLen, 0)
		c.Check(job.GetCronJobID(), Equals, 0)
	}
}

type schedulerTestLogger struct {
	TestLogger
	errors []string
}

func (l *schedulerTestLogger) Errorf(format string, args ...interface{}) {
	l.errors = append(l.errors, fmt.Sprintf(format, args...))
}

func (s *SuiteScheduler) TestStartStop(c *C) {
	job := &TestJob{}
	job.Schedule = "@every 1s"

	sc := NewScheduler(&TestLogger{})
	err := sc.AddJob(job)
	c.Assert(err, IsNil)

	sc.Start()
	c.Assert(sc.IsRunning(), Equals, true)

	time.Sleep(time.Second * 2)

	sc.Stop()
	c.Assert(sc.IsRunning(), Equals, false)
}

func (s *SuiteScheduler) TestMergeMiddlewaresSame(c *C) {
	mA, mB, mC := &TestMiddleware{}, &TestMiddleware{}, &TestMiddleware{}

	job := &TestJob{}
	job.Schedule = "@every 1s"
	job.Use(mB, mC)

	sc := NewScheduler(&TestLogger{})
	sc.Use(mA)
	sc.AddJob(job)

	m := job.Middlewares()
	c.Assert(m, HasLen, 1)
	c.Assert(m[0], Equals, mB)
}
