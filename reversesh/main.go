package reversesh

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"otelDemo/otel"
	"otelDemo/otelgin/core/server"
)

func main() {
	conf, err := otel.LoadApplicationConf("./reversesh/application.yaml")
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
}
