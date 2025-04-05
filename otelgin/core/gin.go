package core

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
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
	tracer := otel.Tracer(tracerName)

	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
		c.Request = c.Request.WithContext(ctx) // 将ctx传递下去，如果后续需要使用otel作tracing，就有父span了？
		// 这里可以认为一定能拿到上一个服务中的baggagees中的网络孪生体的id
		bag := baggage.FromContext(c.Request.Context())
		ctID := bag.Member(consts.CybertwinKey).Value()

		// 设置attribute，包含ct信息
		span.SetAttributes(
			attribute.String(consts.CybertwinKey, ctID),
		)
		
		c.Next()
		defer span.End()

	}
}
