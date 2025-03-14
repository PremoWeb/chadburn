package core

import (
	"fmt"
	"io"
	"reflect"

	"github.com/gobs/args"
)

// RunJob represents a job that runs a command in a Docker container
type RunJob struct {
	BareJob   `mapstructure:",squash"`
	Client    DockerClient `json:"-"`
	Container string       `hash:"true"`
	Image     string       `hash:"true"`
	User      string       `default:"root" hash:"true"`
	TTY       bool         `default:"false" hash:"true"`
	Delete    bool         `default:"true" hash:"true"`
	Network   string       `hash:"true"`
	Volume    []string     `hash:"true"`
}

// NewRunJob creates a new RunJob
func NewRunJob(c DockerClient) *RunJob {
	return &RunJob{Client: c}
}

func (j *RunJob) Run(ctx *Context) error {
	if j.Image == "" && j.Container == "" {
		return fmt.Errorf("unable to execute a job-run without a container or image")
	}

	if j.Image != "" {
		return j.runContainer(ctx)
	}

	return j.startContainer(ctx)
}

// Returns a hash of all the job attributes. Used to detect changes
func (j *RunJob) Hash() string {
	var hash string
	getHash(reflect.TypeOf(j).Elem(), reflect.ValueOf(j).Elem(), &hash)
	return hash
}

func (j *RunJob) startContainer(ctx *Context) error {
	// Check if container exists and is running
	container, err := j.Client.InspectContainer(j.Container)
	if err != nil {
		return fmt.Errorf("error inspecting container: %s", err)
	}

	if !container.State.Running {
		// Start the container
		err = j.Client.StartContainer(j.Container)
		if err != nil {
			return fmt.Errorf("error starting container: %s", err)
		}
	}

	// Create variable context
	varContext := VariableContext{
		Container: ContainerInfo{
			Name: j.Container,
			ID:   j.Container, // We use the container name as ID for now
		},
	}

	// Get processed command with variables replaced
	processedCommand := j.GetProcessedCommand(varContext)

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

func (j *RunJob) runContainer(ctx *Context) error {
	// Pull the image
	if err := j.Client.PullImage(j.Image); err != nil {
		return fmt.Errorf("error pulling image: %s", err)
	}

	// Create container config
	config := &ContainerConfig{
		Image:        j.Image,
		Cmd:          args.GetArgs(j.GetCommand()),
		AttachStdout: true,
		AttachStderr: true,
		Tty:          j.TTY,
		User:         j.User,
	}

	// Add host config if network or volumes are specified
	if j.Network != "" || len(j.Volume) > 0 {
		config.HostConfig = &HostConfig{
			NetworkMode: j.Network,
			Binds:       j.Volume,
		}
	}

	// Create the container
	container, err := j.Client.CreateContainer(config)
	if err != nil {
		return fmt.Errorf("error creating container: %s", err)
	}

	// Start the container
	err = j.Client.StartContainer(container.ID)
	if err != nil {
		return fmt.Errorf("error starting container: %s", err)
	}

	// Wait for the container to finish
	exitCode, err := j.Client.WaitContainer(container.ID)
	if err != nil {
		return fmt.Errorf("error waiting for container: %s", err)
	}

	// Remove the container if Delete is true
	if j.Delete {
		err = j.Client.RemoveContainer(container.ID)
		if err != nil {
			ctx.Logger.Errorf("error removing container: %s", err)
		}
	}

	if exitCode != 0 {
		return fmt.Errorf("error non-zero exit code: %d", exitCode)
	}

	return nil
}

func (j *RunJob) pullImage() error {
	// Pull the image directly
	if err := j.Client.PullImage(j.Image); err != nil {
		return err
	}

	return nil
}
