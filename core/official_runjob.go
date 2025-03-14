package core

import (
	"fmt"
	"reflect"

	"github.com/gobs/args"
)

// OfficialRunJob is an adapter for RunJob that uses the DockerClient interface
type OfficialRunJob struct {
	BareJob     `mapstructure:",squash"`
	Client      DockerClient `json:"-"`
	Image       string       `hash:"true"`
	Network     string       `default:"bridge" hash:"true"`
	Pull        string       `default:"missing" hash:"true"`
	User        string       `default:"root" hash:"true"`
	TTY         bool         `default:"false" hash:"true"`
	Environment []string     `hash:"true"`
	Volumes     []string     `hash:"true"`
	WorkingDir  string       `hash:"true"`
}

// NewOfficialRunJob creates a new OfficialRunJob
func NewOfficialRunJob(c DockerClient) *OfficialRunJob {
	return &OfficialRunJob{Client: c}
}

// Run executes the job
func (j *OfficialRunJob) Run(ctx *Context) error {
	// Pull image if needed
	if j.Pull != "never" {
		if err := j.pullImage(); err != nil {
			return err
		}
	}

	// Create and start container
	containerID, err := j.startContainer(ctx)
	if err != nil {
		return err
	}

	// Wait for container to finish
	exitCode, err := j.Client.WaitContainer(containerID)
	if err != nil {
		return err
	}

	// Check exit code
	if exitCode != 0 {
		return fmt.Errorf("error non-zero exit code: %d", exitCode)
	}

	// Remove container
	return j.Client.RemoveContainer(containerID)
}

// Hash returns a hash of all the job attributes
func (j *OfficialRunJob) Hash() string {
	var hash string
	getHash(reflect.TypeOf(j).Elem(), reflect.ValueOf(j).Elem(), &hash)
	return hash
}

// pullImage pulls the Docker image
func (j *OfficialRunJob) pullImage() error {
	// Only pull if needed
	if j.Pull == "always" || j.Pull == "missing" {
		return j.Client.PullImage(j.Image)
	}
	return nil
}

// startContainer creates and starts a container
func (j *OfficialRunJob) startContainer(ctx *Context) (string, error) {
	// Parse command
	var cmds []string
	if j.Command != "" {
		cmds = args.GetArgs(j.Command)
	}

	// Create container config
	config := &ContainerConfig{
		Image:        j.Image,
		Cmd:          cmds,
		Tty:          j.TTY,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		User:         j.User,
		WorkingDir:   j.WorkingDir,
		Env:          j.Environment,
		HostConfig: &HostConfig{
			Binds:       j.Volumes,
			NetworkMode: j.Network,
		},
	}

	// Create container
	container, err := j.Client.CreateContainer(config)
	if err != nil {
		return "", fmt.Errorf("error creating container: %s", err)
	}

	// Start container
	if err := j.Client.StartContainer(container.ID); err != nil {
		return container.ID, fmt.Errorf("error starting container: %s", err)
	}

	// TODO: Implement logs streaming for containers
	// This would require extending the DockerClient interface with a ContainerLogs method

	return container.ID, nil
}
