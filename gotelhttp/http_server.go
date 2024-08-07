package gotelhttp

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"
)

func HandleFunc(pattern string, handlerFunc http.HandlerFunc, mux ...*http.ServeMux) {
	route := extractRoute(pattern)
	// Configure the "http.route" for the HTTP instrumentation.
	withRouteTag := otelhttp.WithRouteTag(route, handlerFunc)
	withCorrectName := spanNameInjector(route, withRouteTag)
	wrapped := otelhttp.NewHandler(withCorrectName, route)

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
