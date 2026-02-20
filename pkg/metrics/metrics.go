package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	once sync.Once

	DBQueries = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "zenx",
		Subsystem: "db",
		Name:      "queries_total",
		Help:      "Total database queries by operation.",
	}, []string{"operation"})

	CacheOps = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "zenx",
		Subsystem: "cache",
		Name:      "operations_total",
		Help:      "Total cache operations by operation and result.",
	}, []string{"operation", "result"})

	JobRuns = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "zenx",
		Subsystem: "jobs",
		Name:      "runs_total",
		Help:      "Total job run attempts by job and status.",
	}, []string{"job", "status"})

	WebsocketConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "zenx",
		Subsystem: "websocket",
		Name:      "connections",
		Help:      "Current active websocket connections.",
	})
)

func Handler() http.Handler {
	once.Do(func() {})
	return promhttp.Handler()
}
