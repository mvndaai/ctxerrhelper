package ctxerrlogrus

import (
	"github.com/mvndaai/ctxerr"
	"github.com/sirupsen/logrus"
)

// NewContextHook creates a logrus hook that can be used to add ctxerr fields to logrus entries with a context.
func NewContextHook() *ContextHook { return &ContextHook{} }

// ContextHook implements logrus.Hook
type ContextHook struct {
	LogLevels []logrus.Level
	// PrependConflicts will prepend the key if it already exists in the entry.Data map
	PrependConflicts bool
	// ConflitPrefix will be prepended to the key if it already exists in the entry.Data map if prepend is enabled
	ConflitPrefix string
}

// DefaultConflitPrefix is the default prefix for keys that already exist on the entry when "PrependConflicts" enabled
const DefaultConflitPrefix = "ctxerr."

// Levels returns the log levels that this hook is enabled for
func (hook ContextHook) Levels() []logrus.Level {
	if len(hook.LogLevels) > 0 {
		return hook.LogLevels
	}
	return logrus.AllLevels
}

// Fire adds ctxerr fields to the logrus entry
func (hook ContextHook) Fire(entry *logrus.Entry) error {
	fields := ctxerr.Fields(entry.Context)
	for k, v := range fields {
		if _, ok := entry.Data[k]; ok && hook.PrependConflicts {
			if hook.ConflitPrefix == "" {
				hook.ConflitPrefix = DefaultConflitPrefix
			}
			k = hook.ConflitPrefix + k
		}
		entry.Data[k] = v
	}
	return nil
}
