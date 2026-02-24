package runtime

import (
	"context"
	"runtime"
	"sync/atomic"
	"time"
)

type Watchdog struct {
	goroutines atomic.Int64
	heapBytes  atomic.Uint64
	deadlocked atomic.Bool
}

func NewWatchdog() *Watchdog { return &Watchdog{} }

func (w *Watchdog) Start(ctx context.Context) {
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			w.goroutines.Store(int64(runtime.NumGoroutine()))
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			w.heapBytes.Store(ms.HeapAlloc)
		}
	}
}

func (w *Watchdog) Snapshot() (goroutines int64, heapBytes uint64, deadlocked bool) {
	return w.goroutines.Load(), w.heapBytes.Load(), w.deadlocked.Load()
}

func (w *Watchdog) MarkPotentialDeadlock(v bool) { w.deadlocked.Store(v) }
