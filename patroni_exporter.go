package patroni_exporter

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/go-resty/resty/v2"
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

func newPatroniExporter(opts patroniOpts) error {
	return nil
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

	patroni := newPatroniExporter(opts)
	if err != nil {
		level.Error(logger).Log("msg", "Error initialize patroni_exporter", "err", err)
		os.Exit(1)
	}

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

	restyClient := resty.New()
	resp, err := restyClient.R().EnableTrace().Get(fmt.Sprintf("%s:%s/patroni", opts.host, opts.port))
	if err != nil {
		fmt.Errorf(err)
	}
}
