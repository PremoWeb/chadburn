package core

import (
	. "gopkg.in/check.v1"
)

type SuiteRunService struct {
	mockClient *MockDockerClient
}

var _ = Suite(&SuiteRunService{})

func (s *SuiteRunService) SetUpTest(c *C) {
	s.mockClient = &MockDockerClient{}
}

func (s *SuiteRunService) TestRun(c *C) {
	// Skip this test for now
	c.Skip("Skipping test for now")
}
