package httpclient

import (
	"otelDemo/analyzer/config"
	"otelDemo/utils/httpc"
)

var (
	OTELClient *httpc.BaseClient
	FlareAdmin *FlareAdminClient
)

type FlareAdminClient struct {
	*httpc.BaseClient
}

func init() {
	// 向 Tempo 后端发起请求的客户端
	OTELClient, _ = httpc.NewClient(config.ApplicationConfig.Httpclient.Tempo.Dev.Host)
	// 向 FlareAdmin 后端发起请求的客户端
	client, _ := httpc.NewClient(config.ApplicationConfig.Httpclient.FlareAdmin.Dev.Host)
	FlareAdmin = &FlareAdminClient{
		BaseClient: client,
	}

	//FlareAdmin = &adminClient
	//FlareAdmin= *FlareAdminClient(client()
}
