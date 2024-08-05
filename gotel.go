package gotel

import (
	"context"
	"github.com/NicolasBissig/gotel/lib"
	"os"
	"os/signal"
)

// OptionFunc is a function that sets some option on the OpenTelemetry SDK.
type OptionFunc func() error

type InitializedOtelSdk struct {
	Shutdown       func(context.Context) error
	CancelFunction context.CancelFunc
	RootContext    context.Context
}

func Setup(opts ...OptionFunc) (*InitializedOtelSdk, error) {
	result := &InitializedOtelSdk{}

	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	result.RootContext = ctx

	// Set up OpenTelemetry.
	otelShutdown, err := lib.SetupOTelSDK(ctx)
	if err != nil {
		return nil, err
	}

	result.Shutdown = otelShutdown

	return result, nil
}
