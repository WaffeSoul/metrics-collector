package main

import "github.com/WaffeSoul/metrics-collector/internal/agent"

func main() {
	parseFlags()
	collect := agent.NewCollector(addr, pollInterval, reportInterval)
	go collect.UpdateMetrict()
	collect.UpdateMetricToServer()
}
