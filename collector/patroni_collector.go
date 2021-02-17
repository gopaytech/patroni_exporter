package collector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("patroni", createPatroniCollectorFactory)
}

type patroniCollector struct {
	logger log.Logger
}

func createPatroniCollectorFactory(logger log.Logger) prometheus.Collector {
	// need to be exporterd
	// state
	// scope
	// role
	// patroni version
	return &patroniCollector{
		logger: logger,
	}
}

func (p *patroniCollector) Describe(ch chan<- *prometheus.Desc) {

}

func (p *patroniCollector) Collect(ch chan<- prometheus.Metric) {

}
