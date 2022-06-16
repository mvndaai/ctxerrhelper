package stacktrace

import "context"

func CreateHook(ctx context.Context, code string, wrapping error) context.Context {
	return ctx
}
