package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

// Config holds the configuration for the telemetry.
type Config struct {
	// App configuration
	OtlpEndpoint   string `env:"OTEL_TRACING_OTLP_ENDPOINT" envDefault:""`
	ServiceName    string `env:"OTEL_TRACING_SERVICE_NAME" envDefault:"service"`
	ServiceVersion string `env:"OTEL_TRACING_SERVICE_VERSION" envDefault:"1.0.0"`
	Insecure       bool   `env:"OTEL_TRACING_INSECURE_MODE" envDefault:"true"`
}

// NewConfigFromEnv creates a new telemetry config from the environment.
func NewConfigFromEnv() (Config, error) {

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse telemetry config: %w", err)
	}

	return cfg, nil
}
