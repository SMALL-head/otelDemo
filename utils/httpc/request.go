package httpc

import (
	"encoding/json"
	"io"
	"net/http"
)

// SendHttpRequest 发送traceId请求
func SendHttpRequest(r *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MarshalResp 将http.Response的body反序列化为data
func MarshalResp(resp *http.Response, data interface{}) error {
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(respBytes, data); err != nil {
		return err
	}
	return nil
}
