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
