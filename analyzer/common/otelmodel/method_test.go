package otelmodel_test

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"otelDemo/analyzer/common/otelmodel"
	"testing"
)

func TestMatchPattern(t *testing.T) {
	traceBytes := []byte("{\"trace\":{\"resourceSpans\":[{\"resource\":{\"attributes\":[{\"key\":\"service.name\",\"value\":{\"stringValue\":\"svc3\"}}]},\"scopeSpans\":[{\"scope\":{\"name\":\"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp\",\"version\":\"0.58.0\"},\"spans\":[{\"traceId\":\"UjJIhALnUOLgdmVBdK3ASQ==\",\"spanId\":\"Fxkl8KJCalA=\",\"parentSpanId\":\"c5APpKy18rQ=\",\"name\":\"GET /svc4\",\"kind\":\"SPAN_KIND_CLIENT\",\"startTimeUnixNano\":\"1744357875298104000\",\"endTimeUnixNano\":\"1744357875311743375\",\"attributes\":[{\"key\":\"cybertwin_id\",\"value\":{\"stringValue\":\"simct\"}},{\"key\":\"net.peer.name\",\"value\":{\"stringValue\":\"127.0.0.1\"}},{\"key\":\"net.peer.port\",\"value\":{\"intValue\":\"8084\"}},{\"key\":\"http.response_content_length\",\"value\":{\"intValue\":\"29\"}},{\"key\":\"http.method\",\"value\":{\"stringValue\":\"GET\"}},{\"key\":\"http.url\",\"value\":{\"stringValue\":\"http://127.0.0.1:8084/svc4\"}},{\"key\":\"http.status_code\",\"value\":{\"intValue\":\"200\"}}],\"status\":{}}]}]},{\"resource\":{\"attributes\":[{\"key\":\"service.name\",\"value\":{\"stringValue\":\"svc3\"}}]},\"scopeSpans\":[{\"scope\":{\"name\":\"svc3\"},\"spans\":[{\"traceId\":\"UjJIhALnUOLgdmVBdK3ASQ==\",\"spanId\":\"c5APpKy18rQ=\",\"name\":\"/tosvc4\",\"kind\":\"SPAN_KIND_SERVER\",\"startTimeUnixNano\":\"1744357875295559000\",\"endTimeUnixNano\":\"1744357875313250833\",\"attributes\":[{\"key\":\"cybertwin_id\",\"value\":{\"stringValue\":\"simct\"}}],\"status\":{}}]}]},{\"resource\":{\"attributes\":[{\"key\":\"service.name\",\"value\":{\"stringValue\":\"svc4\"}}]},\"scopeSpans\":[{\"scope\":{\"name\":\"svc4\"},\"spans\":[{\"traceId\":\"UjJIhALnUOLgdmVBdK3ASQ==\",\"spanId\":\"0QiaKBSuHLU=\",\"parentSpanId\":\"Fxkl8KJCalA=\",\"name\":\"/svc4\",\"kind\":\"SPAN_KIND_SERVER\",\"startTimeUnixNano\":\"1744357875308397000\",\"endTimeUnixNano\":\"1744357875310311208\",\"attributes\":[{\"key\":\"cybertwin_id\",\"value\":{\"stringValue\":\"simct\"}}],\"status\":{}}]}]}]}}")
	traceData := &otelmodel.TraceData{}
	err := json.Unmarshal(traceBytes, traceData)
	require.NoError(t, err)
	traceDataTree, err := otelmodel.TransferTraceData2Tree(traceData)
	require.NoError(t, err)
	pTree, err := otelmodel.Pattern2Tree("test", []byte("{\"edges\": [{\"label\": \"HTTP GET /foo\", \"source\": \"a\", \"target\": \"b\"}], \"nodes\": [{\"id\": \"a\", \"label\": \"svc3~/tosvc4\"}, {\"id\": \"b\", \"label\": \"svc4~/svc4\"}]}"))
	require.NoError(t, err)
	res := otelmodel.MatchPattern(pTree, traceDataTree)
	require.Equal(t, true, res)
}

func TestMatchPattern2(t *testing.T) {
	sharedTData := &otelmodel.TraceDataTree{
		TraceId: "123",
		Root: &otelmodel.TraceDataNode{
			SpanId:  "1",
			Value:   "v1",
			Service: "svc1",
			Children: []*otelmodel.TraceDataNode{
				{
					SpanId:  "2",
					Value:   "v2",
					Service: "svc2",
					Children: []*otelmodel.TraceDataNode{
						{
							SpanId:  "2-1",
							Value:   "v2-1",
							Service: "svc4",
							Children: []*otelmodel.TraceDataNode{
								{
									SpanId:   "2-1-1",
									Value:    "v2-1-1",
									Service:  "svc5",
									Children: []*otelmodel.TraceDataNode{},
								},
							},
						},
						{
							SpanId:  "2-2",
							Value:   "v2-2",
							Service: "svc3",
							Children: []*otelmodel.TraceDataNode{
								{
									SpanId:   "2-2-1",
									Value:    "v2-2-1",
									Service:  "svc6",
									Children: []*otelmodel.TraceDataNode{},
								},
							},
						},
						{
							SpanId:   "2-3",
							Value:    "v2-3",
							Service:  "svc4",
							Children: []*otelmodel.TraceDataNode{},
						},
					},
				},
			},
		},
	}
	cases := []struct {
		Name     string
		TData    *otelmodel.TraceDataTree
		PTree    *otelmodel.PatternTree
		Expected bool
	}{
		{
			Name:  "case1",
			TData: sharedTData,
			PTree: &otelmodel.PatternTree{
				Name: "pattern1",
				Root: &otelmodel.PatternTreeNode{
					Value:   "v2",
					Service: "svc2",
					Children: []*otelmodel.PatternTreeNode{
						{
							Value:    "v2-1",
							Service:  "svc4",
							Children: []*otelmodel.PatternTreeNode{},
						},
						{
							Value:   "v2-3",
							Service: "svc4",
						},
					},
				},
			},
			Expected: true,
		},
		{
			Name:  "case2",
			TData: sharedTData,
			PTree: &otelmodel.PatternTree{
				Name: "pattern1",
				Root: &otelmodel.PatternTreeNode{
					Value:   "v2",
					Service: "svc2",
					Children: []*otelmodel.PatternTreeNode{
						{
							Value:    "v2-1",
							Service:  "svc4",
							Children: []*otelmodel.PatternTreeNode{},
						},
						{
							Value:   "v2-4",
							Service: "svc4",
						},
					},
				},
			},
			Expected: false,
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			res := otelmodel.MatchPattern(c.PTree, c.TData)
			require.Equal(t, c.Expected, res)
		})
	}
}
