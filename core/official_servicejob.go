package core

import (
	"fmt"
	"io"
	"reflect"

	"github.com/gobs/args"
)

// OfficialServiceJob is an adapter for ServiceJob that uses the DockerClient interface
type OfficialServiceJob struct {
	BareJob   `mapstructure:",squash"`
	Client    DockerClient `json:"-"`
	Container string       `hash:"true"`
	User      string       `default:"root" hash:"true"`
	TTY       bool         `default:"false" hash:"true"`
}

// NewOfficialServiceJob creates a new OfficialServiceJob
func NewOfficialServiceJob(c DockerClient) *OfficialServiceJob {
	return &OfficialServiceJob{Client: c}
}

// Run executes the job
func (j *OfficialServiceJob) Run(ctx *Context) error {
	// Create variable context
	varContext := VariableContext{
		Container: ContainerInfo{
			Name: j.Container,
			ID:   j.Container, // We use the container name as ID for now
		},
	}

	// Get processed command with variables replaced
	processedCommand := j.GetProcessedCommand(varContext)

	// Create and start exec
	execID, err := j.buildExec(processedCommand)
	if err != nil {
		return err
	}

	if err := j.startExec(execID, ctx.Execution); err != nil {
		return err
	}

	return j.inspectExec(execID)
}

// Hash returns a hash of all the job attributes
func (j *OfficialServiceJob) Hash() string {
	var hash string
	getHash(reflect.TypeOf(j).Elem(), reflect.ValueOf(j).Elem(), &hash)
	return hash
}

// buildExec creates an exec instance in the container
func (j *OfficialServiceJob) buildExec(processedCommand string) (string, error) {
	// Check if container is running
	container, err := j.Client.InspectContainer(j.Container)
	if err != nil {
		return "", err
	}

	if !container.State.Running {
		return "", fmt.Errorf("unable to exec because container %q is not running", j.Container)
	}

	// Parse command
	var cmds []string
	if processedCommand != "" {
		cmds = args.GetArgs(processedCommand)
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
	execID, err := j.Client.CreateExec(j.Container, cmds, config)
	if err != nil {
		return "", fmt.Errorf("error creating exec: %s", err)
	}

	return execID, nil
}

// startExec starts an exec instance
func (j *OfficialServiceJob) startExec(execID string, e *Execution) error {
	reader, err := j.Client.StartExec(execID, true, true)
	if err != nil {
		return fmt.Errorf("error starting exec: %s", err)
	}
	defer reader.Close()

	// Copy output to the execution streams
	if e.OutputStream != nil {
		_, err = io.Copy(e.OutputStream, reader)
		if err != nil {
			return fmt.Errorf("error copying output: %s", err)
		}
	}

	return nil
}

// inspectExec inspects an exec instance
func (j *OfficialServiceJob) inspectExec(execID string) error {
	inspect, err := j.Client.InspectExec(execID)
	if err != nil {
		return fmt.Errorf("error inspecting exec: %s", err)
	}

	if inspect.ExitCode != 0 {
		return fmt.Errorf("error non-zero exit code: %d", inspect.ExitCode)
	}

	return nil
}

// GetProcessedCommand returns the processed command with variables replaced
func (j *OfficialServiceJob) GetProcessedCommand(varContext VariableContext) string {
	// Use the BareJob's implementation
	return j.BareJob.GetProcessedCommand(varContext)
}
