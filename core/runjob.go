package core

import (
	"fmt"
	"reflect"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gobs/args"
)

var dockercfg *docker.AuthConfigurations

func init() {
	dockercfg, _ = docker.NewAuthConfigurationsFromDockerCfg()
}

type RunJob struct {
	BareJob   `mapstructure:",squash"`
	Client    *docker.Client `json:"-"`
	Container string         `hash:"true"`
	Image     string         `hash:"true"`
	User      string         `default:"root" hash:"true"`
	TTY       bool           `default:"false" hash:"true"`
	Delete    bool           `default:"true" hash:"true"`
	Network   string         `hash:"true"`
	Volume    []string       `hash:"true"`
}

func NewRunJob(c *docker.Client) *RunJob {
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
	err := j.Client.StartContainer(j.Container, nil)
	if err != nil {
		return err
	}

	if ctx.Execution.OutputStream != nil {
		err = j.Client.Logs(docker.LogsOptions{
			Container:    j.Container,
			OutputStream: ctx.Execution.OutputStream,
			ErrorStream:  ctx.Execution.ErrorStream,
			Stdout:       true,
			Stderr:       true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (j *RunJob) runContainer(ctx *Context) error {
	// Create variable context
	varContext := VariableContext{
		Container: ContainerInfo{
			Name: j.Container,
			ID:   j.Container, // We use the container name as ID for now
		},
	}

	// Get processed command with variables replaced
	processedCommand := j.GetProcessedCommand(varContext)

	c, err := j.buildContainer(processedCommand)
	if err != nil {
		return err
	}

	if err := j.createContainer(c); err != nil {
		return err
	}

	defer j.removeContainer(c)

	if err := j.startContainerWithExec(c, ctx.Execution); err != nil {
		return err
	}

	return j.watchContainer(c)
}

func (j *RunJob) buildContainer(processedCommand string) (*docker.Container, error) {
	var cmds []string
	if processedCommand != "" {
		// Use a try/catch approach to handle potential errors from args.GetArgs
		defer func() {
			if r := recover(); r != nil {
				// If there's a panic in GetArgs, just use the command as is
				cmds = []string{processedCommand}
			}
		}()

		cmds = args.GetArgs(processedCommand)
	}

	var err error
	var binds []string
	for _, volume := range j.Volume {
		binds = append(binds, volume)
	}

	config := &docker.Config{
		Image:        j.Image,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          j.TTY,
		Cmd:          cmds,
		User:         j.User,
	}

	hostConfig := &docker.HostConfig{
		Binds:       binds,
		NetworkMode: j.Network,
	}

	name := fmt.Sprintf("chadburn-%s", randomID())
	c, err := j.Client.CreateContainer(docker.CreateContainerOptions{
		Name:       name,
		Config:     config,
		HostConfig: hostConfig,
	})

	if err != nil {
		if err != docker.ErrNoSuchImage {
			return c, err
		}

		if err = j.pullImage(); err != nil {
			return nil, err
		}

		c, err = j.Client.CreateContainer(docker.CreateContainerOptions{
			Name:       name,
			Config:     config,
			HostConfig: hostConfig,
		})

		if err != nil {
			return c, err
		}
	}

	return c, nil
}

func (j *RunJob) createContainer(c *docker.Container) error {
	return nil
}

func (j *RunJob) removeContainer(c *docker.Container) error {
	if !j.Delete {
		return nil
	}

	return j.Client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            c.ID,
		RemoveVolumes: true,
		Force:         true,
	})
}

func (j *RunJob) startContainerWithExec(c *docker.Container, e *Execution) error {
	err := j.Client.StartContainer(c.ID, nil)
	if err != nil {
		return err
	}

	if e.OutputStream != nil {
		err = j.Client.Logs(docker.LogsOptions{
			Container:    c.ID,
			OutputStream: e.OutputStream,
			ErrorStream:  e.ErrorStream,
			Stdout:       true,
			Stderr:       true,
			Follow:       true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (j *RunJob) watchContainer(c *docker.Container) error {
	statusCode, err := j.Client.WaitContainer(c.ID)
	if err != nil {
		return err
	}

	if statusCode != 0 {
		return fmt.Errorf("error non-zero exit code: %d", statusCode)
	}

	return nil
}

func (j *RunJob) pullImage() error {
	o, a := buildPullOptions(j.Image)
	if err := j.Client.PullImage(o, a); err != nil {
		return err
	}

	return nil
}
