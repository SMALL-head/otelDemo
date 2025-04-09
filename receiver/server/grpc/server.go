package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
	"net"
	"otelDemo/utils/otelutils"
)

type TraceAnalyzerServer struct {
	tracev1.UnimplementedTraceServiceServer
}

type TraceAnalyzerGrpcServer struct {
	Addr   string
	Server *grpc.Server
}

func (s *TraceAnalyzerGrpcServer) Run() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		logrus.Fatalf("Failed to listen: %v", err)
	}
	logrus.Infof("[Run] - 监听地址: %s", s.Addr)
	if err := s.Server.Serve(lis); err != nil {
		logrus.Fatalf("Failed to serve: %v", err)
	}
	return nil
}

func (s *TraceAnalyzerGrpcServer) Close() {
	s.Server.GracefulStop()
	logrus.Infof("[TraceAnalyzerGrpcServer] - 关闭成功")
}

func (s *TraceAnalyzerServer) Export(ctx context.Context, request *tracev1.ExportTraceServiceRequest) (*tracev1.ExportTraceServiceResponse, error) {
	logrus.Infof("[export] - 被调用")
	for _, each := range request.ResourceSpans {
		for _, span := range each.ScopeSpans {
			for _, eachSpan := range span.Spans {
				logrus.Infof("[export] - spanid = %s", otelutils.OtelIDToString(eachSpan.SpanId))
			}
		}
	}

	return &tracev1.ExportTraceServiceResponse{}, nil
}

func NewTraceAnalyzerServer(addr string) *TraceAnalyzerGrpcServer {
	server := grpc.NewServer()
	tracev1.RegisterTraceServiceServer(server, &TraceAnalyzerServer{})
	return &TraceAnalyzerGrpcServer{
		Addr:   addr,
		Server: server,
	}
}
