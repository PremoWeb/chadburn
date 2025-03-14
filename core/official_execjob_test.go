package core

import (
	. "gopkg.in/check.v1"
)

type SuiteOfficialExecJob struct{}

var _ = Suite(&SuiteOfficialExecJob{})

func (s *SuiteOfficialExecJob) TestOfficialExecJob(c *C) {
	// Create a mock Docker client
	mockClient := &MockDockerClient{}

	// Create a new OfficialExecJob
	job := NewOfficialExecJob(mockClient)
	job.Container = "test-container"
	job.User = "root"
	job.TTY = true
	job.Command = "echo 'hello world'"

	// Verify the job properties
	c.Assert(job.Container, Equals, "test-container")
	c.Assert(job.User, Equals, "root")
	c.Assert(job.TTY, Equals, true)
	c.Assert(job.Command, Equals, "echo 'hello world'")
}
