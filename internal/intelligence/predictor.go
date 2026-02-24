package intelligence

import "time"

type Signal struct {
	ScaleOut bool
	Reason   string
	At       time.Time
}

func Predict(ins Insight) Signal {
	if len(ins.Bottlenecks)+len(ins.CongestedWS)+len(ins.FailureTrend) > 2 {
		return Signal{ScaleOut: true, Reason: "combined pressure indicators", At: time.Now()}
	}
	return Signal{ScaleOut: false, Reason: "stable", At: time.Now()}
}
