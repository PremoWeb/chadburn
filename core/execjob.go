package core

import (
	"fmt"
	"io"
	"reflect"

	"github.com/gobs/args"
)

// ExecJob represents a job that executes a command in a running Docker container
type ExecJob struct {
	BareJob   `mapstructure:",squash"`
	Client    DockerClient `json:"-"`
	Container string       `hash:"true"`
	User      string       `default:"root" hash:"true"`
	TTY       bool         `default:"false" hash:"true"`
	Workdir   string       `default:"" hash:"true"`
	// ContainerName is used for variable processing and is not included in the hash
	ContainerName string `hash:"-"`
}

// NewExecJob creates a new ExecJob
func NewExecJob(c DockerClient) *ExecJob {
	return &ExecJob{Client: c}
}

func (j *ExecJob) Run(ctx *Context) error {
	// Create variable context
	varContext := VariableContext{
		Container: ContainerInfo{
			Name: j.Container,
			ID:   j.Container, // We use the container name as ID for now
		},
	}

	// Get processed command with variables replaced
	processedCommand := j.GetProcessedCommand(varContext)

	// Parse command
	cmds := args.GetArgs(processedCommand)

	// Create exec config
	config := &ExecConfig{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          j.TTY,
		User:         j.User,
		WorkingDir:   j.Workdir,
	}

	// Create exec instance
	execID, err := j.Client.CreateExec(j.Container, cmds, config)
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

	switch inspect.ExitCode {
	case 0:
		return nil
	case -1:
		return ErrUnexpected
	default:
		return fmt.Errorf("error non-zero exit code: %d", inspect.ExitCode)
	}
}

// Returns a hash of all the job attributes. Used to detect changes
func (j *ExecJob) Hash() string {
	var hash string
	getHash(reflect.TypeOf(j).Elem(), reflect.ValueOf(j).Elem(), &hash)
	return hash
}
