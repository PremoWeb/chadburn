package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PremoWeb/Chadburn/core"
)

type GotifyConfig struct {
	GotifyWebhook     string `gcfg:"gotify-webhook" mapstructure:"gotify-webhook"`
	GotifyOnlyOnError bool   `gcfg:"gotify-only-on-error" mapstructure:"gotify-only-on-error"`
	GotifyPriority    int64  `gcfg:"gotify-priority" mapstructure:"gotify-priority"`
}

// NewGotify returns a Gotify middleware if the given configuration is not empty
func NewGotify(c *GotifyConfig) core.Middleware {
	var m core.Middleware
	if !IsEmpty(c) {
		m = &Gotify{*c}
	}

	return m
}

type Gotify struct {
	GotifyConfig
}

// ContinueOnStop return allways true, we want always report the final status
func (m *Gotify) ContinueOnStop() bool {
	return true
}

func (m *Gotify) Run(ctx *core.Context) error {
	err := ctx.Run()
	ctx.Stop(err)

	if ctx.Execution.Failed() || !m.GotifyOnlyOnError {
		m.pushMessage(ctx)
	}

	return err
}

func (m *Gotify) pushMessage(ctx *core.Context) {
	content, _ := json.Marshal(m.buildMessage(ctx))

	r, err := http.Post(m.GotifyWebhook, "application/json", bytes.NewReader(content))
	if err != nil {
		ctx.Logger.Errorf("Gotify error calling %q error: %q", m.GotifyWebhook, err)
	} else if r.StatusCode != 200 {
		ctx.Logger.Errorf("Gotify error non-200 status code calling %q", m.GotifyWebhook)
	}
}

func (m *Gotify) buildMessage(ctx *core.Context) *gotifyMessage {
	msg := &gotifyMessage{Title: ctx.Job.GetName(), Priority: m.GotifyPriority, Extras: gotifyMessageExtras{ClientDisplay: gotifyMessageExtrasDisplay{ContentType: "text/markdown"}}}

	msg.Message = fmt.Sprintf(
		"Job *%q* finished in *%s*, command `%s`",
		ctx.Job.GetName(), ctx.Execution.Duration(), ctx.Job.GetCommand(),
	)

	if ctx.Execution.Failed() {
		msg.Message = "FAILED: " + msg.Message
	} else if ctx.Execution.Skipped() {
		msg.Message = "Skipped: " + msg.Message
	}
	return msg
}

type gotifyMessage struct {
	Title    string              `json:"title"`
	Message  string              `json:"message"`
	Priority int64               `json:"priority"`
	Extras   gotifyMessageExtras `json:"extras"`
}

type gotifyMessageExtras struct {
	ClientDisplay gotifyMessageExtrasDisplay `json:"client::display"`
}

type gotifyMessageExtrasDisplay struct {
	ContentType string `json:"contentType"`
}
