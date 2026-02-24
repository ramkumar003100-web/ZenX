package testing

import (
	"net/http"
	"sync"
)

func RunLoad(url string, n, c int) int {
	var wg sync.WaitGroup
	sem := make(chan struct{}, c)
	success := 0
	var mu sync.Mutex
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode < 500 {
				mu.Lock()
				success++
				mu.Unlock()
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}()
	}
	wg.Wait()
	return success
}
