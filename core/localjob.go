package core

import (
	"os/exec"
	"reflect"

	"github.com/gobs/args"
)

type LocalJob struct {
	BareJob     `mapstructure:",squash"`
	Dir         string
	Environment []string
	// Variables for command processing
	ContainerName string `hash:"-"`
	ContainerID   string `hash:"-"`
}

func NewLocalJob() *LocalJob {
	return &LocalJob{}
}

// Returns a hash of all the job attributes. Used to detect changes
func (j *LocalJob) Hash() string {
	var hash string
	getHash(reflect.TypeOf(j).Elem(), reflect.ValueOf(j).Elem(), &hash)
	return hash
}

func (j *LocalJob) Run(ctx *Context) error {
	cmd, err := j.buildCommand(ctx)
	if err != nil {
		return err
	}

	return cmd.Run()
}

func (j *LocalJob) buildCommand(ctx *Context) (*exec.Cmd, error) {
	// Create variable context
	varContext := VariableContext{
		Container: ContainerInfo{
			Name: j.ContainerName,
			ID:   j.ContainerID,
		},
	}

	// Get processed command with variables replaced
	processedCommand := j.GetProcessedCommand(varContext)

	args := args.GetArgs(processedCommand)
	bin, err := exec.LookPath(args[0])
	if err != nil {
		return nil, err
	}

	return &exec.Cmd{
		Path:   bin,
		Args:   args,
		Stdout: ctx.Execution.OutputStream,
		Stderr: ctx.Execution.ErrorStream,
		Env:    j.Environment,
		Dir:    j.Dir,
	}, nil
}
