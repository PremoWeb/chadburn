package cli

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/PremoWeb/Chadburn/core"
)

var ErrNoContainerWithChadburnEnabled = errors.New("couldn't find containers with label 'chadburn.enabled=true'")

type DockerHandler struct {
	dockerClient  core.DockerClient
	notifier      dockerLabelsUpdate
	logger        core.Logger
	lifecycleJobs map[string]*LifecycleJobConfig // Map of lifecycle jobs
	ctx           context.Context
	cancel        context.CancelFunc
}

type dockerLabelsUpdate interface {
	dockerLabelsUpdate(map[string]map[string]string)
}

// GetLifecycleJobs returns the lifecycle jobs
func (c *DockerHandler) GetLifecycleJobs() map[string]*LifecycleJobConfig {
	return c.lifecycleJobs
}

// SetLifecycleJobs sets the lifecycle jobs
func (c *DockerHandler) SetLifecycleJobs(jobs map[string]*LifecycleJobConfig) {
	c.lifecycleJobs = jobs
}

// GetInternalDockerClient returns the internal Docker client
func (c *DockerHandler) GetInternalDockerClient() core.DockerClient {
	return c.dockerClient
}

func NewDockerHandler(notifier dockerLabelsUpdate, logger core.Logger) (*DockerHandler, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a new official Docker client
	client, err := core.NewDockerClient()
	if err != nil {
		cancel() // Cancel the context if there's an error
		return nil, err
	}

	c := &DockerHandler{
		dockerClient:  client,
		notifier:      notifier,
		logger:        logger,
		lifecycleJobs: make(map[string]*LifecycleJobConfig),
		ctx:           ctx,
		cancel:        cancel,
	}

	go c.watch()
	go c.watchEvents() // Start watching Docker events
	return c, nil
}

// watch polls for changes in Docker containers
func (c *DockerHandler) watch() {
	// Poll for changes
	tick := time.Tick(10000 * time.Millisecond)

	// Add backoff for connection issues
	connectionBackoff := 100 * time.Millisecond
	maxConnectionBackoff := 5 * time.Second

	for {
		select {
		case <-tick:
			labels, err := c.GetDockerLabels()
			// Do not print or care if there is no container up right now
			if err != nil && !errors.Is(err, ErrNoContainerWithChadburnEnabled) {
				c.logger.Debugf("%v", err)

				// Check if it's a connection error
				if strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
					// Apply exponential backoff
					if connectionBackoff < maxConnectionBackoff {
						connectionBackoff *= 2
					}
					c.logger.Noticef("Docker daemon connection issue. Waiting %v before next attempt...", connectionBackoff)
					time.Sleep(connectionBackoff)
					continue
				}
			} else {
				// Reset backoff on success
				connectionBackoff = 100 * time.Millisecond
			}

			c.notifier.dockerLabelsUpdate(labels)
		case <-c.ctx.Done():
			return
		}
	}
}

// watchEvents listens for Docker events and triggers lifecycle jobs
func (c *DockerHandler) watchEvents() {
	// Add a backoff mechanism to prevent rapid restarts
	backoff := 100 * time.Millisecond
	maxBackoff := 5 * time.Second

	for {
		// Create channels to receive Docker events and errors
		eventCh := make(chan *core.DockerEvent)
		errCh := make(chan error)

		// Start watching events
		c.dockerClient.WatchEvents(c.ctx, eventCh, errCh)

		c.logger.Noticef("Started watching Docker events")

		// Check if context is done
		select {
		case <-c.ctx.Done():
			return
		default:
			// Continue processing events
		}

		// Listen for events
		shouldReconnect := false
		for !shouldReconnect {
			select {
			case event := <-eventCh:
				// Successfully received an event, reset backoff
				backoff = 100 * time.Millisecond

				// Only process container events
				if event.Type != "container" {
					continue
				}

				// Get container info
				container, err := c.dockerClient.InspectContainer(event.ID)
				if err != nil {
					c.logger.Debugf("Failed to inspect container %s: %v", event.ID, err)
					continue
				}

				// Process the event
				switch event.Action {
				case "start":
					c.logger.Debugf("Container %s started", container.Name)
					c.processLifecycleEvent(container.Name, event.ID, core.ContainerStart)
				case "die", "stop":
					c.logger.Debugf("Container %s stopped", container.Name)
					c.processLifecycleEvent(container.Name, event.ID, core.ContainerStop)
				}
			case err := <-errCh:
				c.logger.Errorf("Error watching events: %v", err)

				// Handle all connection errors with backoff
				if err == io.EOF || strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
					// Handle connection errors by reconnecting after a delay
					shouldReconnect = true
					// Apply exponential backoff with a maximum limit
					if backoff < maxBackoff {
						backoff *= 2
					}
					c.logger.Noticef("Docker events connection issue. Reconnecting in %v...", backoff)
					time.Sleep(backoff)
				}
			case <-c.ctx.Done():
				return
			}
		}
	}
}

// processLifecycleEvent processes a container lifecycle event
func (c *DockerHandler) processLifecycleEvent(containerName, containerID string, eventType core.LifecycleEventType) {
	// Check if we have any lifecycle jobs for this container
	for name, job := range c.lifecycleJobs {
		// Check if this job should run for this container and event type
		if job.Container == containerName && job.EventType == eventType && !job.Executed {
			c.logger.Noticef("Running lifecycle job %s for container %s on %s event", name, containerName, eventType)

			// Create execution context
			ctx := &core.Context{
				Execution: core.NewExecution(),
				Logger:    c.logger,
				Job:       &job.LifecycleJob,
			}

			// Run the job
			err := job.Run(ctx)
			if err != nil {
				c.logger.Errorf("Failed to run lifecycle job %s: %v", name, err)
			} else {
				c.logger.Noticef("Lifecycle job %s completed successfully", name)
			}
		}
	}
}

func (c *DockerHandler) GetDockerLabels() (map[string]map[string]string, error) {
	// First, get containers with the required label
	conts, err := c.dockerClient.ListContainers(map[string][]string{
		"label": {requiredLabelFilter},
	})
	if err != nil {
		return nil, err
	}

	// Also get containers with job-run labels
	jobRunConts, err := c.dockerClient.ListContainers(map[string][]string{
		"label": {labelPrefix + "." + jobRun},
	})
	if err != nil {
		return nil, err
	}

	// Combine the two lists, avoiding duplicates
	contMap := make(map[string]core.Container)
	for _, cont := range conts {
		contMap[cont.Name] = cont
	}
	for _, cont := range jobRunConts {
		contMap[cont.Name] = cont
	}

	if len(contMap) == 0 {
		return nil, ErrNoContainerWithChadburnEnabled
	}

	var labels = make(map[string]map[string]string)

	for name, c := range contMap {
		if len(c.Labels) > 0 {
			containerLabels := make(map[string]string)
			for k, v := range c.Labels {
				// only include relevant labels
				if strings.HasPrefix(k, labelPrefix) {
					containerLabels[k] = v
				}
			}

			if len(containerLabels) > 0 {
				labels[name] = containerLabels
			}
		}
	}

	return labels, nil
}
