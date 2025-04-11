package openapi

import (
	"fmt"
	"net/http"
	"otelDemo/analyzer/common/consts"
)

func MakeTraceIdRequest(tempoHost, traceid, start, end string) (*http.Request, error) {
	if start == "" && end == "" {
		return http.NewRequest("GET", tempoHost+fmt.Sprintf(consts.TempoTraceIDAPIV2Template, traceid), nil)
	} else {
		return http.NewRequest("GET", tempoHost+fmt.Sprintf(consts.TempoTraceIDAPIWithDateV2Template, traceid, start, end), nil)
	}
}
