package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	log     = logrus.New()
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

func newRootCommand() *cobra.Command {
	var (
		configFile   string
		port         int
		megacliPath  string
		logLevel     string
		timeout      int
	)

	cmd := &cobra.Command{
		Use:   "megaraid-exporter",
		Short: "Prometheus exporter for MegaRAID controllers",
		Long:  "A Prometheus exporter that collects metrics from MegaRAID controllers using MegaCLI64",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(configFile, port, megacliPath, logLevel, timeout)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	cmd.Flags().IntVarP(&port, "port", "p", 9272, "HTTP port to listen on")
	cmd.Flags().StringVar(&megacliPath, "megacli-path", "/usr/sbin/megacli64", "Path to megacli64 binary")
	cmd.Flags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	cmd.Flags().IntVar(&timeout, "timeout", 30, "Command timeout in seconds")

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("megaraid-exporter version %s\n", version)
		},
	})

	return cmd
}

func run(configFile string, port int, megacliPath, logLevel string, timeout int) error {
	// Setup logging
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}
	log.SetLevel(level)
	log.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Warnf("Failed to read config file: %v", err)
		}
	}

	// Override with command line flags
	viper.SetDefault("port", port)
	viper.SetDefault("megacli_path", megacliPath)
	viper.SetDefault("command_timeout", fmt.Sprintf("%ds", timeout))

	log.WithFields(logrus.Fields{
		"version":      version,
		"port":         viper.GetInt("port"),
		"megacli_path": viper.GetString("megacli_path"),
		"timeout":      viper.GetString("command_timeout"),
	}).Info("Starting MegaRAID exporter")

	// Verify megacli64 is available
	if err := verifyMegaCLI(viper.GetString("megacli_path")); err != nil {
		return fmt.Errorf("MegaCLI verification failed: %v", err)
	}

	// Create and register collector
	collector := NewMegaRAIDCollector(
		viper.GetString("megacli_path"),
		viper.GetDuration("command_timeout"),
		log,
	)
	
	prometheus.MustRegister(collector)
	
	// Setup HTTP server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html>
<head><title>MegaRAID Exporter</title></head>
<body>
<h1>MegaRAID Exporter</h1>
<p><a href="/metrics">Metrics</a></p>
<p><a href="/health">Health</a></p>
<p>Version: %s</p>
</body>
</html>`, version)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Info("Received shutdown signal")
		
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Errorf("HTTP server shutdown error: %v", err)
		}
		cancel()
	}()

	log.Infof("HTTP server listening on :%d", viper.GetInt("port"))
	
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %v", err)
	}

	<-ctx.Done()
	log.Info("Exporter stopped")
	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","version":"%s"}`, version)
}

func verifyMegaCLI(path string) error {
	// Implementation would verify megacli64 is executable and working
	log.Infof("Verifying MegaCLI at: %s", path)
	return nil
}

// Placeholder collector - actual implementation would be in separate files
type MegaRAIDCollector struct {
	megacliPath string
	timeout     time.Duration
	logger      *logrus.Logger
}

func NewMegaRAIDCollector(path string, timeout time.Duration, logger *logrus.Logger) *MegaRAIDCollector {
	return &MegaRAIDCollector{
		megacliPath: path,
		timeout:     timeout,
		logger:      logger,
	}
}

func (c *MegaRAIDCollector) Describe(ch chan<- *prometheus.Desc) {
	// Implementation needed
}

func (c *MegaRAIDCollector) Collect(ch chan<- prometheus.Metric) {
	// Implementation needed
}
