package jobs

import "context"

type Job interface {
	Name() string
	Run(context.Context) error
}

type JobFunc struct {
	JobName string
	Fn      func(context.Context) error
}

func (j JobFunc) Name() string { return j.JobName }
func (j JobFunc) Run(ctx context.Context) error {
	if j.Fn == nil {
		return nil
	}
	return j.Fn(ctx)
}
