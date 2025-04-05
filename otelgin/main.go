package main

import (
	""
	"github.com/gin-gonic/gin"
	"otelDemo/otelgin/core"
)

func main() {
	server := core.NewOtelGinServer(gin.ReleaseMode, "test111")
	if err := server.Run(":8080"); err != nil {
		logrus.Fatal(err)
	}
}
