package jobs

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type DistributedQueue struct {
	client     *redis.Client
	queueKey   string
	dlqKey     string
	maxRetries int
}

type payload struct {
	Name  string          `json:"name"`
	Body  json.RawMessage `json:"body"`
	Retry int             `json:"retry"`
	RunAt int64           `json:"run_at"`
}

func NewDistributedQueue(client *redis.Client, queueName string, maxRetries int) *DistributedQueue {
	if queueName == "" {
		queueName = "zenx:jobs"
	}
	return &DistributedQueue{client: client, queueKey: queueName, dlqKey: queueName + ":dlq", maxRetries: maxRetries}
}

func (q *DistributedQueue) Enqueue(ctx context.Context, name string, body any, delay time.Duration) error {
	raw, _ := json.Marshal(body)
	p := payload{Name: name, Body: raw, RunAt: time.Now().Add(delay).UnixMilli()}
	blob, _ := json.Marshal(p)
	return q.client.ZAdd(ctx, q.queueKey, redis.Z{Score: float64(p.RunAt), Member: blob}).Err()
}

func (q *DistributedQueue) Consume(ctx context.Context, handlers map[string]func(context.Context, json.RawMessage) error, poll time.Duration) {
	t := time.NewTicker(poll)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				q.popAndRun(ctx, handlers)
			}
		}
	}()
}

func (q *DistributedQueue) popAndRun(ctx context.Context, handlers map[string]func(context.Context, json.RawMessage) error) {
	now := time.Now().UnixMilli()
	items, err := q.client.ZRangeByScore(ctx, q.queueKey, &redis.ZRangeBy{Min: "-inf", Max: stringInt(now), Offset: 0, Count: 1}).Result()
	if err != nil || len(items) == 0 {
		return
	}
	blob := items[0]
	removed, err := q.client.ZRem(ctx, q.queueKey, blob).Result()
	if err != nil || removed == 0 {
		return
	}
	var p payload
	if json.Unmarshal([]byte(blob), &p) != nil {
		return
	}
	h, ok := handlers[p.Name]
	if !ok {
		_ = q.client.RPush(ctx, q.dlqKey, blob).Err()
		return
	}
	if err := h(ctx, p.Body); err != nil {
		p.Retry++
		if p.Retry > q.maxRetries {
			_ = q.client.RPush(ctx, q.dlqKey, blob).Err()
			return
		}
		p.RunAt = time.Now().Add(time.Duration(1<<p.Retry) * time.Second).UnixMilli()
		requeue, _ := json.Marshal(p)
		_ = q.client.ZAdd(ctx, q.queueKey, redis.Z{Score: float64(p.RunAt), Member: requeue}).Err()
	}
}

func stringInt(v int64) string { return strconv.FormatInt(v, 10) }
