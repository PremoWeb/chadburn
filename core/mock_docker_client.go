package core

import (
	"context"
	"io"
)

// MockDockerClient is a mock implementation of the DockerClient interface for testing
type MockDockerClient struct {
	// Add fields to track method calls and return values
}

// ListContainers lists containers with the given filters
func (c *MockDockerClient) ListContainers(filterMap map[string][]string) ([]Container, error) {
	return []Container{}, nil
}

// InspectContainer inspects a container by ID
func (c *MockDockerClient) InspectContainer(id string) (*Container, error) {
	container := &Container{
		ID:     id,
		Name:   id,
		Labels: map[string]string{},
	}
	container.State.Running = true
	container.Config.Labels = map[string]string{}
	return container, nil
}

// CreateContainer creates a new container
func (c *MockDockerClient) CreateContainer(config *ContainerConfig) (*Container, error) {
	return &Container{}, nil
}

// StartContainer starts a container
func (c *MockDockerClient) StartContainer(id string) error {
	return nil
}

// StopContainer stops a container
func (c *MockDockerClient) StopContainer(id string) error {
	return nil
}

// RemoveContainer removes a container
func (c *MockDockerClient) RemoveContainer(id string) error {
	return nil
}

// WaitContainer waits for a container to exit and returns its exit code
func (c *MockDockerClient) WaitContainer(id string) (int, error) {
	return 0, nil
}

// CreateExec creates an exec instance in a container
func (c *MockDockerClient) CreateExec(containerID string, cmd []string, config *ExecConfig) (string, error) {
	return "", nil
}

// StartExec starts an exec instance
func (c *MockDockerClient) StartExec(execID string, attachStdout, attachStderr bool) (io.ReadCloser, error) {
	// Create a pipe to capture the output
	r, w := io.Pipe()

	// Write some data to the pipe
	go func() {
		w.Write([]byte("test output"))
		w.Close()
	}()

	return r, nil
}

// InspectExec inspects an exec instance
func (c *MockDockerClient) InspectExec(execID string) (*ExecInspect, error) {
	return &ExecInspect{}, nil
}

// PullImage pulls an image from a registry
func (c *MockDockerClient) PullImage(imageName string) error {
	return nil
}

// WatchEvents watches Docker events and sends them to the provided channel
func (c *MockDockerClient) WatchEvents(ctx context.Context, eventCh chan<- *DockerEvent, errCh chan<- error) {
	// No-op for mock
}

// CreateService creates a new service
func (c *MockDockerClient) CreateService(config *ServiceConfig) (string, error) {
	return "", nil
}

// InspectService inspects a service
func (c *MockDockerClient) InspectService(id string) (*Service, error) {
	return &Service{}, nil
}

// ListTasks lists tasks for a service
func (c *MockDockerClient) ListTasks(serviceID string) ([]Task, error) {
	return []Task{}, nil
}

// RemoveService removes a service
func (c *MockDockerClient) RemoveService(id string) error {
	return nil
}

// CreateExecJob creates a new ExecJob
func (c *MockDockerClient) CreateExecJob() *OfficialExecJob {
	return NewOfficialExecJob(c)
}

// CreateRunJob creates a new RunJob
func (c *MockDockerClient) CreateRunJob() *OfficialRunJob {
	return NewOfficialRunJob(c)
}

// CreateServiceJob creates a new ServiceJob
func (c *MockDockerClient) CreateServiceJob() *OfficialServiceJob {
	return NewOfficialServiceJob(c)
}

// CreateLifecycleJob creates a new LifecycleJob
func (c *MockDockerClient) CreateLifecycleJob() *OfficialLifecycleJob {
	return NewOfficialLifecycleJob(c)
}

// Close closes the Docker client
func (c *MockDockerClient) Close() error {
	return nil
}
