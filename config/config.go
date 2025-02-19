package config

import (
	"fmt"

	"github.com/faizal-asep-outlook/env"
)

// Config holds the configuration for the telemetry.
type Config struct {
	// App configuration
	OtlpEndpoint   string `env:"OTEL_TRACING_OTLP_ENDPOINT" default:""`
	ServiceName    string `env:"OTEL_TRACING_SERVICE_NAME" default:"service"`
	ServiceVersion string `env:"OTEL_TRACING_SERVICE_VERSION" default:"1.0.0"`
	Insecure       bool   `env:"OTEL_TRACING_INSECURE_MODE" default:"true"`
}

// NewConfigFromEnv creates a new telemetry config from the environment.
func NewConfigFromEnv() (Config, error) {

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse telemetry config: %w", err)
	}

	return cfg, nil
}
