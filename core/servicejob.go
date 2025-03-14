package core

import (
	"fmt"
	"io"

	"github.com/gobs/args"
)

// ServiceJob represents a job that runs a command in a Docker container as a service
type ServiceJob struct {
	BareJob   `mapstructure:",squash"`
	Container string       `hash:"true"`
	Client    DockerClient `json:"-"`
	TTY       bool         `default:"false" hash:"true"`
	User      string       `default:"root" hash:"true"`
}

// NewServiceJob creates a new ServiceJob
func NewServiceJob(c DockerClient) *ServiceJob {
	return &ServiceJob{Client: c}
}

func (j *ServiceJob) Run(ctx *Context) error {
	// Create variable context
	varContext := VariableContext{
		Container: ContainerInfo{
			Name: j.Container,
			ID:   j.Container, // We use the container name as ID for now
		},
	}

	// Get processed command with variables replaced
	processedCommand := j.GetProcessedCommand(varContext)

	// Check if container exists and is running
	container, err := j.Client.InspectContainer(j.Container)
	if err != nil {
		return fmt.Errorf("error inspecting container: %s", err)
	}

	if !container.State.Running {
		return fmt.Errorf("unable to exec because container %q is not running", j.Container)
	}

	// Create exec config
	config := &ExecConfig{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          j.TTY,
		User:         j.User,
	}

	// Create exec instance
	execID, err := j.Client.CreateExec(j.Container, args.GetArgs(processedCommand), config)
	if err != nil {
		return fmt.Errorf("error creating exec: %s", err)
	}

	// Start exec
	reader, err := j.Client.StartExec(execID, true, true)
	if err != nil {
		return fmt.Errorf("error starting exec: %s", err)
	}
	defer reader.Close()

	// Copy output to the execution streams
	if ctx.Execution.OutputStream != nil {
		_, err = io.Copy(ctx.Execution.OutputStream, reader)
		if err != nil {
			return fmt.Errorf("error copying output: %s", err)
		}
	}

	// Inspect exec
	inspect, err := j.Client.InspectExec(execID)
	if err != nil {
		return fmt.Errorf("error inspecting exec: %s", err)
	}

	if inspect.ExitCode != 0 {
		return fmt.Errorf("error non-zero exit code: %d", inspect.ExitCode)
	}

	return nil
}

func (j *ServiceJob) GetProcessedCommand(varContext VariableContext) string {
	// Use the BareJob's implementation
	return j.BareJob.GetProcessedCommand(varContext)
}

// ... rest of the existing code ...
