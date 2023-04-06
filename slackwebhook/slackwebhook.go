package slackwebhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mvndaai/ctxerr"
)

type (
	// Message is the json structure used by the webhook service
	Message struct {
		Text        string              `json:"text"`
		Username    string              `json:"username,omitempty"`
		Mrkdwn      bool                `json:"mrkdwn,omitempty"`
		Attachments []MessageAttachment `json:"attachments,omitempty"`
		Icon        string              `json:"icon_emoji"`
		Channel     string              `json:"channel,omitempty"`
	}

	// MessageAttachment is the attachment section of the message used for webhooks
	MessageAttachment struct {
		Color    string   `json:"color,omitempty"`   // Can either be one of 'good', 'warning', 'danger', or any hex color code
		Title    string   `json:"title,omitempty"`   // The title may not contain markup and will be escaped for you
		Pretext  string   `json:"pretext,omitempty"` // "Optional text that should appear above the formatted data",
		Text     string   `json:"text"`              // May contain standard message markup and must be escaped as normal. May be multi-line.
		MrkdwnIn []string `json:"mrkdwn_in,omitempty"`
	}
)

type contexter interface {
	Context() context.Context
}

type ContextHook func(context.Context, *Message)

type Config struct {
	// HTTPClient is used to make request to slack
	HTTPClient *http.Client
	// WebhookURL is the url from the slack webhook app to a channel
	WebhookURL string
	// Icon is a slack icon like :thumbsup: displayed when sending messages
	Icon string
	// Username is the name displayed when sending messages
	Username string
	// ColorError is the color errors are displayed in IsWarning is false or non-existent
	ColorError string
	// ColorWarning is the color warning are display in when IsWarning is true
	ColorWarning string
	// IsWarning tells if something is a warning
	IsWarning func(error) bool
	// Ignore tells if an error should be ignored and not sent to slack
	Ignore func(error) bool
	// LogError is a way to log an error not using ctxerr.Handle to avoid circular errors
	LogError func(error)
	// PrettyIndent is the indentation used when doing pretty JSON
	PrettyIndent string
	// NotPretty removes pretty print
	NotPretty bool
	// Fields function to get feilds to log. Defaults to ctxerr.AllFields
	Fields func(error) map[string]any
	// ContextHooks allow adding hooks to update things using the context
	ContextHooks []ContextHook
}

const (
	PrettyIndentTab    = "\t"
	PrettyIndentSpaces = "    "
)

const (
	ColorError   = "danger"
	ColorWarning = "warning"
)

// ToMessage converts an error to a message that can be sent to slack
func (c Config) ToMessage(err error) *Message {
	if err == nil {
		return nil
	}

	m := &Message{
		Text:     err.Error(),
		Username: c.Username,
		Icon:     c.Icon,
	}

	ff := c.Fields
	if ff == nil {
		ff = ctxerr.AllFields
	}
	fields := ff(err)
	if len(fields) > 0 {
		a := MessageAttachment{Color: c.ColorError}
		if c.IsWarning != nil && c.IsWarning(err) {
			a.Color = c.ColorWarning
		}

		if c.NotPretty {
			if jsonBody, err := json.Marshal(fields); err == nil {
				a.Text = fmt.Sprintf("```%s```", jsonBody)
			}
		} else {
			if jsonBody, err := json.MarshalIndent(fields, "", c.PrettyIndent); err == nil {
				a.Text = fmt.Sprintf("```%s```", jsonBody)
			}
		}
		if a.Text == "" {
			a.Text = fmt.Sprintf("```%s```", fields)
		}

		m.Attachments = append(m.Attachments, a)
	}

	if v, ok := err.(contexter); ok {
		ctx := v.Context()
		for _, hook := range c.ContextHooks {
			hook(ctx, m)
		}
	}
	return m
}

// SendSlackMessage sends a message to the slack webhook url in the config
func (c Config) SendSlackMessage(m *Message) {
	if m == nil {
		if c.LogError != nil {
			c.LogError(fmt.Errorf("no Message"))
		}
		return
	}

	if c.WebhookURL == "" {
		if c.LogError != nil {
			c.LogError(fmt.Errorf("no WebhookURL"))
		}
		return
	}

	if c.HTTPClient == nil {
		if c.LogError != nil {
			c.LogError(fmt.Errorf("nil HTTPClient"))
		}
		return
	}

	slackMessageBytes, err := json.Marshal(m)
	if err != nil {
		if c.LogError != nil {
			c.LogError(err)
		}
		return
	}

	_, err = c.HTTPClient.Post(c.WebhookURL, "application/json", bytes.NewBuffer(slackMessageBytes))
	if err != nil && c.LogError != nil {
		if c.LogError != nil {
			c.LogError(err)
		}
		return
	}
}

// HandleHook is a hook that can be added to ctxerr.AddHandleHook
func (c Config) HandleHook(err error) {
	if c.Ignore != nil && c.Ignore(err) {
		return
	}
	m := c.ToMessage(err)
	c.SendSlackMessage(m)
}
