package openapi_test

import (
	"github.com/stretchr/testify/require"
	"otelDemo/analyzer/common/otelmodel"
	"otelDemo/utils/httpc"
	"otelDemo/utils/otelutils/openapi"
	"testing"
)

func TestMakeTraceIdRequest(t *testing.T) {
	request, err := openapi.MakeTraceIdRequest("http://127.0.0.1:3200", "36d6e37df638d2a416ff09e6a5c81c1e", "", "")
	require.NoError(t, err)
	resp, err := httpc.SendHttpRequest(request)
	require.NoError(t, err)
	data := &otelmodel.TraceData{}
	err = httpc.MarshalResp(resp, data)
	require.NoError(t, err)
	//logrus.Infof("[TEST] - data size = %d", len(data.Batches)) // 大概也许是成功了，但是这个单元测试case中的traceid需要根据实际情况传递，因此并不好直接运行
}
