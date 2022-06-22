package stacktrace

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/mvndaai/ctxerr"
)

const FieldKeyStackTrace = "error_stack_trace"

func CreateHook(ctx context.Context, code string, wrapping error) context.Context {
	// if ctxerr.HasField(FieldKeyStackTrace) {
	// 	return ctx
	// }

	// stackSlice := make([]byte, 512)
	stackSlice := make([]byte, 1024)
	s := runtime.Stack(stackSlice, false)
	fmt.Println("-----------------------")
	fmt.Println(s)
	fmt.Println(string(stackSlice[0:s]))
	fmt.Println("-----------------------")

	return ctxerr.SetField(ctx, FieldKeyStackTrace, string(stackSlice[0:s]))
}

func CreateDebugHook(ctx context.Context, code string, wrapping error) context.Context {
	return ctxerr.SetField(ctx, "error_stack_trace_debug", string(debug.Stack()))
}

// TODO see how pkg errors gets thier stack
