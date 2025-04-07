package grpc

import (
	"context"
	"fmt"
	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

type TraceAnalyzerServer struct {
	tracev1.UnimplementedTraceServiceServer
}

func (s *TraceAnalyzerServer) Export(ctx context.Context, request *tracev1.ExportTraceServiceRequest) (*tracev1.ExportTraceServiceResponse, error) {
	for _, each := range request.ResourceSpans {
		fmt.Println(each.String())
	}
	return &tracev1.ExportTraceServiceResponse{}, nil
}
