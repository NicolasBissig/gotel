package gotelhttp

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"net/http"
)

func NewClient() http.Client {
	return http.Client{
		Transport: otelhttp.NewTransport(nil),
	}
}

func NewRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

func InstrumentDefaultClient() {
	oldtransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = otelhttp.NewTransport(oldtransport)
}

func Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
