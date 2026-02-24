package security

import (
	"context"
	"sync/atomic"
	"time"
)

type CertReloader struct{ cert atomic.Pointer[[]byte] }

func (r *CertReloader) Start(ctx context.Context, interval time.Duration, fetch func(context.Context) ([]byte, error)) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if c, err := fetch(ctx); err == nil {
					r.cert.Store(&c)
				}
			}
		}
	}()
}
