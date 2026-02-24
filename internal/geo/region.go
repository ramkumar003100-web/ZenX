package geo

import "sort"

type Node struct {
	ID, Region string
	LatencyMS  float64
	Healthy    bool
}

type Router struct{ nodes []Node }

func (r *Router) SetNodes(nodes []Node) { r.nodes = nodes }
func (r *Router) ByRegion(region string) []Node {
	out := []Node{}
	for _, n := range r.nodes {
		if n.Region == region && n.Healthy {
			out = append(out, n)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].LatencyMS < out[j].LatencyMS })
	return out
}
