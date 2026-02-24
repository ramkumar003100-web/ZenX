package plugins

import "context"

type SandboxPolicy struct {
	AllowNetwork   bool
	AllowFileWrite bool
}
type Executor interface {
	Exec(context.Context, string, ...string) error
}
