package server

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"otelDemo/otelgin/common/consts"
)

// OtelGinEngine 添加了链路追踪能力的gin server
type OtelGinEngine struct {
	*gin.Engine
}

func NewOtelGinServer(ginMode string, tracerName string) *gin.Engine {
	gin.SetMode(ginMode)
	engine := gin.New()
	gin.DefaultWriter = io.Discard
	engine.Use(gin.Recovery())
	// 为所有的方法添加otel路由中间件
	engine.Use(otelMiddleware(tracerName))
	return engine
}

func otelMiddleware(tracerName string) gin.HandlerFunc {
	propagator := otel.GetTextMapPropagator()
	tracer := otel.Tracer(tracerName)
	return func(c *gin.Context) {
		// 从请求头中提取traceid和baggage信息，填充入ctx里，这样对于服务内的恒宇handler都能够有这个信息
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		// 生成span，注意这个的spankind设置为server，有利于tempo中的调用关系分析
		ctx, span := tracer.Start(ctx, c.Request.URL.Path, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()
		c.Request = c.Request.WithContext(ctx)
		bag := baggage.FromContext(ctx)
		ctID := bag.Member(consts.CybertwinKey).Value()
		if ctID == "" {
			// 如果没有ctid，有两种情况：
			// 1. 这个请求不是网络孪生体的请求，这种情况直接返回
			// 2. 是网络孪生体发出第一个请求，这个时候ctid位于header中

			// 尝试从header中获取ctid
			ctID = c.Request.Header.Get(consts.HeaderCybertwin)
			if ctID == "" {
				c.JSON(400, gin.H{
					"error": "ctid not found",
				})
				c.Abort()
				return
			} else {
				// 如果header中有ctid，添加到baggage中
				b, _ := baggage.NewMember(consts.CybertwinKey, ctID)
				bag, _ = bag.SetMember(b)
				c.Request = c.Request.WithContext(baggage.ContextWithBaggage(c.Request.Context(), bag)) // 重新组装ctx
			}
		}
		span.SetAttributes(attribute.String(consts.CybertwinKey, ctID))
		c.Next()
	}
}
