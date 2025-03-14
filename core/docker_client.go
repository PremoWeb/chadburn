package core

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// DockerClient is an interface for Docker client operations
type DockerClient interface {
	ListContainers(filters map[string][]string) ([]Container, error)
	InspectContainer(id string) (*Container, error)
	WatchEvents(ctx context.Context, eventCh chan<- *DockerEvent, errCh chan<- error)
	Close() error
}

// Container represents a Docker container
type Container struct {
	ID     string
	Name   string
	Labels map[string]string
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

// OfficialDockerClient implements DockerClient using the official Docker client
type OfficialDockerClient struct {
	client *client.Client
	ctx    context.Context
}

// NewDockerClient creates a new Docker client
func NewDockerClient() (DockerClient, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &OfficialDockerClient{
		client: cli,
		ctx:    ctx,
	}, nil
}

// ListContainers lists containers with the given filters
func (c *OfficialDockerClient) ListContainers(filterMap map[string][]string) ([]Container, error) {
	// Convert filters to Docker filter format
	filterArgs := filters.NewArgs()
	for k, values := range filterMap {
		for _, v := range values {
			filterArgs.Add(k, v)
		}
	}

	// List containers
	containers, err := c.client.ContainerList(c.ctx, container.ListOptions{
		Filters: filterArgs,
	})
	if err != nil {
		return nil, err
	}

	// Convert to our Container type
	result := make([]Container, len(containers))
	for i, container := range containers {
		name := ""
		if len(container.Names) > 0 {
			name = strings.TrimPrefix(container.Names[0], "/")
		}
		result[i] = Container{
			ID:     container.ID,
			Name:   name,
			Labels: container.Labels,
		}
	}

	return result, nil
}

// InspectContainer inspects a container by ID
func (c *OfficialDockerClient) InspectContainer(id string) (*Container, error) {
	containerInfo, err := c.client.ContainerInspect(c.ctx, id)
	if err != nil {
		return nil, err
	}

	return &Container{
		ID:     containerInfo.ID,
		Name:   strings.TrimPrefix(containerInfo.Name, "/"),
		Labels: containerInfo.Config.Labels,
	}, nil
}

// WatchEvents watches Docker events and sends them to the provided channel
func (c *OfficialDockerClient) WatchEvents(ctx context.Context, eventCh chan<- *DockerEvent, errCh chan<- error) {
	// Create a context that can be cancelled
	if ctx == nil {
		ctx = c.ctx
	}

	// Watch Docker events with empty filters
	messages, errs := c.client.Events(ctx, events.ListOptions{})

	// Process events
	go func() {
		for {
			select {
			case err := <-errs:
				if err != nil && err != io.EOF {
					errCh <- err
				}
				if err != nil && err == io.EOF {
					errCh <- io.EOF
					return
				}
			case message := <-messages:
				// Only process container events
				if message.Type == events.ContainerEventType {
					event := &DockerEvent{
						Action: string(message.Action),
						ID:     message.ID,
						Type:   string(message.Type),
						Actor: DockerActor{
							ID:         message.Actor.ID,
							Attributes: message.Actor.Attributes,
						},
						Attributes: message.Actor.Attributes,
					}
					eventCh <- event
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Close closes the Docker client
func (c *OfficialDockerClient) Close() error {
	return c.client.Close()
}
