package core

import (
	"context"
	"io"
	"time"
)

// DockerClient is an interface for Docker client operations
type DockerClient interface {
	// Container operations
	ListContainers(filters map[string][]string) ([]Container, error)
	InspectContainer(id string) (*Container, error)
	CreateContainer(config *ContainerConfig) (*Container, error)
	StartContainer(id string) error
	StopContainer(id string) error
	RemoveContainer(id string) error
	WaitContainer(id string) (int, error)

	// Exec operations
	CreateExec(containerID string, cmd []string, config *ExecConfig) (string, error)
	StartExec(execID string, attachStdout, attachStderr bool) (io.ReadCloser, error)
	InspectExec(execID string) (*ExecInspect, error)

	// Image operations
	PullImage(image string) error

	// Service operations
	CreateService(config *ServiceConfig) (string, error)
	InspectService(id string) (*Service, error)
	ListTasks(serviceID string) ([]Task, error)
	RemoveService(id string) error

	// Event operations
	WatchEvents(ctx context.Context, eventCh chan<- *DockerEvent, errCh chan<- error)

	// Job creation operations
	CreateExecJob() *OfficialExecJob
	CreateRunJob() *OfficialRunJob
	CreateServiceJob() *OfficialServiceJob
	CreateLifecycleJob() *OfficialLifecycleJob

	// Cleanup
	Close() error
}

// Container represents a Docker container
type Container struct {
	ID     string
	Name   string
	Labels map[string]string
	State  struct {
		Running bool
	}
	Config struct {
		Labels map[string]string
	}
}

// ContainerConfig represents configuration for creating a container
type ContainerConfig struct {
	Image        string
	Cmd          []string
	Env          []string
	WorkingDir   string
	User         string
	Volumes      map[string]struct{}
	ExposedPorts map[string]struct{}
	Labels       map[string]string
	Tty          bool
	AttachStdin  bool
	AttachStdout bool
	AttachStderr bool
	NetworkMode  string
	HostConfig   *HostConfig
}

// HostConfig represents the container's host configuration
type HostConfig struct {
	Binds       []string
	NetworkMode string
}

// ExecConfig represents configuration for creating an exec instance
type ExecConfig struct {
	User         string
	Tty          bool
	AttachStdin  bool
	AttachStdout bool
	AttachStderr bool
	Cmd          []string
	WorkingDir   string
}

// ExecInspect represents the result of inspecting an exec instance
type ExecInspect struct {
	ID       string
	Running  bool
	ExitCode int
}

// ServiceConfig represents configuration for creating a service
type ServiceConfig struct {
	Name          string
	Image         string
	Cmd           []string
	Env           []string
	Labels        map[string]string
	RestartPolicy *RestartPolicy
	Networks      []string
	Mounts        []ServiceMount
}

// RestartPolicy represents the restart policy for a service
type RestartPolicy struct {
	Condition   string
	Delay       time.Duration
	MaxAttempts int
}

// ServiceMount represents a mount for a service
type ServiceMount struct {
	Source string
	Target string
	Type   string
}

// Service represents a Docker service
type Service struct {
	ID          string
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Spec        ServiceConfig
	Endpoint    ServiceEndpoint
	UpdateState ServiceUpdateState
}

// ServiceEndpoint represents the endpoint of a service
type ServiceEndpoint struct {
	Ports []ServicePort
}

// ServicePort represents a port of a service
type ServicePort struct {
	Protocol      string
	TargetPort    uint32
	PublishedPort uint32
}

// ServiceUpdateState represents the update state of a service
type ServiceUpdateState struct {
	State     string
	Message   string
	StartedAt time.Time
}

// Task represents a Docker task
type Task struct {
	ID           string
	ServiceID    string
	NodeID       string
	Status       TaskStatus
	DesiredState string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TaskStatus represents the status of a task
type TaskStatus struct {
	Timestamp time.Time
	State     string
	Message   string
	Err       string
}

// DockerEvent represents a Docker event
type DockerEvent struct {
	Action     string
	ID         string
	Type       string
	Actor      DockerActor
	Attributes map[string]string
}

// DockerActor represents the actor in a Docker event
type DockerActor struct {
	ID         string
	Attributes map[string]string
}
