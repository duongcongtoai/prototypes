package opentracing

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	zipkintracer "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

const endpointUrl = "http://localhost:9411/api/v2/spans"

func NewEchoZipkinMiddleWare(serviceName, spanName string) (echo.MiddlewareFunc, opentracing.Tracer, error) {
	newtracer, err := NewTracer(serviceName)

	if err != nil {
		return nil, nil, err
	}
	opentracingTracer := zipkintracer.Wrap(newtracer)

	return echo.WrapMiddleware(createStandardMiddlewareFunc(opentracingTracer)), opentracingTracer, err
	// return echo.WrapMiddleware(zipkinhttp.NewServerMiddleware(
	// 	newtracer,
	// 	zipkinhttp.SpanName(spanName),
	// )), newtracer, nil
}

func NewTracer(serviceName string) (*zipkin.Tracer, error) {
	reporter := reporterhttp.NewReporter(endpointUrl)
	localEndpoint := &model.Endpoint{ServiceName: serviceName}
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}
	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
		zipkin.WithSharedSpans(false),
	)
	return t, err
}

func createStandardMiddlewareFunc(t opentracing.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return nethttp.Middleware(t, next)
	}
}
