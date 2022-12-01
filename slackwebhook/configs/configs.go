package configs

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mvndaai/ctxerr"
	"github.com/mvndaai/ctxerrhelper/slackwebhook"
)

// WarningHTTPStatusCode choose warning or error based on http status codes
func WarningHTTPStatusCode(err error) bool {
	if f := ctxerr.AllFields(err); f != nil {
		code := f[ctxerr.FieldKeyStatusCode]
		if strings.HasPrefix(fmt.Sprint(code), "4") {
			return true
		}
	}
	return false
}

func ConfigSplitWarnings(webhookURL, username, icon string) slackwebhook.Config {
	return slackwebhook.Config{
		WebhookURL:   webhookURL,
		Username:     username,
		Icon:         icon,
		ColorError:   slackwebhook.ColorError,
		ColorWarning: slackwebhook.ColorWarning,
		PrettyIndent: slackwebhook.PrettyIndentTab,
		IsWarning:    WarningHTTPStatusCode,
		HTTPClient:   http.DefaultClient,
	}
}

func ConfigHTTPErrorOnly(webhookURL, username, icon string) slackwebhook.Config {
	return slackwebhook.Config{
		WebhookURL:   webhookURL,
		Username:     username,
		Icon:         icon,
		PrettyIndent: slackwebhook.PrettyIndentTab,
		Ignore:       WarningHTTPStatusCode,
		HTTPClient:   http.DefaultClient,
	}
}

func ConfigHTTPWarningOnly(webhookURL, username, icon string) slackwebhook.Config {
	return slackwebhook.Config{
		WebhookURL:   webhookURL,
		Username:     username,
		Icon:         icon,
		PrettyIndent: slackwebhook.PrettyIndentTab,
		Ignore:       func(err error) bool { return !WarningHTTPStatusCode(err) },
		HTTPClient:   http.DefaultClient,
	}
}
