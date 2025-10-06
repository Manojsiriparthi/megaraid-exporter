package main

import (
	"flag"
	"log"
	"net/http"
	
	"github.com/manojsiriparthi/megaraid-exporter/config"
	"github.com/manojsiriparthi/megaraid-exporter/pkg/megacli"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":8080", "Address to listen on for web interface and telemetry.")
		megacliPath   = flag.String("megacli.path", "", "Path to MegaCli64 binary (auto-discovered if not specified)")
	)
	flag.Parse()

	// Initialize configuration
	cfg := config.NewConfig()
	cfg.Port = *listenAddress
	
	if *megacliPath != "" {
		cfg.SetMegaCLIPath(*megacliPath)
	}

	// Initialize path monitor
	pathMonitor := megacli.NewPathMonitor(cfg)
	pathMonitor.Start()

	// Wait a moment for initial path check
	time.Sleep(1 * time.Second)

	if err := pathMonitor.ValidateOrError(); err != nil {
		log.Fatalf("MegaCLI validation failed: %v", err)
	}

	log.Printf("Starting MegaRAID exporter on %s", *listenAddress)
	log.Printf("Using MegaCLI path: %s", pathMonitor.GetPath())

	// Setup HTTP handlers
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := pathMonitor.GetStatus()
		if pathMonitor.IsHealthy() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
