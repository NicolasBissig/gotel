package goteloapi

import (
	"context"
	"fmt"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

// GotelMiddlewares just registers the GotelMiddleware
var GotelMiddlewares = []strictnethttp.StrictHTTPMiddlewareFunc{
	GotelMiddleware,
}

// GotelMiddleware is a middleware that sets the span name to the operation ID provided from the OpenAPI spec
func GotelMiddleware(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		span := trace.SpanFromContext(ctx)
		// set the operation name to the operation ID
		span.SetName(fmt.Sprintf("%s %s", r.Method, operationID))
		span.SetAttributes(attribute.String("openapi.operation_id", operationID))

		// TODO: check for error and change status + add the error to the span
		return f(ctx, w, r, request)
	}
}
