package goteloapi

import (
	"context"
	"fmt"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func GotelMiddleware(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		span := trace.SpanFromContext(ctx)
		// set the operation name to the operation ID
		span.SetName(fmt.Sprintf("%s %s", r.Method, operationID))
		return f(ctx, w, r, request)
	}
}
