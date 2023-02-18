# ctxerrhelper

Helper packages to integrate with [ctxerr](https://github.com/mvndaai/ctxerr)

Each package has its own `go.mod` file to avoid importing unwanted package dependencies.  Use `go work use ./<dir>` to add a new one and `go work sync` to update.


| Package  | Integration | |
| - | - | - | 
|  [logrus](/logrus) | https://github.com/sirupsen/logrus |  [![DOC](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/mvndaai/ctxerrhelper/logrus) |
|  [slackwebhook](/slackwebhook) | [Incoming WebHooks](https://liveauctioneers.slack.com/apps/A0F7XDUAZ-incoming-webhooks) |  [![DOC](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/mvndaai/ctxerrhelper/slackwebhook) |

