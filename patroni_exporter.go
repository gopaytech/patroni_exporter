package patroni_exporter

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/gopaytech/patroni_exporter/client"
	"github.com/gopaytech/patroni_exporter/collector"
	options "github.com/gopaytech/patroni_exporter/opts"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9933").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	opts          = options.PatroniOpts{}
)

func main() {
	kingpin.Flag("patroni.host", "Patroni host or IP Address").Default("localhost").StringVar(&opts.Host)
	kingpin.Flag("patroni.port", "Patroni port").Default("8008").StringVar(&opts.Port)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting patroni_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	patroniClient := client.NewPatroniClient(opts)
	// if err != nil {
	// level.Error(logger).Log("msg", "Error initialize patroni_exporter", "err", err)
	// os.Exit(1)
	// }

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
