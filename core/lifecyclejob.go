package core

import (
	"reflect"

	docker "github.com/fsouza/go-dockerclient"
)

// LifecycleEventType represents the type of container lifecycle event
type LifecycleEventType string

const (
	// ContainerStart represents a container start event
	ContainerStart LifecycleEventType = "start"
	// ContainerStop represents a container stop event
	ContainerStop LifecycleEventType = "stop"
)

// LifecycleJob represents a job that runs once on container lifecycle events
type LifecycleJob struct {
	BareJob   `mapstructure:",squash"`
	Client    *docker.Client     `json:"-"`
	Container string             `hash:"true"` // Container ID or name to monitor
	EventType LifecycleEventType `hash:"true"` // Type of event to trigger on (start, stop)
	Executed  bool               `hash:"-"`    // Whether this job has been executed
}

// NewLifecycleJob creates a new LifecycleJob
func NewLifecycleJob(c *docker.Client) *LifecycleJob {
	return &LifecycleJob{
		Client:    c,
		EventType: ContainerStart, // Default to start event
		Executed:  false,
	}
}

// Run executes the job
func (j *LifecycleJob) Run(ctx *Context) error {
	if j.Executed {
		// Skip if already executed
		return nil
	}

	// Create variable context with container info
	varContext := VariableContext{
		Container: ContainerInfo{
			Name: j.Container,
			ID:   j.Container, // We use the container name as ID for now
		},
	}

	// Get processed command with variables replaced
	processedCommand := j.GetProcessedCommand(varContext)

	// Execute the command locally
	localJob := &LocalJob{}
	// Set fields individually instead of copying the BareJob struct
	localJob.Name = j.Name
	localJob.Schedule = j.Schedule
	localJob.Command = processedCommand // Use the processed command
	localJob.ContainerID = j.Container
	localJob.ContainerName = j.Container

	// Run the local job
	err := localJob.Run(ctx)
	if err != nil {
		return err
	}

	// Mark as executed
	j.Executed = true
	return nil
}

// Returns a hash of all the job attributes. Used to detect changes
func (j *LifecycleJob) Hash() string {
	var hash string
	getHash(reflect.TypeOf(j).Elem(), reflect.ValueOf(j).Elem(), &hash)
	return hash
}

// ShouldRun determines if the job should run based on the event type
func (j *LifecycleJob) ShouldRun(eventType LifecycleEventType) bool {
	return !j.Executed && j.EventType == eventType
}

// Reset resets the executed state of the job
func (j *LifecycleJob) Reset() {
	j.Executed = false
}

// SetVolumeMounts is a no-op for LifecycleJob
func (j *LifecycleJob) SetVolumeMounts(volumes []string) {
	// No-op for LifecycleJob
}
