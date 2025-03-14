package core

import (
	"sync"
	"sync/atomic"
)

type BareJob struct {
	Schedule string `hash:"true"`
	Name     string `hash:"true"`
	Command  string `hash:"true"`

	middlewareContainer
	running int32
	lock    sync.Mutex
	history []*Execution
	cronID  int
}

func (j *BareJob) GetName() string {
	return j.Name
}

func (j *BareJob) GetSchedule() string {
	return j.Schedule
}

func (j *BareJob) GetCommand() string {
	return j.Command
}

// GetProcessedCommand returns the command with variables replaced
func (j *BareJob) GetProcessedCommand(context VariableContext) string {
	processed, err := ProcessVariables(j.Command, context)
	if err != nil {
		// If there's an error processing variables, return the original command
		return j.Command
	}
	return processed
}

func (j *BareJob) Running() int32 {
	return atomic.LoadInt32(&j.running)
}

func (j *BareJob) NotifyStart() {
	atomic.AddInt32(&j.running, 1)
}

func (j *BareJob) NotifyStop() {
	atomic.AddInt32(&j.running, -1)
}

func (j *BareJob) GetCronJobID() int {
	return j.cronID
}

func (j *BareJob) SetCronJobID(id int) {
	j.cronID = id
}

func (j *BareJob) SetVolumeMounts(volumes []string) {
	// BareJob doesn't use volumes, but we need to implement the interface
}
