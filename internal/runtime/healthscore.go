package runtime

func HealthScore(component map[string]float64, goroutines int64, deadlocked bool) float64 {
	if deadlocked {
		return 0
	}
	if len(component) == 0 {
		return 1
	}
	total := 0.0
	for _, v := range component {
		total += clamp(v)
	}
	score := total / float64(len(component))
	if goroutines > 5000 {
		score *= 0.7
	}
	return clamp(score)
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
