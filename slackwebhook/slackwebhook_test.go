package slackwebhook_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mvndaai/ctxerr"
	"github.com/mvndaai/ctxerrhelper/slackwebhook"
	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	t.Parallel()
	var m slackwebhook.Message
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		defer r.Body.Close()
		if err := json.Unmarshal(b, &m); err != nil {
			t.Error(err)
		}
	}))
	defer s.Close()

	webhookURL := s.URL
	in := ctxerr.NewInstance()
	conf := slackwebhook.Config{
		WebhookURL: webhookURL,
		HTTPClient: http.DefaultClient,
		LogError:   func(err error) { t.Error(err) },
	}
	in.AddHandleHook(conf.HandleHook)

	ctx := context.Background()
	err := fmt.Errorf("inside")
	err = in.WrapHTTP(ctx, err, "code", "action", http.StatusBadRequest, "wrap")
	in.Handle(err)

	assert.Equal(t, err.Error(), m.Text)
}

func TestToMessage(t *testing.T) {
	t.Parallel()
	ctxerrIn := ctxerr.Instance{}
	ctxerrIn.AddCreateHook(ctxerr.SetCodeHook)

	testFunc := func() {}

	tests := []struct {
		name     string
		config   slackwebhook.Config
		err      error
		expected *slackwebhook.Message
	}{
		{
			name:   "fmt wrapped",
			config: slackwebhook.Config{},
			err:    fmt.Errorf("b: %w", fmt.Errorf("a")),
			expected: &slackwebhook.Message{
				Text: "b: a",
			},
		},
		{
			name:   "basic ctxerr.New",
			config: slackwebhook.Config{},
			err:    ctxerrIn.New(context.Background(), "code", "msg"),
			expected: &slackwebhook.Message{
				Text: "msg",
				Attachments: []slackwebhook.MessageAttachment{{
					Text: "```{\n\"error_code\": \"code\"\n}```",
				}},
			},
		},
		{
			name:   "pretty tab indent",
			config: slackwebhook.Config{PrettyIndent: slackwebhook.PrettyIndentTab},
			err:    ctxerrIn.New(context.Background(), "code", "msg"),
			expected: &slackwebhook.Message{
				Text: "msg",
				Attachments: []slackwebhook.MessageAttachment{{
					Text: "```{\n\t\"error_code\": \"code\"\n}```",
				}},
			},
		},
		{
			name:   "pretty spaces indent",
			config: slackwebhook.Config{PrettyIndent: slackwebhook.PrettyIndentSpaces},
			err:    ctxerrIn.New(context.Background(), "code", "msg"),
			expected: &slackwebhook.Message{
				Text: "msg",
				Attachments: []slackwebhook.MessageAttachment{{
					Text: "```{\n    \"error_code\": \"code\"\n}```",
				}},
			},
		},
		{
			name:   "not pretty",
			config: slackwebhook.Config{NotPretty: true},
			err:    ctxerrIn.New(context.Background(), "code", "msg"),
			expected: &slackwebhook.Message{
				Text: "msg",
				Attachments: []slackwebhook.MessageAttachment{{
					Text: "```{\"error_code\":\"code\"}```",
				}},
			},
		},
		{
			name:   "non json",
			config: slackwebhook.Config{},
			err: func() error {
				ctx := context.Background()
				ctx = ctxerr.SetField(ctx, "func", testFunc)
				return ctxerrIn.New(ctx, "code", "msg")
			}(),
			expected: &slackwebhook.Message{
				Text: "msg",
				Attachments: []slackwebhook.MessageAttachment{{
					Text: fmt.Sprintf("```%s```", map[string]any{"error_code": "code", "func": testFunc}),
				}},
			},
		},
		{
			name: "warning",
			config: slackwebhook.Config{
				NotPretty:    true,
				IsWarning:    func(_ error) bool { return true },
				ColorWarning: "orange",
			},
			err: ctxerrIn.New(context.Background(), "code", "msg"),
			expected: &slackwebhook.Message{
				Text: "msg",
				Attachments: []slackwebhook.MessageAttachment{{
					Text:  "```{\"error_code\":\"code\"}```",
					Color: "orange",
				}},
			},
		},
		{
			name: "error color",
			config: slackwebhook.Config{
				NotPretty:    true,
				IsWarning:    func(_ error) bool { return false },
				ColorWarning: "orange",
				ColorError:   "red",
			},
			err: ctxerrIn.New(context.Background(), "code", "msg"),
			expected: &slackwebhook.Message{
				Text: "msg",
				Attachments: []slackwebhook.MessageAttachment{{
					Text:  "```{\"error_code\":\"code\"}```",
					Color: "red",
				}},
			},
		},
		{
			name:     "nil error",
			config:   slackwebhook.Config{},
			err:      nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m := tt.config.ToMessage(tt.err)
			assert.EqualValues(t, tt.expected, m)
		})
	}
}

