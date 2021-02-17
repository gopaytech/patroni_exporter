package collector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("patroni", createPatroniCollectorFactory)
}

var (
	possiblePatroniState = [...]string{"running", "stopped", "restarted"}
)

type patroniCollector struct {
	state  *prometheus.Desc
	logger log.Logger
}

func createPatroniCollectorFactory(logger log.Logger) prometheus.Collector {
	state := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "state"),
		"The current state of Patroni service",
		nil,
		nil)
	return &patroniCollector{
		state:  state,
		logger: logger,
	}
}

func (p *patroniCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.state
}

func (p *patroniCollector) Collect(ch chan<- prometheus.Metric) {
	// restyClient := resty.New()
	// resp, err := restyClient.R().EnableTrace().Get(fmt.Sprintf("%s:%s/patroni", opts.host, opts.port))
	// if err != nil {
	// fmt.Errorf(err)
	// }
}
