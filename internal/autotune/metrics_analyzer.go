package autotune

import "time"

type RuntimeMetrics struct {
	CPUPercent  float64
	MemoryMB    float64
	P95Latency  time.Duration
	QueueDepth  int64
	ErrorRate   float64
	WSConnCount int64
}

type Recommendation struct {
	WorkerPoolSize int
	DBMaxOpenConns int
	RedisPoolSize  int
	RateLimitRPS   float64
	CacheTTL       time.Duration
	WSScaleTarget  int64
}

func Analyze(m RuntimeMetrics, prev Recommendation) Recommendation {
	r := prev
	if r.WorkerPoolSize == 0 {
		r.WorkerPoolSize = 8
		r.DBMaxOpenConns = 20
		r.RedisPoolSize = 20
		r.RateLimitRPS = 100
		r.CacheTTL = time.Minute
		r.WSScaleTarget = 1000
	}
	if m.QueueDepth > int64(r.WorkerPoolSize*10) {
		r.WorkerPoolSize += 2
	}
	if m.CPUPercent > 80 || m.MemoryMB > 2048 {
		r.WorkerPoolSize = max(2, r.WorkerPoolSize-2)
		r.DBMaxOpenConns = max(5, r.DBMaxOpenConns-2)
		r.RedisPoolSize = max(5, r.RedisPoolSize-2)
	}
	if m.P95Latency > 500*time.Millisecond {
		r.RateLimitRPS *= 0.9
		r.CacheTTL += 30 * time.Second
	}
	if m.P95Latency < 120*time.Millisecond && m.ErrorRate < 0.01 {
		r.RateLimitRPS *= 1.05
		if r.CacheTTL > 30*time.Second {
			r.CacheTTL -= 5 * time.Second
		}
	}
	if m.WSConnCount > r.WSScaleTarget {
		r.WSScaleTarget += 500
	}
	return r
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
