package main

import "github.com/WaffeSoul/metrics-collector/internal/agent"

func main() {
	collect := agent.NewCollector("http://localhost:8080", 2, 10)
	go collect.UpdateMetrict()
	collect.UpdateMetricToServer()
}
