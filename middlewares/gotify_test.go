package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "gopkg.in/check.v1"
)

type SuiteGotify struct {
	BaseSuite
}

var _ = Suite(&SuiteGotify{})

func (s *SuiteGotify) TestNewGotifyEmpty(c *C) {
	c.Assert(NewGotify(&GotifyConfig{}), IsNil)
}

func (s *SuiteGotify) TestRunSuccess(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m gotifyMessage
		err := json.NewDecoder(r.Body).Decode(&m)
		c.Assert(err, Equals, nil)
		c.Assert(!strings.Contains(m.Message, "FAILED"), Equals, true)
	}))

	defer ts.Close()

	s.ctx.Start()
	s.ctx.Stop(nil)

	m := NewGotify(&GotifyConfig{GotifyWebhook: ts.URL})
	c.Assert(m.Run(s.ctx), IsNil)
}

func (s *SuiteGotify) TestRunSuccessFailed(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m gotifyMessage
		err := json.NewDecoder(r.Body).Decode(&m)
		c.Assert(err, Equals, nil)
		c.Assert(strings.Contains(m.Message, "FAILED"), Equals, true)
	}))

	defer ts.Close()

	s.ctx.Start()
	s.ctx.Stop(errors.New("foo"))

	m := NewGotify(&GotifyConfig{GotifyWebhook: ts.URL})
	c.Assert(m.Run(s.ctx), IsNil)
}

func (s *SuiteGotify) TestRunSuccessOnError(c *C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(true, Equals, false)
	}))

	defer ts.Close()

	s.ctx.Start()
	s.ctx.Stop(nil)

	m := NewGotify(&GotifyConfig{GotifyWebhook: ts.URL, GotifyOnlyOnError: true})
	c.Assert(m.Run(s.ctx), IsNil)
}
