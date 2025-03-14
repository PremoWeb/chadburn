package cli

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/PremoWeb/Chadburn/core"
)

// OfficialDockerHandler is a Docker handler that uses the official Docker client
type OfficialDockerHandler struct {
	dockerClient  core.DockerClient
	notifier      dockerLabelsUpdate
	logger        core.Logger
	lifecycleJobs map[string]*LifecycleJobConfig // Map of lifecycle jobs
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewOfficialDockerHandler creates a new Docker handler using the official Docker client
func NewOfficialDockerHandler(notifier dockerLabelsUpdate, logger core.Logger) (*OfficialDockerHandler, error) {
	ctx, cancel := context.WithCancel(context.Background())

	client, err := core.NewDockerClient()
	if err != nil {
		cancel() // Cancel the context if there's an error
		return nil, err
	}

	handler := &OfficialDockerHandler{
		dockerClient:  client,
		notifier:      notifier,
		logger:        logger,
		lifecycleJobs: make(map[string]*LifecycleJobConfig),
		ctx:           ctx,
		cancel:        cancel,
	}

	go handler.watch()
	go handler.watchEvents()

	return handler, nil
}

// GetLifecycleJobs returns the lifecycle jobs
func (c *OfficialDockerHandler) GetLifecycleJobs() map[string]*LifecycleJobConfig {
	return c.lifecycleJobs
}

// SetLifecycleJobs sets the lifecycle jobs
func (c *OfficialDockerHandler) SetLifecycleJobs(jobs map[string]*LifecycleJobConfig) {
	c.lifecycleJobs = jobs
}

// Close closes the Docker client
func (c *OfficialDockerHandler) Close() error {
	c.cancel()
	return c.dockerClient.Close()
}

// watch polls for changes in Docker containers
func (c *OfficialDockerHandler) watch() {
	// Poll for changes
	tick := time.Tick(10000 * time.Millisecond)
	for {
		select {
		case <-tick:
			labels, err := c.GetDockerLabels()
			// Do not print or care if there is no container up right now
			if err != nil && !errors.Is(err, ErrNoContainerWithChadburnEnabled) {
				c.logger.Debugf("%v", err)
			}
			c.notifier.dockerLabelsUpdate(labels)
		case <-c.ctx.Done():
			return
		}
	}
}

// watchEvents watches Docker events
func (c *OfficialDockerHandler) watchEvents() {
	eventCh := make(chan *core.DockerEvent)
	errCh := make(chan error)

	// Start watching events
	c.dockerClient.WatchEvents(c.ctx, eventCh, errCh)

	c.logger.Noticef("Started watching Docker events")

	// Process events
	for {
		select {
		case err := <-errCh:
			if err != nil {
				c.logger.Errorf("Docker events error: %v", err)
				if err == context.Canceled {
					return
				}
			}
		case event := <-eventCh:
			// Process container events
			if event.Type == "container" {
				// Get container name from attributes
				containerName, ok := event.Attributes["name"]
				if !ok {
					continue
				}

				// Process the event
				switch event.Action {
				case "start":
					c.logger.Debugf("Container %s started", containerName)
					c.processLifecycleEvent(containerName, event.ID, core.ContainerStart)
				case "die", "stop":
					c.logger.Debugf("Container %s stopped", containerName)
					c.processLifecycleEvent(containerName, event.ID, core.ContainerStop)
				}
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// processLifecycleEvent processes a container lifecycle event
func (c *OfficialDockerHandler) processLifecycleEvent(containerName, containerID string, eventType core.LifecycleEventType) {
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

// GetDockerLabels gets Docker labels from containers
func (c *OfficialDockerHandler) GetDockerLabels() (map[string]map[string]string, error) {
	// First, get containers with the required label
	containers, err := c.dockerClient.ListContainers(map[string][]string{
		"label": {requiredLabelFilter},
	})
	if err != nil {
		return nil, err
	}

	// Also get containers with job-run labels
	jobRunContainers, err := c.dockerClient.ListContainers(map[string][]string{
		"label": {labelPrefix + "." + jobRun},
	})
	if err != nil {
		return nil, err
	}

	// Combine the two lists, avoiding duplicates
	contMap := make(map[string]core.Container)
	for _, cont := range containers {
		contMap[cont.Name] = cont
	}
	for _, cont := range jobRunContainers {
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
