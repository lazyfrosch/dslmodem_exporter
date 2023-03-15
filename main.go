package main

import (
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log/level"
	"github.com/lazyfrosch/dslmodem_exporter/pkg/zyxel"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const (
	namespace = "dslmodem"
)

func main() {
	var (
		webConfig   = webflag.AddFlags(kingpin.CommandLine, ":9120")
		metricsPath = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()

		url = kingpin.Flag("dslmodem.url", "URL to the DSL modem. (env:DSLEXPORTER_MODEM_URL)").
			Default("http://192.168.1.1").Envar("DSLEXPORTER_MODEM_URL").String()
		username = kingpin.Flag("dslmodem.username", "Username for authenticating to the modem. (env:DSLEXPORTER_MODEM_USERNAME)").
				Default("admin").Envar("DSLEXPORTER_MODEM_USERNAME").String()
		password = kingpin.Flag("dslmodem.password", "Password for authenticating to the modem. (env:DSLEXPORTER_MODEM_PASSWORD)").
				Default("1234").Envar("DSLEXPORTER_MODEM_PASSWORD").String()
		interval = kingpin.Flag("dslmodem.interval", "Interval to pull data from the modem. (env:DSLEXPORTER_MODEM_INTERVAL)").
				Default("15s").Envar("DSLEXPORTER_MODEM_INTERVAL").Duration()

		readTimeout = kingpin.Flag("http.read-timeout", "Timeout for reading from HTTP sockets").
				Default("5s").Envar("HTTP_READ_TIMEOUT").Duration()
		writeTimeout = kingpin.Flag("http.write-timeout", "Timeout for reading from HTTP sockets").
				Default("5s").Envar("HTTP_WRITE_TIMEOUT").Duration()
		idleTimeout = kingpin.Flag("http.idle-timeout", "Timeout for reading from HTTP sockets").
				Default("120s").Envar("HTTP_IDLE_TIMEOUT").Duration()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("dslmodem_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)

	_ = level.Info(logger).Log("msg", "Starting dslmodem_exporter", "version", version.Info())
	_ = level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	zyxel.Logger = logger
	client, err := zyxel.NewClient(*url, *username, *password)
	if err != nil {
		_ = level.Error(logger).Log("msg", "Error creating the client", "err", err)
		os.Exit(1)
	}

	go updateLoop(client, *interval, logger)

	prometheus.MustRegister(version.NewCollector("dslmodem_exporter"))

	http.HandleFunc("/", homepage)
	http.Handle(*metricsPath, promhttp.Handler())

	srv := &http.Server{
		ReadTimeout:       *readTimeout,
		ReadHeaderTimeout: *readTimeout,
		WriteTimeout:      *writeTimeout,
		IdleTimeout:       *idleTimeout,
	}

	if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
		_ = level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}

func homepage(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(`<html>
	<head><title>DSL Modem Exporter</title></head>
	<body>
	<h1>DSL Modem Exporter</h1>
	<p><a href="/metrics">Metrics</a></p>
	</body>
	</html>`))
}
