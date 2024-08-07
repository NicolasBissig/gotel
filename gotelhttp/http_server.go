package gotelhttp

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"
)

func Handle(mux *http.ServeMux, pattern string, handlerFunc http.HandlerFunc) {
	route := extractRoute(pattern)
	// Configure the "http.route" for the HTTP instrumentation.
	withRouteTag := otelhttp.WithRouteTag(route, handlerFunc)
	withCorrectName := spanNameInjector(route, withRouteTag)
	mux.Handle(pattern, withCorrectName)
}

func spanNameInjector(route string, handlerFunc http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		span.SetName(r.Method + " " + route)
		handlerFunc.ServeHTTP(w, r)
	}
}

// extractRoute turns a pattern like "GET /rolldice" into "/rolldice".
func extractRoute(pattern string) string {
	return pattern[strings.Index(pattern, " ")+1:]
}
