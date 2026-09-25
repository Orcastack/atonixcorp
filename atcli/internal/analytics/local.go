package analytics

import "sync"

var mu sync.Mutex
var LocalMetrics = map[string]int{}

func Inc(metric string) {
	mu.Lock()
	LocalMetrics[metric]++
	mu.Unlock()
}

func FlushLocalMetrics() {
	SendEvent("cli-metrics", LocalMetrics)
}
