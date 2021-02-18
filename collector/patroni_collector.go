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
	possiblePatroniRole  = [...]string{"MASTER", "REPLICA"}
)

type patroniCollector struct {
	stateDesc  *prometheus.Desc
	roleDesc   *prometheus.Desc
	staticDesc *prometheus.Desc
	logger     log.Logger
	client     client.PatroniClient
}

func createPatroniCollectorFactory(client client.PatroniClient, logger log.Logger) prometheus.Collector {
	stateDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "state"),
		"The current state of Patroni service",
		[]string{"state", "scope"},
		nil)
	roleDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "role"),
		"The current PostgreSQL role of Patroni node",
		[]string{"role", "scope"},
		nil)
	staticDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "cluster_node", "static"),
		"The collection of static value as reported by Patroni",
		[]string{"version"},
		nil)
	return &patroniCollector{
		stateDesc:  stateDesc,
		roleDesc:   roleDesc,
		staticDesc: staticDesc,
		logger:     logger,
		client:     client,
	}
}

func (p *patroniCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.stateDesc
	ch <- p.roleDesc
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
		ch <- prometheus.MustNewConstMetric(p.stateDesc, prometheus.GaugeValue, stateValue, possibleState, patroniResponse.Patroni.Scope)
	}

	for _, possibleRole := range possiblePatroniRole {
		stateValue := 0.0
		if strings.ToUpper(patroniResponse.Role) == possibleRole {
			stateValue = 1.0
		}
		ch <- prometheus.MustNewConstMetric(p.roleDesc, prometheus.GaugeValue, stateValue, possibleRole, patroniResponse.Patroni.Scope)
	}
	ch <- prometheus.MustNewConstMetric(p.staticDesc, prometheus.GaugeValue, 1.0, patroniResponse.Patroni.Version)
}
