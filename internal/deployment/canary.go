package deployment

import "math"

type Canary struct {
	StableVersion, CanaryVersion string
	CanaryPercent                float64
}

func (c Canary) RouteByHash(hash uint64) string {
	bucket := float64(hash%100) / 100.0
	if bucket < math.Max(0, math.Min(1, c.CanaryPercent)) {
		return c.CanaryVersion
	}
	return c.StableVersion
}
