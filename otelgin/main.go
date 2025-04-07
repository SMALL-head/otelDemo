package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"otelDemo/otelgin/core/server"
)

func main() {
	s := server.NewOtelGinServer(gin.ReleaseMode, "test111")
	if err := s.Run(":8080"); err != nil {
		logrus.Fatal(err)
	}
}
