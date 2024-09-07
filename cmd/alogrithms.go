package cmd

type Algorithm string

const (
	LeastTime          Algorithm = "least-time"
	WeightedRoundRobin Algorithm = "weighted-round-robin"
	ConnectionPerTime  Algorithm = "connection-per-time"
	RoundRobin         Algorithm = "round-robin"
)
