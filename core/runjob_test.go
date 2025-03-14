package core

import (
	. "gopkg.in/check.v1"
)

type SuiteRunJob struct {
	mockClient *MockDockerClient
}

var _ = Suite(&SuiteRunJob{})

func (s *SuiteRunJob) SetUpTest(c *C) {
	s.mockClient = &MockDockerClient{}
}

func (s *SuiteRunJob) TestRun(c *C) {
	job := &OfficialRunJob{Client: s.mockClient}
	job.Image = "test"
	job.Command = `echo -a "foo bar"`
	job.User = "foo"
	job.TTY = true

	e := NewExecution()

	err := job.Run(&Context{Execution: e})
	c.Assert(err, IsNil)
}
