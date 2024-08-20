package gotelhttp

import (
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"
)

// NewHandler wraps a http.Handler with OpenTelemetry instrumentation.
// The span name is set to the HTTP method only as no route is available.
func NewHandler(handler http.Handler) http.Handler {
	return otelhttp.NewHandler(traceContextInjector(spanName(handler)), "Unknown operation")
}

func HandleFunc(pattern string, handlerFunc http.HandlerFunc, mux ...*http.ServeMux) {
	route := extractRoute(pattern)
	// Configure the "http.route" for the HTTP instrumentation.
	withRouteTag := otelhttp.WithRouteTag(route, handlerFunc)
	withCorrectName := routeAsSpanName(route, withRouteTag)
	wrapped := otelhttp.NewHandler(traceContextInjector(withCorrectName), route)

	if len(mux) == 0 {
		http.Handle(pattern, wrapped)
	} else {
		for _, m := range mux {
			m.Handle(pattern, wrapped)
		}
	}
}

func Handle(pattern string, handler http.Handler, mux ...*http.ServeMux) {
	HandleFunc(pattern, handler.ServeHTTP, mux...)
}

type ServeMux struct {
	*http.ServeMux
}

func NewServeMux() *ServeMux {
	mux := http.NewServeMux()

	return &ServeMux{
		ServeMux: mux,
	}
}

func (mux *ServeMux) Handle(pattern string, handler http.Handler) {
	HandleFunc(pattern, handler.ServeHTTP, mux.ServeMux)
}

func (mux *ServeMux) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	HandleFunc(pattern, handlerFunc, mux.ServeMux)
}

// routeAsSpanName sets the span name according to https://opentelemetry.io/docs/specs/semconv/http/http-spans/#name for a specific route
func routeAsSpanName(route string, handlerFunc http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		span.SetName(r.Method + " " + route)
		handlerFunc.ServeHTTP(w, r)
	}
}

// spanName sets the span name according to https://opentelemetry.io/docs/specs/semconv/http/http-spans/#name but without a route available
func spanName(handlerFunc http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())

		method := r.Method
		if method == "" {
			method = "HTTP"
		}

		span.SetName(method)

		handlerFunc.ServeHTTP(w, r)
	}
}

// extractRoute turns a pattern like "GET /rolldice" into "/rolldice".
func extractRoute(pattern string) string {
	return pattern[strings.Index(pattern, " ")+1:]
}

func traceContextInjector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if w.Header().Get("traceparent") != "" {
			next.ServeHTTP(w, r)
			return
		}

		// get the current span
		span := trace.SpanFromContext(r.Context())
		// set the traceparent in w3c format: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
		traceparent := fmt.Sprintf("00-%s-%s-%s", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String(), span.SpanContext().TraceFlags().String())
		w.Header().Set("traceparent", traceparent)

		// set the tracestate, if present
		tracestate := span.SpanContext().TraceState()
		if tracestate.Len() > 0 {
			w.Header().Set("tracestate", tracestate.String())
		}

		next.ServeHTTP(w, r)
	})
}
