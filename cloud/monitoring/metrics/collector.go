package metrics

type Collector interface {
	Collect() ([]Metric, error)
}
