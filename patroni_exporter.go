package patroni_exporter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/go-resty/resty/v2"
	"github.com/gopaytech/patroni_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type patroniOpts struct {
	host string
	port string
}

type patroniField struct {
	Scope   string `json:"scope"`
	Version string `json:"version"`
}

type patroniJSONResp struct {
	State   string       `json:"state"`
	Role    string       `json:"role"`
	Patroni patroniField `json:"patroni"`
}

type PatroniClient interface {
	GetMetrics() error
}

type patroniClient struct {
	resty resty.Client
	host  string
	port  string
}

func newPatroniClient(opts patroniOpts) PatroniClient {
	r := resty.New()
	r.SetHostURL(fmt.Sprintf("%s:%s", opts.host, opts.host))
	return &patroniClient{
		resty: r,
	}
}

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9933").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	opts          = patroniOpts{}
)

func main() {
	kingpin.Flag("patroni.host", "Patroni host or IP Address").Default("localhost").StringVar(&opts.host)
	kingpin.Flag("patroni.port", "Patroni port").Default("8008").StringVar(&opts.port)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting patroni_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	patroniClient := newPatroniClient(opts)
	if err != nil {
		level.Error(logger).Log("msg", "Error initialize patroni_exporter", "err", err)
		os.Exit(1)
	}

	collector := collector.NewPatroniCollector(patroniClient, logger)
	prometheus.MustRegister(collector)
	prometheus.MustRegister(version.NewCollector("patroni_exporter"))

	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>Patroni Exporter</title></head>
		<body>
		<h1>Patroni Exporter</h1>
		<p><a href="` + *metricsPath + `"></p>
		</body>
		</html>`))
	})
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}

func (p *patroniClient) GetMetrics() (patroniJSONResp, error) {
	resp, err := p.resty.R().EnableTrace().Get("/patroni")
	if err != nil {
		return patroniJSONResp{}, err
	}

	var objmap patroniJSONResp

	err := json.Unmarshal(resp.Body, &objmap)
	if err != nil {
		return patroniJSONResp{}, err
	}

	return objmap, nil
}
