package core

import (
	"crypto/rand"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/armon/circbuf"
)

var (
	// ErrSkippedExecution pass this error to `Execution.Stop` if you wish to mark
	// it as skipped.
	ErrSkippedExecution   = errors.New("skipped execution")
	ErrUnexpected         = errors.New("error unexpected, docker has returned exit code -1, maybe wrong user?")
	ErrMaxTimeRunning     = errors.New("the job has exceed the maximum allowed time running.")
	ErrLocalImageNotFound = errors.New("couldn't find image on the host")
)

// maximum size of a stdout/stderr stream to be kept in memory and optional stored/sent via mail
const maxStreamSize = 10 * 1024 * 1024

type Job interface {
	GetName() string
	GetSchedule() string
	GetCommand() string
	GetProcessedCommand(VariableContext) string
	Middlewares() []Middleware
	Use(...Middleware)
	Run(*Context) error
	Running() int32
	NotifyStart()
	NotifyStop()
	GetCronJobID() int
	SetCronJobID(int)
	SetVolumeMounts([]string)
}

type Context struct {
	Scheduler   *Scheduler
	Logger      Logger
	Job         Job
	Execution   *Execution
	middlewares []Middleware
	current     int
	executed    bool
}

// NewContext creates a new Context
func NewContext(s *Scheduler, j Job, e *Execution) *Context {
	return &Context{
		Scheduler:   s,
		Job:         j,
		Logger:      s.Logger,
		Execution:   e,
		middlewares: j.Middlewares(),
	}
}

func (c *Context) Run() error {
	c.Job.NotifyStart()
	c.Execution.Start()

	for {
		m, end := c.getNext()
		if end {
			break
		}

		// Check if execution is still running
		if c.Execution.Error() != nil && !m.ContinueOnStop() {
			continue
		}

		return m.Run(c)
	}

	// Check if execution is still running
	if c.Execution.Error() != nil {
		return c.Execution.Error()
	}

	c.executed = true
	return c.Job.Run(c)
}

func (c *Context) getNext() (Middleware, bool) {
	if c.current >= len(c.middlewares) {
		return nil, true
	}

	c.current++
	return c.middlewares[c.current-1], false
}

func (c *Context) Stop(err error) {
	if c.Execution.Error() != nil {
		return
	}

	c.Execution.Stop(err)
	c.Job.NotifyStop()
}

func (c *Context) Log(msg string) {
	format := "[Job %q (%s)] %s"
	args := []interface{}{c.Job.GetName(), c.Execution.ID, msg}

	switch {
	case c.Execution.Failed():
		c.Logger.Errorf(format, args...)
	case c.Execution.Skipped():
		c.Logger.Warningf(format, args...)
	default:
		c.Logger.Noticef(format, args...)
	}
}

// Execution contains all the information relative to a Job execution.
type Execution struct {
	ID           string
	Date         time.Time
	OutputStream *circbuf.Buffer
	ErrorStream  *circbuf.Buffer
	mutex        sync.Mutex
	current      int
	start        time.Time
	err          error
	end          time.Time
}

func NewExecution() *Execution {
	stdout, _ := circbuf.NewBuffer(maxStreamSize)
	stderr, _ := circbuf.NewBuffer(maxStreamSize)

	return &Execution{
		ID:           randomID(),
		OutputStream: stdout,
		ErrorStream:  stderr,
		Date:         time.Now(),
	}
}

func (e *Execution) Start() {
	e.start = time.Now()
}

func (e *Execution) Stop(err error) {
	e.err = err
	e.end = time.Now()
}

func (e *Execution) Error() error {
	return e.err
}

func (e *Execution) Failed() bool {
	if e.err != nil && e.err != ErrSkippedExecution {
		return true
	}

	return false
}

func (e *Execution) Skipped() bool {
	return e.err == ErrSkippedExecution
}

func (e *Execution) Duration() time.Duration {
	if e.start.IsZero() {
		return 0
	}

	end := time.Now()
	return end.Sub(e.start)
}

// Middleware can wrap any job execution, allowing to execution code before
// or/and after of each `Job.Run`
type Middleware interface {
	// Run is called instead of the original `Job.Run`, you MUST call to `ctx.Run`
	// inside of the middleware `Run` function otherwise you will broken the
	// Job workflow.
	Run(*Context) error
	// ContinueOnStop,  If return true the Run function will be called even if
	// the execution is stopped
	ContinueOnStop() bool
}

type middlewareContainer struct {
	m     map[string]Middleware
	order []string
}

func (c *middlewareContainer) Use(ms ...Middleware) {
	if c.m == nil {
		c.m = make(map[string]Middleware, 0)
	}

	for _, m := range ms {
		if m == nil {
			continue
		}

		t := reflect.TypeOf(m).String()
		if _, ok := c.m[t]; ok {
			continue
		}

		c.order = append(c.order, t)
		c.m[t] = m
	}
}

func (c *middlewareContainer) Middlewares() []Middleware {
	var ms []Middleware
	for _, t := range c.order {
		ms = append(ms, c.m[t])
	}

	return ms
}

type Logger interface {
	Criticalf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Noticef(format string, args ...interface{})
	Warningf(format string, args ...interface{})
}

func randomID() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", b)
}

// Helper function to parse a registry from a repository
func parseRegistry(repository string) string {
	parts := strings.Split(repository, "/")
	if len(parts) < 2 {
		return ""
	}

	if strings.ContainsAny(parts[0], ".:") || len(parts) > 2 {
		return parts[0]
	}

	return ""
}

const HashmeTagName = "hash"

func getHash(t reflect.Type, v reflect.Value, hash *string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldv := v.Field(i)
		kind := field.Type.Kind()

		if kind == reflect.Struct {
			getHash(field.Type, fieldv, hash)
			continue
		}

		hashmeTag := field.Tag.Get(HashmeTagName)
		if hashmeTag == "true" {
			if kind == reflect.String {
				*hash += fieldv.String()
			} else if kind == reflect.Int32 || kind == reflect.Int || kind == reflect.Int64 || kind == reflect.Int16 || kind == reflect.Int8 {
				*hash += strconv.FormatInt(fieldv.Int(), 10)
			} else if kind == reflect.Bool {
				*hash += strconv.FormatBool(fieldv.Bool())
			} else {
				panic("Unsupported field type")
			}
		}
	}
}

// Start starts the execution
func (c *Context) Start() {
	c.Execution.Start()
}

// Next executes the next middleware in the chain
func (c *Context) Next() error {
	return c.Run()
}

// IsRunning returns true if the execution is running
func (e *Execution) IsRunning() bool {
	return e.start.IsZero() == false && e.end.IsZero() == true && e.err == nil
}
