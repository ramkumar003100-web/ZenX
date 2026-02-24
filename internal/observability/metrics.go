package observability

import "github.com/prometheus/client_golang/prometheus"

var (
	RouteLatency      = prometheus.NewHistogramVec(prometheus.HistogramOpts{Namespace: "zenx", Subsystem: "http", Name: "route_latency_seconds", Help: "Per-route latency", Buckets: prometheus.DefBuckets}, []string{"route", "method"})
	QueueDepth        = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "zenx", Subsystem: "queue", Name: "depth", Help: "Queue depth"}, []string{"queue"})
	WSRoomConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "zenx", Subsystem: "websocket", Name: "room_connections", Help: "Websocket room connections"}, []string{"room"})
	BusinessCounter   = prometheus.NewCounterVec(prometheus.CounterOpts{Namespace: "zenx", Subsystem: "business", Name: "events_total", Help: "Business events"}, []string{"event"})
)

func Register(reg prometheus.Registerer) {
	reg.MustRegister(RouteLatency, QueueDepth, WSRoomConnections, BusinessCounter)
}
