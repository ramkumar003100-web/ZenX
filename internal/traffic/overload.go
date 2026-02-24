package traffic

import "sync/atomic"

type OverloadDetector struct {
	thresholdCPU float64
	thresholdQPS int64
	overloaded   atomic.Bool
}

func NewOverloadDetector(cpu float64, qps int64) *OverloadDetector {
	return &OverloadDetector{thresholdCPU: cpu, thresholdQPS: qps}
}

func (o *OverloadDetector) Evaluate(cpu float64, qps int64) bool {
	over := cpu > o.thresholdCPU || qps > o.thresholdQPS
	o.overloaded.Store(over)
	return over
}

func (o *OverloadDetector) IsOverloaded() bool { return o.overloaded.Load() }
