package observability

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

type AsyncLogger struct {
	out  chan slog.Record
	base *slog.Logger
	wg   sync.WaitGroup
}

func NewAsyncLogger(size int) *AsyncLogger {
	l := &AsyncLogger{out: make(chan slog.Record, size), base: slog.New(slog.NewJSONHandler(os.Stdout, nil))}
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for r := range l.out {
			_ = l.base.Handler().Handle(context.Background(), r)
		}
	}()
	return l
}
func (l *AsyncLogger) Close() { close(l.out); l.wg.Wait() }
