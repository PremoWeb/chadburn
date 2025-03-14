package cli

import (
	"errors"
	"strings"
	"time"

	"github.com/PremoWeb/Chadburn/core"
	docker "github.com/fsouza/go-dockerclient"
)

var ErrNoContainerWithChadburnEnabled = errors.New("Couldn't find containers with label 'chadburn.enabled=true'")

type DockerHandler struct {
	dockerClient  *docker.Client
	notifier      dockerLabelsUpdate
	logger        core.Logger
	lifecycleJobs map[string]*LifecycleJobConfig // Map of lifecycle jobs
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

// TODO: Implement an interface so the code does not have to use third parties directly
func (c *DockerHandler) GetInternalDockerClient() *docker.Client {
	return c.dockerClient
}

func (c *DockerHandler) buildDockerClient() (*docker.Client, error) {
	d, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func NewDockerHandler(notifier dockerLabelsUpdate, logger core.Logger) (*DockerHandler, error) {
	c := &DockerHandler{}
	var err error
	c.dockerClient, err = c.buildDockerClient()
	c.notifier = notifier
	c.logger = logger
	c.lifecycleJobs = make(map[string]*LifecycleJobConfig)
	if err != nil {
		return nil, err
	}
	// Do a sanity check on docker
	_, err = c.dockerClient.Info()
	if err != nil {
		return nil, err
	}

	go c.watch()
	go c.watchEvents() // Start watching Docker events
	return c, nil
}

func (c *DockerHandler) watch() {
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
		}
	}
}

// watchEvents listens for Docker events and triggers lifecycle jobs
func (c *DockerHandler) watchEvents() {
	// Create a channel to receive Docker events
	events := make(chan *docker.APIEvents)

	// Add event listener
	err := c.dockerClient.AddEventListener(events)
	if err != nil {
		c.logger.Errorf("Failed to add event listener: %v", err)
		return
	}

	// Remove event listener when done
	defer c.dockerClient.RemoveEventListener(events)

	c.logger.Noticef("Started watching Docker events")

	// Listen for events
	for event := range events {
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

		// Get container name without leading slash
		containerName := strings.TrimPrefix(container.Name, "/")

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
	conts, err := c.dockerClient.ListContainers(docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": {requiredLabelFilter},
		},
	})
	if err != nil {
		return nil, err
	}

	// Also get containers with job-run labels
	jobRunConts, err := c.dockerClient.ListContainers(docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": {labelPrefix + "." + jobRun},
		},
	})
	if err != nil {
		return nil, err
	}

	// Combine the two lists, avoiding duplicates
	contMap := make(map[string]docker.APIContainers)
	for _, cont := range conts {
		if len(cont.Names) > 0 {
			name := strings.TrimPrefix(cont.Names[0], "/")
			contMap[name] = cont
		}
	}
	for _, cont := range jobRunConts {
		if len(cont.Names) > 0 {
			name := strings.TrimPrefix(cont.Names[0], "/")
			contMap[name] = cont
		}
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
