package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

// Config holds the configuration for the telemetry.
type Config struct {
	// App configuration
	OtlpEndpoint   string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" envDefault:"127.0.0.1:4317"`
	ServiceName    string `env:"SERVICE_NAME" envDefault:"komposer"`
	ServiceVersion string `env:"SERVICE_VERSION" envDefault:"1.0.0"`
	Enabled        bool   `env:"TELEMETRY_ENABLED" envDefault:"true"`
	Insecure       bool   `env:"INSECURE_MODE" envDefault:"true"`
}

// NewConfigFromEnv creates a new telemetry config from the environment.
func NewConfigFromEnv() (Config, error) {

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse telemetry config: %w", err)
	}

	return cfg, nil
}
