package main

import (
	"sync"

	"github.com/WaffeSoul/metrics-collector/internal/agent"
)

func main() {
	parseFlags()
	collect := agent.NewCollector(cfg.Addr, cfg.PollInterval, cfg.ReportInterval, cfg.KeyHash, cfg.Rate)
	fieldsCh := make(chan agent.Fields)
	var wg sync.WaitGroup
	go collect.UpdateMetrict(fieldsCh)
	go collect.UpdataGopsutil(fieldsCh)

	for w := int64(0); w < collect.Rate; w++ {
		wg.Add(1)
		go collect.UpdateMetricToServer(fieldsCh)
	}
	wg.Wait()
}
