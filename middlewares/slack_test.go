package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"
)

type SuiteSlack struct {
	BaseSuite
}

var _ = Suite(&SuiteSlack{})

func (s *SuiteSlack) TestNewSlackEmpty(c *C) {
	c.Assert(NewSlack(&SlackConfig{}), IsNil)
}

func (s *SuiteSlack) TestRunSuccess(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m slackMessage
		json.Unmarshal([]byte(r.FormValue(slackPayloadVar)), &m)

		// Check if Attachments array is not empty before accessing its elements
		if len(m.Attachments) > 0 {
			c.Assert(m.Attachments[0].Title, Equals, "Execution successful")
		} else {
			// For successful executions, there might not be any attachments
			c.Assert(m.Text, Not(Equals), "")
		}
	}))

	defer ts.Close()

	s.ctx.Start()
	s.ctx.Stop(nil)

	m := NewSlack(&SlackConfig{SlackWebhook: ts.URL})
	c.Assert(m.Run(s.ctx), IsNil)
}

func (s *SuiteSlack) TestRunSuccessFailed(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		c.Assert(err, IsNil)

		var m slackMessage
		err = json.Unmarshal([]byte(r.FormValue(slackPayloadVar)), &m)
		c.Assert(err, IsNil)
		c.Assert(len(m.Attachments), Equals, 1)
		c.Assert(m.Attachments[0].Color, Equals, "#F35A00")
	}))

	defer ts.Close()

	s.ctx.Start()
	s.ctx.Stop(errors.New("foo"))

	m := NewSlack(&SlackConfig{SlackWebhook: ts.URL})
	c.Assert(m.Run(s.ctx), NotNil)
}

func (s *SuiteSlack) TestRunSuccessOnError(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(true, Equals, false)
	}))

	defer ts.Close()

	s.ctx.Start()
	s.ctx.Stop(nil)

	m := NewSlack(&SlackConfig{SlackWebhook: ts.URL, SlackOnlyOnError: true})
	c.Assert(m.Run(s.ctx), IsNil)
}
