package jobs

import (
	"container/heap"
	"context"
	"math"
	"sync"
	"time"
)

type queuedJob struct {
	job      Job
	runAt    time.Time
	retry    int
	maxRetry int
	index    int
}

type pq []*queuedJob

func (p pq) Len() int           { return len(p) }
func (p pq) Less(i, j int) bool { return p[i].runAt.Before(p[j].runAt) }
func (p pq) Swap(i, j int)      { p[i], p[j] = p[j], p[i]; p[i].index, p[j].index = i, j }
func (p *pq) Push(x any)        { *p = append(*p, x.(*queuedJob)) }
func (p *pq) Pop() any {
	old := *p
	n := len(old)
	item := old[n-1]
	*p = old[:n-1]
	return item
}

type Queue struct {
	mu      sync.Mutex
	jobs    pq
	notify  chan struct{}
	workers int
}

func NewQueue(workers int) *Queue {
	q := &Queue{notify: make(chan struct{}, 1), workers: workers}
	heap.Init(&q.jobs)
	return q
}

func (q *Queue) Enqueue(job Job, delay time.Duration, retries int) {
	q.mu.Lock()
	heap.Push(&q.jobs, &queuedJob{job: job, runAt: time.Now().Add(delay), maxRetry: retries})
	q.mu.Unlock()
	select {
	case q.notify <- struct{}{}:
	default:
	}
}

func (q *Queue) ScheduleRecurring(ctx context.Context, interval time.Duration, job Job, retries int) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				q.Enqueue(job, 0, retries)
			}
		}
	}()
}

func (q *Queue) Start(ctx context.Context) {
	for i := 0; i < q.workers; i++ {
		go q.worker(ctx)
	}
}

func (q *Queue) worker(ctx context.Context) {
	for {
		job := q.next(ctx)
		if job == nil {
			return
		}
		if err := job.job.Run(ctx); err != nil && job.retry < job.maxRetry {
			job.retry++
			backoff := time.Duration(math.Pow(2, float64(job.retry))) * time.Second
			q.Enqueue(job.job, backoff, job.maxRetry)
		}
	}
}

func (q *Queue) next(ctx context.Context) *queuedJob {
	for {
		q.mu.Lock()
		if len(q.jobs) > 0 {
			next := q.jobs[0]
			wait := time.Until(next.runAt)
			if wait <= 0 {
				heap.Pop(&q.jobs)
				q.mu.Unlock()
				return next
			}
			q.mu.Unlock()
			t := time.NewTimer(wait)
			select {
			case <-ctx.Done():
				t.Stop()
				return nil
			case <-q.notify:
				t.Stop()
			case <-t.C:
			}
			continue
		}
		q.mu.Unlock()

		select {
		case <-ctx.Done():
			return nil
		case <-q.notify:
		}
	}
}
