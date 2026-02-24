package ai

import (
	"context"

	"zenx/internal/jobs"
)

func EnqueueInference(q *jobs.Queue, c Client, prompt string) {
	q.Enqueue(jobs.JobFunc{JobName: "ai_inference", Fn: func(ctx context.Context) error { _, err := c.Complete(ctx, prompt); return err }}, 0, 3)
}
