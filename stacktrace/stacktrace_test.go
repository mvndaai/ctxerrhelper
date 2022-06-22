package stacktrace_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mvndaai/ctxerr"
	"github.com/mvndaai/ctxerrhelper/stacktrace"
)

func TestSetHook(t *testing.T) {
	ctxerr.AddCreateHook(stacktrace.CreateHook)
	ctxerr.AddCreateHook(stacktrace.CreateDebugHook)

	err := ctxerr.New(context.Background(), "c", "msg")
	// ctxerr.Handle(err)

	f := ctxerr.AllFields(err)
	t.Log(f)
	b, _ := json.Marshal(f)

	t.Error(string(b))

	// t.Error(err)
}
