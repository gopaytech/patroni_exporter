package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/gopaytech/patroni_exporter/client"
	"github.com/gopaytech/patroni_exporter/collector"
	options "github.com/gopaytech/patroni_exporter/opts"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9933").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	opts          = options.PatroniOpts{}
)

func main() {
	kingpin.Flag("patroni.host", "Patroni host or IP Address").Default("http://localhost").StringVar(&opts.Host)
	kingpin.Flag("patroni.port", "Patroni port").Default("8008").StringVar(&opts.Port)

	promlogConfig := &promlog.Config{}
	logger := promlog.New(promlogConfig)
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	level.Info(logger).Log("msg", "Starting patroni_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	httpClient := resty.New()
	patroniClient := client.NewPatroniClient(httpClient, opts)
	patroniCollector := collector.NewPatroniCollector(patroniClient, logger)
	prometheus.MustRegister(patroniCollector)
	prometheus.MustRegister(version.NewCollector("patroni_exporter"))

	setupHandler()
	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}

func setupHandler() {
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
}