func TestNegative(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		conf        slackwebhook.Config
		message     *slackwebhook.Message
		errContains string
	}{
		{
			name:        "nil message",
			conf:        slackwebhook.Config{},
			message:     nil,
			errContains: "no Message",
		},
		{
			name:        "no webhook",
			conf:        slackwebhook.Config{},
			message:     &slackwebhook.Message{},
			errContains: "no WebhookURL",
		},
		{
			name:        "no webhook",
			conf:        slackwebhook.Config{WebhookURL: "example.com"},
			message:     &slackwebhook.Message{},
			errContains: "nil HTTPClient",
		},
		{
			name:        "no webhook",
			conf:        slackwebhook.Config{WebhookURL: "example.com", HTTPClient: http.DefaultClient},
			message:     &slackwebhook.Message{},
			errContains: "unsupported protocol scheme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.conf.LogError = func(logErr error) { err = logErr }
			tt.conf.SendSlackMessage(tt.message)
			assert.Contains(t, err.Error(), tt.errContains)
		})
	}
}

func TestIngore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		conf          slackwebhook.Config
		err           error
		expectMessage bool
	}{
		{
			name:          "ignore all",
			conf:          slackwebhook.Config{Ignore: func(error) bool { return true }},
			err:           fmt.Errorf("err"),
			expectMessage: false,
		},
		{
			name:          "ignore none",
			conf:          slackwebhook.Config{Ignore: func(error) bool { return false }},
			err:           fmt.Errorf("err"),
			expectMessage: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var externalCall bool
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { externalCall = true }))
			defer s.Close()
			tt.conf.WebhookURL = s.URL
			tt.conf.HTTPClient = http.DefaultClient
			tt.conf.HandleHook(tt.err)
			assert.Equal(t, tt.expectMessage, externalCall)
		})
	}
}

func TestConfigFields(t *testing.T) {
	conf := slackwebhook.Config{
		NotPretty: true,
		Fields:    func(error) map[string]any { return map[string]any{"a": "b"} },
	}
	err := fmt.Errorf("err")
	m := conf.ToMessage(err)
	assert.Equal(t, err.Error(), m.Text)
	assert.EqualValues(t, "```{\"a\":\"b\"}```", m.Attachments[0].Text)
}

func TestContextHook(t *testing.T) {
	type ctxKeyT string
	var ctxKey ctxKeyT = "a"

	ctxhook := func(ctx context.Context, m *slackwebhook.Message) {
		v := ctx.Value(ctxKey)
		if v != nil {
			m.Username = fmt.Sprint(v)
		}
	}

	conf := slackwebhook.Config{
		NotPretty:    true,
		Fields:       func(error) map[string]any { return map[string]any{"a": "b"} },
		ContextHooks: []slackwebhook.ContextHook{ctxhook},
	}

	err := fmt.Errorf("err")
	m := conf.ToMessage(err)
	assert.Equal(t, "", m.Username)

	expectedUsername := "ctx-username"
	ctx := context.WithValue(context.Background(), ctxKey, expectedUsername)
	err = ctxerr.New(ctx, "", "")
	m = conf.ToMessage(err)
	assert.Equal(t, expectedUsername, m.Username)
}
