package gotel

import (
	"context"
	"fmt"
	"github.com/NicolasBissig/gotel/internal"
	"os"
	"os/signal"
)

// OptionFunc is a function that sets some option on the OpenTelemetry SDK.
type OptionFunc func() error

type Sdk struct {
	Shutdown    func(context.Context) error
	RootContext context.Context
}

// Setup initializes the OpenTelemetry SDK.
// Calling Shutdown is not necessary, as it is called automatically when the RootContext is done, i.e. when the application is terminated.
//
// Example usage:
//
//	sdk, err := gotel.Setup()
//	if err != nil {
//	    return fmt.Errorf("failed to set up OpenTelemetry: %w", err)
//	}
func Setup(opts ...OptionFunc) (*Sdk, error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	// Set up OpenTelemetry.
	otelShutdown, err := internal.SetupOTelSDK(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to set up OpenTelemetry: %w", err)
	}

	sdk := &Sdk{
		Shutdown:    otelShutdown,
		RootContext: ctx,
	}

	go func() {
		<-sdk.RootContext.Done()
		_ = sdk.Shutdown(context.Background())
	}()

	return sdk, nil
}
