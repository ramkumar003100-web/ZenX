package anomaly

import "math"

func ZScore(v, mean, std float64) float64 {
	if std == 0 {
		return 0
	}
	return (v - mean) / std
}
func IsOutlier(v, mean, std, threshold float64) bool {
	return math.Abs(ZScore(v, mean, std)) > threshold
}
