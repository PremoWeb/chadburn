package core

import (
	. "gopkg.in/check.v1"
)

const ContainerFixture = "test-container"

type SuiteExecJob struct {
	mockClient *MockDockerClient
}

var _ = Suite(&SuiteExecJob{})

func (s *SuiteExecJob) SetUpTest(c *C) {
	s.mockClient = &MockDockerClient{}
}

func (s *SuiteExecJob) TestRun(c *C) {
	job := &OfficialExecJob{Client: s.mockClient}
	job.Container = ContainerFixture
	job.Command = `echo -a "foo bar"`
	job.User = "foo"
	job.TTY = true

	e := NewExecution()

	err := job.Run(&Context{Execution: e})
	c.Assert(err, IsNil)
}
