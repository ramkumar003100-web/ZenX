package intelligence

import "time"

type Snapshot struct {
	RouteLatency map[string]time.Duration
	QueryLatency map[string]time.Duration
	JobFailures  map[string]int64
	WSBacklog    map[string]int64
}

type Insight struct {
	Bottlenecks  []string
	SlowQueries  []string
	CongestedWS  []string
	FailureTrend []string
}

func Analyze(s Snapshot) Insight {
	ins := Insight{}
	for r, l := range s.RouteLatency {
		if l > 700*time.Millisecond {
			ins.Bottlenecks = append(ins.Bottlenecks, r)
		}
	}
	for q, l := range s.QueryLatency {
		if l > 300*time.Millisecond {
			ins.SlowQueries = append(ins.SlowQueries, q)
		}
	}
	for j, f := range s.JobFailures {
		if f > 10 {
			ins.FailureTrend = append(ins.FailureTrend, j)
		}
	}
	for room, b := range s.WSBacklog {
		if b > 100 {
			ins.CongestedWS = append(ins.CongestedWS, room)
		}
	}
	return ins
}
