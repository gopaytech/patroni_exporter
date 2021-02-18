package collector

import (
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/gopaytech/patroni_exporter/client"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("patroni", createPatroniCollectorFactory)
}

var (
	possiblePatroniState = [...]string{"RUNNING", "STOPPED", "PROMOTED", "UNKNOWN"}
)

type patroniCollector struct {
	state  *prometheus.Desc
	logger log.Logger
	client client.PatroniClient
}

func createPatroniCollectorFactory(client client.PatroniClient, logger log.Logger) prometheus.Collector {
	state := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "state"),
		"The current state of Patroni service",
		[]string{"state"},
		nil)
	return &patroniCollector{
		state:  state,
		logger: logger,
		client: client,
	}
}

func (p *patroniCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.state
}

func (p *patroniCollector) Collect(ch chan<- prometheus.Metric) {
	patroniResponse, err := p.client.GetMetrics()
	if err != nil {
		level.Error(p.logger).Log("msg", "Unable to get metrics from Patroni", "err", fmt.Sprintf("errornya: %v", err))
		return
	}
	for _, possibleState := range possiblePatroniState {
		stateValue := 0.0
		if strings.ToUpper(patroniResponse.State) == possibleState {
			stateValue = 1.0
		}
		ch <- prometheus.MustNewConstMetric(p.state, prometheus.GaugeValue, stateValue, possibleState)
	}
}
