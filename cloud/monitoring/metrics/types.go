package metrics

type Metric struct {
	Name   string
	Value  float64
	Labels map[string]string
}
