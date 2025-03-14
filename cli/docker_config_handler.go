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
	dockerClient *docker.Client
	notifier     dockerLabelsUpdate
	logger       core.Logger
}

type dockerLabelsUpdate interface {
	dockerLabelsUpdate(map[string]map[string]string)
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
	if err != nil {
		return nil, err
	}
	// Do a sanity check on docker
	_, err = c.dockerClient.Info()
	if err != nil {
		return nil, err
	}

	go c.watch()
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
