package collector

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type initialCollector struct {
	collectors []prometheus.Collector
	logger     log.Logger
}

var (
	factories = make(map[string]func(logger log.Logger) prometheus.Collector)
)

func registerCollector(collector string, factory func(logger log.Logger) prometheus.Collector) {
	factories[collector] = factory
}

func NewPatroniCollector(logger log.Logger) prometheus.Collector {
	var collectors []prometheus.Collector
	for key, factory := range factories {
		collector := factory(log.With(logger, "collector", key))
		collectors = append(collectors, collector)
	}
	return initialCollector{
		collectors: collectors,
		logger:     logger,
	}
}

// Describe implements the prometheus.Collector interface.
func (p *initialCollector) Describe(ch chan<- *prometheus.Desc) {
	wg := sync.WaitGroup{}
	wg.Add(len(p.collectors))
	for _, c := range p.collectors {
		go func(c prometheus.Collector) {
			c.Describe(ch)
			wg.Done()
		}(c)
	}
	wg.Wait()
}

// Collect implements the prometheus.Collector interface.
func (p *initialCollector) Collect(ch chan<- *prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(p.collectors))
	for _, c := range p.collectors {
		go func(c prometheus.Collector) {
			c.Collect(ch)
			wg.Done()
		}(c)
	}
	wg.Wait()
}
