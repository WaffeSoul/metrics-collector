package main

import "github.com/WaffeSoul/metrics-collector/internal/agent"

func main() {
	parseFlags()
	collect := agent.NewCollector(cfg.Addr, cfg.PollInterval, cfg.ReportInterval, cfg.KeyHash)
	go collect.UpdateMetrict()
	collect.UpdateMetricToServer()
}
