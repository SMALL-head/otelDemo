package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"otelDemo/otel"
	"otelDemo/otelgin/core/client"
	"otelDemo/otelgin/core/server"
)

func main() {
	conf, err := otel.LoadApplicationConf("./svc3/application.yaml")
	if err != nil {
		logrus.Errorf("[main]-加载配置文件失败: %v", err)
		return
	}
	otelShutdown, err := otel.SetupOTelSDK(context.Background(), conf)
	if err != nil {
		logrus.Errorf("[main]-初始化OpenTelemetry失败: %v", err)
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	ginServer := server.NewOtelGinServer(gin.ReleaseMode, conf.ServiceName)
	ginServer.GET("/svc3", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from svc3",
		})
	})

	ginServer.GET("/tosvc4", func(c *gin.Context) {
		// baggage信息通过propagation传递
		clt := client.New()
		record, err2 := clt.ReqWithOtelRecord(c.Request.Context(), "http://127.0.0.1:8084/svc4", "GET", nil, conf.ServiceName)
		if err2 != nil {
			logrus.Errorf("[main]-请求svc4失败: %v", err2)
			c.JSON(500, gin.H{
				"message": "请求svc4失败",
			})
			return
		}
		logrus.Infof("[main]-请求svc4成功: %s", string(record))
		c.JSON(200, gin.H{
			"message": "请求svc4成功",
		})
	})

	if err := ginServer.Run(":8083"); err != nil {
		panic(err)
	}
}
