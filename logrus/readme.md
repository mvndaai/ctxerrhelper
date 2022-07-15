# [logrus](https://github.com/sirupsen/logrus)

## Context Hook

If the entry used [`WithContext`](https://pkg.go.dev/github.com/sirupsen/logrus#WithContext) the hook add the `ctxerr` fields to the log.

```go
import (
	ctxerrlogrus "github.com/mvndaai/ctxerrhelper/logrus"
)

func main() {
    lg := logrus.New()
    lg.AddHook(ctxerrlogrus.NewContextHook())

    ctx := ctxerr.SetField(context.Background(), "foo", "bar")
    lg.WithContext(ctx).Info("msg")
}
```

