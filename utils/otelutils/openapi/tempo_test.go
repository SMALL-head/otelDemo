package openapi_test

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"otelDemo/analyzer/common/otelmodel"
	"otelDemo/utils/httpc"
	"otelDemo/utils/otelutils/openapi"
	"testing"
)

func TestMakeTraceIdRequest(t *testing.T) {
	request, err := openapi.MakeTraceIdRequest("http://127.0.0.1:3200", "b3dfc76089f6b3f4e05c04084f7d2895", "", "")
	require.NoError(t, err)
	resp, err := httpc.SendHttpRequest(request)
	require.NoError(t, err)
	data := &otelmodel.TraceData{}
	err = httpc.MarshalResp(resp, data)
	require.NoError(t, err)
	logrus.Infof("[TEST] - data size = %d", len(data.Trace.ResourceSpans)) // 大概也许是成功了，但是这个单元测试case中的traceid需要根据实际情况传递，因此并不好直接运行
}
