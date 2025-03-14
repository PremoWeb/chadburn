package core

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gobs/args"
)

type ServiceJob struct {
	BareJob   `mapstructure:",squash"`
	Container string         `hash:"true"`
	Client    *docker.Client `json:"-"`
	TTY       bool           `default:"false" hash:"true"`
	User      string         `default:"root" hash:"true"`
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

	exec, err := j.buildExec(processedCommand)
	if err != nil {
		return err
	}

	if err := j.startExec(exec, ctx.Execution); err != nil {
		return err
	}

	return j.inspectExec(exec.ID)
}

func (j *ServiceJob) buildExec(processedCommand string) (*docker.Exec, error) {
	container, err := j.Client.InspectContainer(j.Container)
	if err != nil {
		return nil, err
	}

	if !container.State.Running {
		return nil, fmt.Errorf("unable to exec because container %q is not running", j.Container)
	}

	var cmds []string
	if processedCommand != "" {
		cmds = args.GetArgs(processedCommand)
	}

	exec, err := j.Client.CreateExec(docker.CreateExecOptions{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          j.TTY,
		Cmd:          cmds,
		Container:    j.Container,
		User:         j.User,
	})

	if err != nil {
		return exec, fmt.Errorf("error creating exec: %s", err)
	}

	return exec, nil
}

func (j *ServiceJob) startExec(exec *docker.Exec, e *Execution) error {
	err := j.Client.StartExec(exec.ID, docker.StartExecOptions{
		OutputStream: e.OutputStream,
		ErrorStream:  e.ErrorStream,
	})

	if err != nil {
		return err
	}

	return nil
}

func (j *ServiceJob) inspectExec(execID string) error {
	execInspect, err := j.Client.InspectExec(execID)
	if err != nil {
		return err
	}

	if execInspect.ExitCode != 0 {
		return fmt.Errorf("error non-zero exit code: %d", execInspect.ExitCode)
	}

	return nil
}

func (j *ServiceJob) GetProcessedCommand(varContext VariableContext) string {
	// Use the BareJob's implementation
	return j.BareJob.GetProcessedCommand(varContext)
}

func NewServiceJob(c *docker.Client) *ServiceJob {
	return &ServiceJob{Client: c}
}

// ... rest of the existing code ...
