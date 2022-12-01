package configs_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mvndaai/ctxerr"
	"github.com/mvndaai/ctxerrhelper/slackwebhook"
	"github.com/mvndaai/ctxerrhelper/slackwebhook/configs"
	"github.com/stretchr/testify/assert"
)

func TestIngore(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		conf          slackwebhook.Config
		err           error
		expectMessage bool
	}{
		{
			name:          "ConfigHTTPErrorOnly error",
			conf:          configs.ConfigHTTPErrorOnly("", "", ""),
			err:           ctxerr.New(context.Background(), "c"),
			expectMessage: true,
		},
		{
			name:          "ConfigHTTPErrorOnly warning",
			conf:          configs.ConfigHTTPErrorOnly("", "", ""),
			err:           ctxerr.NewHTTP(context.Background(), "c", "", http.StatusBadRequest),
			expectMessage: false,
		},
		{
			name:          "ConfigHTTPWarningOnly error",
			conf:          configs.ConfigHTTPWarningOnly("", "", ""),
			err:           ctxerr.New(context.Background(), "c"),
			expectMessage: false,
		},
		{
			name:          "ConfigHTTPWarningOnly warning",
			conf:          configs.ConfigHTTPWarningOnly("", "", ""),
			err:           ctxerr.NewHTTP(context.Background(), "c", "", http.StatusBadRequest),
			expectMessage: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var externalCall bool
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { externalCall = true }))
			defer s.Close()
			tt.conf.WebhookURL = s.URL
			tt.conf.HandleHook(tt.err)
			assert.Equal(t, tt.expectMessage, externalCall)
		})
	}
}

func TestToMessage(t *testing.T) {
	t.Parallel()
	ctxerrIn := ctxerr.Instance{}
	ctxerrIn.AddCreateHook(ctxerr.SetCodeHook)

	tests := []struct {
		name     string
		config   slackwebhook.Config
		err      error
		expected *slackwebhook.Message
	}{
		{
			name:   "split warning config error",
			config: configs.ConfigSplitWarnings("", "", ""),
			err:    ctxerrIn.New(context.Background(), "code", "msg"),
			expected: &slackwebhook.Message{
				Text: "msg",
				Attachments: []slackwebhook.MessageAttachment{{
					Text:  "```{\n\t\"error_code\": \"code\"\n}```",
					Color: slackwebhook.ColorError,
				}},
			},
		},
		{
			name:   "split warning config warning",
			config: configs.ConfigSplitWarnings("", "", ""),
			err:    ctxerrIn.NewHTTP(context.Background(), "code", "", http.StatusBadRequest, "msg"),
			expected: &slackwebhook.Message{
				Text: "msg",
				Attachments: []slackwebhook.MessageAttachment{{
					Text:  "```{\n\t\"error_code\": \"code\",\n\t\"error_status_code\": 400\n}```",
					Color: slackwebhook.ColorWarning,
				}},
			},
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

func TestRealWebhhokURL(t *testing.T) {
	var webhookURL, username, icon string
	username = "Unit Test: Split"
	icon = ":thumbsup:"

	if webhookURL == "" {
		t.Skip("skipping real test")
	}

	ctxerrIn := ctxerr.Instance{}
	ctxerrIn.AddCreateHook(ctxerr.SetCodeHook)

	config := configs.ConfigSplitWarnings(webhookURL, username, icon)
	ctxerrIn.AddHandleHook(config.HandleHook)
	configError := configs.ConfigHTTPErrorOnly(webhookURL, "Unit Test: Error", ":thumbsdown:")
	configError.ColorError = "#8DDCA4"
	ctxerrIn.AddHandleHook(configError.HandleHook)
	configWarning := configs.ConfigHTTPWarningOnly(webhookURL, "Unit Test: Warning", ":arrow_up:")
	configWarning.ColorError = "#E76D83"
	ctxerrIn.AddHandleHook(configWarning.HandleHook)

	err := ctxerrIn.New(context.Background(), "test_error", "test error")
	ctxerrIn.Handle(err)

	err = ctxerrIn.NewHTTP(context.Background(), "test_warn", "action", http.StatusBadRequest, "test warning")
	ctxerrIn.Handle(err)

	longMsg := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
	err = ctxerrIn.New(context.Background(), "code", longMsg)
	ctxerrIn.Handle(err)

	t.Error("see logs")
}
