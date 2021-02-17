package collector

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/gopaytech/patroni_exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type initialCollector struct {
	collectors []prometheus.Collector
	client     client.PatroniClient
	logger     log.Logger
}

const namespace = "patroni"

var (
	factories = make(map[string]func(client client.PatroniClient, logger log.Logger) prometheus.Collector)
)

func registerCollector(collector string, factory func(client client.PatroniClient, logger log.Logger) prometheus.Collector) {
	factories[collector] = factory
}

func NewPatroniCollector(client client.PatroniClient, logger log.Logger) prometheus.Collector {
	var collectors []prometheus.Collector
	for key, factory := range factories {
		collector := factory(client, log.With(logger, "collector", key))
		collectors = append(collectors, collector)
	}
	return &initialCollector{
		collectors: collectors,
		client:     client,
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
func (p *initialCollector) Collect(ch chan<- prometheus.Metric) {
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
