package internal

import (
	"os"
	"strings"
)

const (
	OTEL_EXPORTER_OTLP_ENDPOINT = "OTEL_EXPORTER_OTLP_ENDPOINT"
)

func exportEndpoint() (string, error) {
	// The default OTLP exporter is set to "https://localhost:4317"
	// See: https://github.com/open-telemetry/opentelemetry-go/issues/4147
	// usually http://localhost:4317 is preferred
	_, present := os.LookupEnv(OTEL_EXPORTER_OTLP_ENDPOINT)
	if !present {
		err := os.Setenv(OTEL_EXPORTER_OTLP_ENDPOINT, "http://localhost:4317")
		if err != nil {
			return "", err
		}
	}

	return os.Getenv(OTEL_EXPORTER_OTLP_ENDPOINT), nil
}

type Protocol string

const (
	ProtocolGRPC Protocol = "grpc"
	ProtocolHTTP Protocol = "http"
)

func lookupProtocol() Protocol {
	envvar := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	if envvar == "" {
		return ProtocolHTTP
	}
	if strings.Contains(envvar, "grpc") {
		return ProtocolGRPC
	}
	return ProtocolHTTP
}
