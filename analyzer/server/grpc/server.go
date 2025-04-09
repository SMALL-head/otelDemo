package grpc

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"net"
	"otelDemo/analyzer/common/otelmodel"
	"otelDemo/analyzer/config"
	"otelDemo/utils/httpc"
	"otelDemo/utils/otelutils"
	"otelDemo/utils/otelutils/openapi"
	"sync"
	"time"
)

type TraceAnalyzerServer struct {
	tracev1.UnimplementedTraceServiceServer

	cfg *config.AnalyzerConfig

	traceId      map[string]bool
	traceIdMutex sync.RWMutex

	analyseCh chan string // 此通道内存放需要分析的traceID，我们并不希望立即分析，希望间隔一小段时间后分析
}

type TraceAnalyzerGrpcServer struct {
	Addr   string
	Server *grpc.Server
	Cfg    *config.AnalyzerConfig
}

func (s *TraceAnalyzerGrpcServer) Run() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		logrus.Fatalf("Failed to listen: %v", err)
	}
	srv := &TraceAnalyzerServer{
		traceId:      make(map[string]bool),
		traceIdMutex: sync.RWMutex{},
		analyseCh:    make(chan string, 100), // 缓冲区大小为100
		cfg:          s.Cfg,
	}
	tracev1.RegisterTraceServiceServer(s.Server, srv)
	logrus.Infof("[Run] - 启动分析协程")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		srv.analyse(ctx.Done())
	}()

	logrus.Infof("[Run] - 监听地址: %s", s.Addr)
	if err := s.Server.Serve(lis); err != nil {
		logrus.Fatalf("Failed to serve: %v", err)
	}
	return nil
}

func (s *TraceAnalyzerServer) analyse(stopCh <-chan struct{}) {
	for {
		select {
		case traceId := <-s.analyseCh:
			logrus.Infof("[analyse] - traceId = %s", traceId)
			go func() {
				time.Sleep(8 * time.Second) // 等待tempo那边把traceid的数据都整理好，然后我们调用tempo的api接口获取这条traceid 的数据
				request, err := openapi.MakeTraceIdRequest(s.cfg.Tempo.Host, traceId, "", "")
				if err != nil {
					logrus.Errorf("[analyse] - 创建请求失败, err = %v", err)
					return
				}
				// todo: 完成接下来的response分析
				response, err := httpc.SendHttpRequest(request)
				data := &otelmodel.TraceData{}
				if err = httpc.MarshalResp(response, data); err != nil {
					logrus.Errorf("[analyse] - 解析响应失败, traceid = %s, err = %v", traceId, err)
					return
				} else {
					// todo: 对data进行分析
					logrus.Infof("[analyse] - data size = %d", len(data.Batches))
				}

			}()
		case <-stopCh:
			logrus.Infof("[analyse] - 关闭分析服务")
			return
		}
	}
}

func (s *TraceAnalyzerGrpcServer) Close() {
	s.Server.GracefulStop()
	logrus.Infof("[TraceAnalyzerGrpcServer] - 关闭成功")
}

func (s *TraceAnalyzerServer) Export(ctx context.Context, request *tracev1.ExportTraceServiceRequest) (*tracev1.ExportTraceServiceResponse, error) {
	for _, each := range request.ResourceSpans {
		for _, span := range each.ScopeSpans {
			for _, eachSpan := range span.Spans {
				traceId := otelutils.OtelIDToString(eachSpan.TraceId)
				logrus.Infof("[export] - traceId = %s", traceId)
				if s.CheckAndSetTraceId(traceId) {
					logrus.Debugf("[export] - traceId = %s 已存在", traceId)
				} else {
					logrus.Debugf("[export] - traceId = %s 不存在", traceId)
					// 不存在的情况下，需要对该traceid做分析，考虑到时间问题，这里做延迟分析
					s.analyseCh <- traceId
				}
			}
		}
	}
	return &tracev1.ExportTraceServiceResponse{}, nil
}

// CheckAndSetTraceId 检查traceId是否存在，如果不存在则设置为true并返回false
func (s *TraceAnalyzerServer) CheckAndSetTraceId(traceId string) bool {
	s.traceIdMutex.Lock()
	exists := s.traceId[traceId]
	s.traceIdMutex.Unlock()

	if !exists {
		s.traceIdMutex.Lock()
		if !s.traceId[traceId] {
			// 双重校验
			s.traceId[traceId] = true
			s.traceIdMutex.Unlock()
			return false
		}
		s.traceIdMutex.Unlock()
	}

	return true
}

func NewTraceAnalyzerServer(cfg *config.AnalyzerConfig) *TraceAnalyzerGrpcServer {
	server := grpc.NewServer()

	return &TraceAnalyzerGrpcServer{
		Addr:   fmt.Sprintf(":%d", cfg.Server.Port),
		Server: server,
		Cfg:    cfg,
	}
}
