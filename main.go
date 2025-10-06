package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress = flag.String("web.listen-address", ":9216", "Address to listen on for web interface and telemetry.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	storCliPath   = flag.String("storcli.path", "/usr/sbin/storcli64", "Path to storcli binary.")
	interval      = flag.Duration("interval", 30*time.Second, "Interval between metric collections.")
)

func main() {
	flag.Parse()

	collector := NewMegaRAIDCollector(*storCliPath)
	prometheus.MustRegister(collector)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>MegaRAID Exporter</title></head>
			<body>
			<h1>MegaRAID Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Starting MegaRAID exporter on %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
