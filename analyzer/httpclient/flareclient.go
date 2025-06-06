package httpclient

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (c *FlareAdminClient) HelloClient() string {
	get, err := c.DoGet("/debug/hello")

	if err != nil {
		return "error"
	}
	defer get.Body.Close()
	if get.StatusCode != http.StatusOK {
		return "status code error"
	}
	all, err := io.ReadAll(get.Body)
	if err != nil {
		return "read error"
	}
	return string(all)
}

func (c *FlareAdminClient) AddMatchResultRecord(patternID int, cybertwinID int, cybertwinLabel string) string {
	req := struct {
		PatternID      int    `json:"pattern_id"`
		CybertwinID    int    `json:"cybertwin_id"`
		CybertwinLabel string `json:"cybertwin_label"`
	}{
		PatternID:      patternID,
		CybertwinID:    cybertwinID,
		CybertwinLabel: cybertwinLabel,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		logrus.Errorf("[AddMatchResultRecord] - [json.Marshal] - failed to marshal json: %v", err)
		return "error"
	}
	response, err := c.DoPost("micro/pattern/result/add", reqBytes, "application/json")

	if err != nil {
		return "error"
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Sprintf("status code error: %d", response.StatusCode)
	}
	all, err := io.ReadAll(response.Body)
	if err != nil {
		return "read error"
	}
	return string(all)
}
