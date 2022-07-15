package ctxerrlogrus

import (
	"github.com/mvndaai/ctxerr"
	"github.com/sirupsen/logrus"
)

func NewContextHook() *ContextHook { return &ContextHook{} }

type ContextHook struct {
	LogLevels []logrus.Level
}

func (hook ContextHook) Levels() []logrus.Level {
	if len(hook.LogLevels) > 0 {
		return hook.LogLevels
	}
	return logrus.AllLevels
}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	fields := ctxerr.Fields(entry.Context)
	for k, v := range fields {
		entry.Data[k] = v
	}
	return nil
}
