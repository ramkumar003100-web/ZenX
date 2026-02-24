package anomaly

import "time"

type Finding struct {
	Name  string
	Score float64
	Time  time.Time
}

type Detector struct {
	TrafficMean float64
	TrafficStd  float64
	LoginMean   float64
	LoginStd    float64
}

func (d Detector) DetectTrafficSpike(v float64) (Finding, bool) {
	s := ZScore(v, d.TrafficMean, d.TrafficStd)
	if s > 3 {
		return Finding{Name: "traffic_spike", Score: s, Time: time.Now()}, true
	}
	return Finding{}, false
}

func (d Detector) DetectUnusualLogin(v float64) (Finding, bool) {
	s := ZScore(v, d.LoginMean, d.LoginStd)
	if s > 2.5 {
		return Finding{Name: "unusual_login", Score: s, Time: time.Now()}, true
	}
	return Finding{}, false
}
