package prom

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tsenart/vegeta/v12/lib/prom"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/tests/internal/cfg"
)

const operation = "prom server"

const (
	idleTimeout        = 240 * time.Second
	writeHeaderTimeout = 15 * time.Second
	readHeaderTimeout  = 30 * time.Second
)

func NewServer(cfg cfg.Config) (*http.Server, *prom.Metrics, error) {
	pm := prom.NewMetrics()
	r := prometheus.NewRegistry()
	if err := pm.Register(r); err != nil {
		return nil, nil, fmt.Errorf("%s: failed to register prom registry: %w", operation, err)
	}

	srv := &http.Server{
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeHeaderTimeout,
		IdleTimeout:       idleTimeout,
		Addr:              cfg.PromExporterAddr,
		Handler:           prom.NewHandler(r, time.Now().UTC()),
	}

	return srv, pm, nil
}
