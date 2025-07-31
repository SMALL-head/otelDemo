package zerotrust

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	EnvZeroTrustKey   = "ZEROTRUST_ENDPOINT"
	NotifyMsgTemplate = "cybertwin %s has trigger dandger microservice calling chain"
)

var (
	endpoint string
)

func init() {
	e, ok := os.LookupEnv(EnvZeroTrustKey)
	if !ok {
		endpoint = "http://zero-trust.immune-security.svc.cluster.local" // 默认值
	} else {
		endpoint = e
	}

}

func GetZeroTrustEndpoint() string {
	return endpoint
}

func GetScoreByCtid(ctid string, svc string) (int, error) {
	if ctid == "" || svc == "" {
		return 0, errors.New("ctid and svc are required")
	}

	endpoint := GetZeroTrustEndpoint()
	url := fmt.Sprintf("%s/query/cybertwin-trust", endpoint)

	// 将 ctid 字符串转换为整数
	var ctidNum int
	_, err := fmt.Sscanf(ctid, "%d", &ctidNum)
	if err != nil {
		return 0, fmt.Errorf("invalid ctid format, must be a number: %v", err)
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"ctid": ctidNum, // 使用转换后的整数
		"svc":  svc,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 发送POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %v", err)
	}

	// 解析响应
	var response struct {
		Success bool   `json:"success"`
		Score   int    `json:"score"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response.Score, nil
}

func NotifyMicroDanger(ctid string, msg string) error {
	if ctid == "" {
		return errors.New("ctid is required")
	}
	endpoint := GetZeroTrustEndpoint()
	url := fmt.Sprintf("%s/notify/microservice", endpoint)
	ctidInt, err := strconv.Atoi(ctid)
	if err != nil {
		return fmt.Errorf("invalid ctid format, must be a number: %v", err)
	}
	// 构建请求体
	req := struct {
		Ctid int    `json:"ctid"`
		Msg  string `json:"msg"`
	}{
		Ctid: ctidInt,
		Msg:  msg,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}
	// 发送POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	// 解析响应
	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	if !response.Success {
		return fmt.Errorf("request failed: %s", response.Message)
	}
	return nil
}
