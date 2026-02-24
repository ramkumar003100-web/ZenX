package geo

func Failover(primary []Node, secondary []Node) []Node {
	for _, n := range primary {
		if n.Healthy {
			return primary
		}
	}
	return secondary
}
