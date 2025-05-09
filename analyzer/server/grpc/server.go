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
	"otelDemo/analyzer/httpclient"
	"otelDemo/analyzer/svc"
	"otelDemo/db/dao"
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

	svc *svc.Svc
}

type TraceAnalyzerGrpcServer struct {
	Addr   string
	Server *grpc.Server
	Cfg    *config.AnalyzerConfig

	traceAnalyzerSvc *TraceAnalyzerServer
}

func (s *TraceAnalyzerGrpcServer) Init() {
	conn, err := dao.NewDBConn(context.Background(), s.Cfg.DataSource.Host)
	if err != nil {
		logrus.Fatalf("[Init] - 数据库连接失败, err = %v", err)
	}
	db := dao.NewDBTXQuery(conn)
	srv := &TraceAnalyzerServer{
		traceId:      make(map[string]bool),
		traceIdMutex: sync.RWMutex{},
		analyseCh:    make(chan string, 100), // 缓冲区大小为100
		cfg:          s.Cfg,
		svc:          svc.New(db),
	}

	s.traceAnalyzerSvc = srv
	tracev1.RegisterTraceServiceServer(s.Server, srv)
}

func (s *TraceAnalyzerGrpcServer) Run() error {
	if s.Server == nil || s.traceAnalyzerSvc == nil {
		logrus.Fatalf("[Run] - grpc server 未初始化")
	}

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		logrus.Fatalf("Failed to listen: %v", err)
	}

	logrus.Infof("[Run] - 启动分析协程")
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		s.traceAnalyzerSvc.analyse(ctx.Done())
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

			go func() {
				time.Sleep(6 * time.Second) // 等待tempo那边把traceid的数据都整理好，然后我们调用tempo的api接口获取这条traceid 的数据
				logrus.Infof("[analyse] - start analyse traceId = %s", traceId)
				request, err := openapi.MakeTraceIdRequest(config.ApplicationConfig.Tempo.Host, traceId, "", "")
				if err != nil {
					logrus.Errorf("[analyse] - 创建请求失败, err = %v", err)
					return
				}
				// todo: 完成接下来的response分析
				response, err := httpc.SendHttpRequest(request)
				tData := &otelmodel.TraceData{}
				if err = httpc.MarshalResp(response, tData); err != nil {
					logrus.Errorf("[analyse] - 解析响应失败, traceid = %s, err = %v", traceId, err)
					return
				} else {
					// todo: 对data进行分析
					patternInfos, err := s.svc.GetAllPattern(context.TODO())
					patterns := make([]*otelmodel.PatternTree, 0)
					cybertwinInfo := otelmodel.GetCybertwinInfoFromTraceData(tData)
					if cybertwinInfo == "" {
						logrus.Warnf("[analyse] - traceId = %s, 没有cybertwin信息", traceId)
					}
					for _, each := range patternInfos {
						tree, err := otelmodel.Pattern2Tree(int(each.ID), each.Name, each.GraphData)
						if err != nil {
							logrus.Errorf("[analyse] - 将pattern转化为tree失败, err = %v, pattern_name = %s", each.Name, err)
							continue
						}
						patterns = append(patterns, tree)
					}
					if err != nil {
						logrus.Errorf("[analyse] - 获取所有pattern失败, err = %v", err)
						return
					}
					tree, err := otelmodel.TransferTraceData2Tree(tData)
					if err != nil {
						logrus.Errorf("[analyse] - 将tracedata转化为tree失败, err = %v", err)
						return
					}
					match(patterns, tree, cybertwinInfo)
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

// Export 实现grpc otel-collector的TracerServer接口，当有trace上报的时候会调用这个接口。
// 注意，这个接口并不是按照traceid为单位的上报数据，也就是说，同一个traceid下的span可能会分多次上报
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

func match(patterns []*otelmodel.PatternTree, tree *otelmodel.TraceDataTree, cybertwinLabel string) {
	for _, pattern := range patterns {
		if otelmodel.MatchPattern(pattern, tree) {
			// todo: 匹配成功，进行后续处理
			logrus.Infof("[analyse] - [match] - 匹配成功, traceId = %s, pattern = %s", tree.TraceId, pattern.Root.Value)
			//res := httpclient.FlareAdmin.HelloClient()
			//logrus.Info("[analyse] - [match] - 匹配成功, res = ", res)
			res := httpclient.FlareAdmin.AddMatchResultRecord(pattern.ID, 0, cybertwinLabel)
			logrus.Infof("[analyse] - [match] - 添加结果, res = %v", res)
		}
	}
}
