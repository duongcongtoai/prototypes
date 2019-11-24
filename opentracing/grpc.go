package opentracing

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"

	"google.golang.org/grpc"
)

func TracedGrpcConn(t opentracing.Tracer, port int) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(
				t,
				otgrpc.LogPayloads(),
			),
		),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
