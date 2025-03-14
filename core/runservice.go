package core

import (
	"fmt"
	"time"

	"github.com/gobs/args"
)

// Note: The ServiceJob is loosely inspired by https://github.com/alexellis/jaas/

// RunServiceJob represents a job that runs a Docker service
type RunServiceJob struct {
	BareJob `mapstructure:",squash"`
	Client  DockerClient `json:"-"`
	User    string       `default:"root"`
	TTY     bool         `default:"false"`
	// do not use bool values with "default:true" because if
	// user would set it to "false" explicitly, it still will be
	// changed to "true" https://github.com/mcuadros/ofelia/issues/135
	// so lets use strings here as workaround
	Delete  string `default:"true"`
	Image   string
	Network string
}

// NewRunServiceJob creates a new RunServiceJob
func NewRunServiceJob(c DockerClient) *RunServiceJob {
	return &RunServiceJob{Client: c}
}

func (j *RunServiceJob) Run(ctx *Context) error {
	if err := j.pullImage(); err != nil {
		return err
	}

	// Create service config
	config := &ServiceConfig{
		Name:  fmt.Sprintf("chadburn-%s", randomID()),
		Image: j.Image,
		Cmd:   args.GetArgs(j.GetCommand()),
		Labels: map[string]string{
			"chadburn.job": j.Name,
		},
	}

	// Create the service
	serviceID, err := j.Client.CreateService(config)
	if err != nil {
		return fmt.Errorf("error creating service: %s", err)
	}

	ctx.Logger.Noticef("Created service %s for job %s\n", serviceID, j.Name)

	// Watch the service
	if err := j.watchContainer(ctx, serviceID); err != nil {
		return err
	}

	// Delete the service if Delete is true
	if j.Delete == "true" {
		return j.deleteService(ctx, serviceID)
	}

	return nil
}

func (j *RunServiceJob) pullImage() error {
	// Pull the image directly
	if err := j.Client.PullImage(j.Image); err != nil {
		return fmt.Errorf("error pulling image %q: %s", j.Image, err)
	}

	return nil
}

func (j *RunServiceJob) watchContainer(ctx *Context, serviceID string) error {
	// Get service info
	service, err := j.Client.InspectService(serviceID)
	if err != nil {
		return fmt.Errorf("error inspecting service: %s", err)
	}

	// Log service info
	ctx.Logger.Debugf("Service %s created with name %s", serviceID, service.Name)

	// Wait for tasks to complete
	return j.waitForTasks(ctx, serviceID)
}

func (j *RunServiceJob) waitForTasks(ctx *Context, serviceID string) error {
	// Poll for task status
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Set timeout
	timeout := time.After(10 * time.Minute)

	for {
		select {
		case <-ticker.C:
			// Get tasks for the service
			tasks, err := j.Client.ListTasks(serviceID)
			if err != nil {
				return fmt.Errorf("error listing tasks: %s", err)
			}

			// Check if any tasks are running
			allCompleted := true
			for _, task := range tasks {
				ctx.Logger.Debugf("Task %s has status %s", task.ID, task.Status.State)

				if task.Status.State != "complete" && task.Status.State != "failed" {
					allCompleted = false
				}

				// If a task failed, return an error
				if task.Status.State == "failed" {
					return fmt.Errorf("task %s failed: %s", task.ID, task.Status.Err)
				}
			}

			// If all tasks are completed, we're done
			if allCompleted && len(tasks) > 0 {
				ctx.Logger.Noticef("Service %s has completed", serviceID)
				return nil
			}

		case <-timeout:
			return fmt.Errorf("timeout waiting for service %s to complete", serviceID)
		}
	}
}

func (j *RunServiceJob) deleteService(ctx *Context, serviceID string) error {
	ctx.Logger.Debugf("Removing service %s", serviceID)

	// Remove the service
	err := j.Client.RemoveService(serviceID)
	if err != nil {
		ctx.Logger.Warningf("Service %s cannot be removed. An error may have happened, or it might have been removed by another process", serviceID)
	}

	return nil
}
