/*
Copyright © 2025 zyc <skyzyc@126.com>
*/
package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"otelDemo/analyzer/common/consts"
	"otelDemo/analyzer/config"
	"otelDemo/analyzer/server/grpc"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "analyzer",
	Short: "轮询查找最新的traceid并且做一系列处理",
	Long:  `没啥好说的`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.analyzer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().Int32P("interval", "i", 1000, "轮询间隔,单位为默认为毫秒,目前不支持修改。已废弃，现在改成响应式的了")
	rootCmd.Flags().StringP("config", "c", "application.yaml", "配置文件,注意目录相对位置")

}

//func run(cmd *cobra.Command, args []string) {
//	configFile := cmd.Flag("config").Value.String()
//	cfg, err := config.LoadConfig(configFile)
//	if err != nil {
//		return
//	}
//	intervalStr := cmd.Flag("interval")
//	interval, err := strconv.ParseInt(intervalStr.Value.String(), 10, 32)
//	if err != nil {
//		logrus.Errorf("[run] - 轮询时间间隔解析失败")
//		return
//	}
//
//	// 每个gorouting需要这个chan来停止
//	stopChs := make([]chan os.Signal, 0)
//
//	wg := sync.WaitGroup{} // 主程序等待所有gorouting都退出后才能退出
//
//	sigIntCh := make(chan os.Signal) // 用于接收程序退出的信号
//	signal.Notify(sigIntCh, os.Interrupt)
//
//	s1 := make(chan os.Signal)
//	stopChs = append(stopChs, s1)
//	go func(stopCh chan os.Signal) {
//		ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
//		for {
//			select {
//			case <-ticker.C:
//				// do something
//
//				logrus.Infof("[run] - ticker 触发")
//			case <-stopCh:
//				wg.Done()
//				return
//			}
//		}
//
//	}(s1)
//
//	wg.Add(len(stopChs))
//
//	<-sigIntCh
//	logrus.Infof("[run] - 等待所有gorouting关闭")
//	for _, eachStopCh := range stopChs {
//		eachStopCh <- os.Interrupt
//	}
//
//	wg.Wait()
//}

func run(cmd *cobra.Command, args []string) {
	// configFile := cmd.Flag("config").Value.String()
	// cfg, err := config.LoadConfig(configFile)
	//if err != nil {
	//	logrus.Fatalf("[run] - 加载配置文件失败, err = %v", err)
	//	return
	//}
	sigIntCh := make(chan os.Signal) // 用于接收程序退出的信号
	signal.Notify(sigIntCh, os.Interrupt)
	cfg := config.AnalyzerConfigDeepCopy(config.ApplicationConfig)
	grpcServer := grpc.NewTraceAnalyzerServer(cfg)
	grpcServer.Init()
	go func() {
		if err := grpcServer.Run(); err != nil {
			logrus.Fatalf("[run] - 运行grpc服务器失败")
		}
	}()
	<-sigIntCh
	grpcServer.Close()
}

func AnalyseTraceId(tempoHost string) {
	// 1. 查询时间窗口内的traceid

}

func makeTraceIdRequest(tempoHost, start, end string) (*http.Request, error) {
	return http.NewRequest("GET", tempoHost+consts.TempoSearchAPI+fmt.Sprintf("?start=%s&end=%s", start, end), nil)
}
