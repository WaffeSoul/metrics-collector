package main

import "github.com/WaffeSoul/metrics-collector/internal/agent"

func main() {
	parseFlags()
	collect := agent.NewCollector(cfg.addr, cfg.pollInterval, cfg.reportInterval)
	go collect.UpdateMetrict()
	collect.UpdateMetricToServer()
}
