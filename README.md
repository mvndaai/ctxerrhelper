# ctxerrhelper

Helper packages to integrate with [ctxerr](https://github.com/mvndaai/ctxerr)

Each package has its own `go.mod` file to avoid importing unwanted package dependencies.  Use `go work use ./<dir>` to add a new one and `go work sync` to update.


| Package  | Integration | |
| - | - | - |
|  [logrus](/logrus) | https://github.com/sirupsen/logrus |  [![DOC](https://img.shields.io/github/v/tag/mvndaai/ctxerrhelper?filter=logrus%2F*)](https://pkg.go.dev/github.com/mvndaai/ctxerrhelper/logrus) |
|  [slackwebhook](/slackwebhook) | [Incoming WebHooks](https://liveauctioneers.slack.com/apps/A0F7XDUAZ-incoming-webhooks) |  [![DOC](https://img.shields.io/github/v/tag/mvndaai/ctxerrhelper?filter=slackwebhook%2F*)](https://pkg.go.dev/github.com/mvndaai/ctxerrhelper/slackwebhook) |
|  [echo](/echo) | https://echo.labstack.com/ |  [![DOC](https://img.shields.io/github/v/tag/mvndaai/ctxerrhelper?filter=echo%2F*)](https://pkg.go.dev/github.com/mvndaai/ctxerrhelper/echo) |
|  [opencensus](/opencensus) | https://pkg.go.dev/go.opencensus.io |  [![DOC](https://img.shields.io/github/v/tag/mvndaai/ctxerrhelper?filter=opencensus%2F*)](https://pkg.go.dev/github.com/mvndaai/ctxerrhelper/opencensus) |
