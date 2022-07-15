package ctxerrlogrus_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mvndaai/ctxerr"
	ctxerrlogrus "github.com/mvndaai/ctxerrhelper/logrus"
	"github.com/sirupsen/logrus"
)

func TestHook(t *testing.T) {
	lg := logrus.New()
	lg.SetFormatter(&logrus.JSONFormatter{})
	sb := &strings.Builder{}
	lg.Out = sb

	lg.AddHook(ctxerrlogrus.NewContextHook())

	ctx := context.Background()
	ctx = ctxerr.SetField(ctx, "foo", "bar")
	lg.WithContext(ctx).Info("msg")

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(sb.String()), &m); err != nil {
		t.Error("could not unmarshall json", err)
	}

	if m["foo"] != "bar" {
		t.Error("could not find field in json")
	}
}

func TestHookWithLogLevels(t *testing.T) {
	lg := logrus.New()
	lg.SetFormatter(&logrus.JSONFormatter{})
	sb := &strings.Builder{}
	lg.Out = sb

	hook := ctxerrlogrus.NewContextHook()
	hook.LogLevels = []logrus.Level{logrus.WarnLevel}
	lg.AddHook(hook)

	ctx := context.Background()
	ctx = ctxerr.SetField(ctx, "foo", "bar")
	entry := lg.WithContext(ctx)

	entry.Info("msg")

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(sb.String()), &m); err != nil {
		t.Error("could not unmarshall json", err)
	}
	if _, ok := m["foo"]; ok {
		t.Error("fields should not have been added on info level:", m)
	}
	if m["level"] != "info" {
		t.Error("level should be info", m["level"])
	}

	sb.Reset()
	entry.Warn("msg")

	if err := json.Unmarshal([]byte(sb.String()), &m); err != nil {
		t.Error("could not unmarshall json", err)
	}
	if m["foo"] != "bar" {
		t.Error("could not find field in json", m)
	}
	if m["level"] != "warning" {
		t.Error("level should be warn not", m["level"])
	}
}

func TestHookWithConflictPrefix(t *testing.T) {
	tests := []struct {
		name        string
		prepend     bool
		prefix      string
		expectedKey string
		expectedVal interface{}
	}{
		{name: "no prepend", prepend: false, prefix: "", expectedKey: "foo", expectedVal: "bar"},
		{name: "default prepend", prepend: true, prefix: "", expectedKey: ctxerrlogrus.DefaultConflitPrefix + "foo", expectedVal: "bar"},
		{name: "default prepend original", prepend: true, prefix: "", expectedKey: "foo", expectedVal: "baz"},
		{name: "set prepend", prepend: true, prefix: "a.", expectedKey: "a.foo", expectedVal: "bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lg := logrus.New()
			lg.SetFormatter(&logrus.JSONFormatter{})
			sb := &strings.Builder{}
			lg.Out = sb

			hook := ctxerrlogrus.NewContextHook()
			hook.ConflitPrefix = tt.prefix
			hook.PrependConflicts = tt.prepend
			lg.AddHook(hook)

			f := func(lgr *logrus.Entry) {
				sb.Reset()
				lgr.Info("msg")

				var m map[string]interface{}
				if err := json.Unmarshal([]byte(sb.String()), &m); err != nil {
					t.Error("could not unmarshall json", err)
				}

				if m[tt.expectedKey] != tt.expectedVal {
					t.Errorf("expected [%s:%s] in json\n%v", tt.expectedKey, tt.expectedVal, m)
				}
			}

			ctx := context.Background()
			ctx = ctxerr.SetField(ctx, "foo", "bar")
			f(lg.WithContext(ctx).WithField("foo", "baz"))
			f(lg.WithField("foo", "baz").WithContext(ctx))
		})
	}
}

func ExampleNewContextHook() {
	lg := logrus.New()
	lg.AddHook(ctxerrlogrus.NewContextHook())

	ctx := ctxerr.SetField(context.Background(), "foo", "bar")
	lg.WithContext(ctx).Info("msg")
}
