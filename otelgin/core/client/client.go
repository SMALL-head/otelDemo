package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"otelDemo/otelgin/common/consts"
)

type Client struct {
	bags baggage.Baggage
}

func New() *Client {
	b, _ := baggage.New()
	return &Client{bags: b}
}

func (c *Client) AddBaggage(k, v string) error {
	b, err := baggage.NewMember(k, v)
	if err != nil {
		return err
	}
	if c.bags, err = c.bags.SetMember(b); err != nil {
		return err
	}

	return nil
}

// ReqWithOtelRecord 发送http请求，并且记录链路信息。此处的ctx中需要有网络孪生体的ctid的baggage信息
func (c *Client) ReqWithOtelRecord(ctx context.Context, url, method string, body interface{}, tracerName string) ([]byte, error) {
	//tracer := otel.Tracer(tracerName)
	//ctx, span := tracer.Start(ctx, tracerName)
	//defer span.End()

	fromContextBaggage := baggage.FromContext(ctx)
	c.bags = fromContextBaggage
	member := fromContextBaggage.Member(consts.CybertwinKey)
	if member.Value() == "" {
		logrus.Errorf("[Req]-获取ctid失败")
		return nil, errors.New("获取ctid失败")
	}

	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			logrus.Errorf("[Req]-序列化body失败")
			return nil, err
		}
	}
	r := bytes.NewReader(bodyBytes)

	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		logrus.Errorf("[Req]-封装request请求失败")
		return nil, err
	}

	// baggage封装，这里应该不用封装，因为baggage会伴随ctx一直走完整个链路
	//ctx = baggage.ContextWithBaggage(ctx, c.bags)

	//collect/record otel metrics
	realClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport,
			otelhttp.WithSpanOptions(trace.WithAttributes(attribute.String(consts.CybertwinKey, member.Value()))),
			otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
				return r.Method + " " + r.URL.Path
			})),
	}

	if resp, err := realClient.Do(req); err != nil {
		logrus.Errorf("[Req]-请求失败")
		return nil, err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			logrus.Errorf("[Req]-请求失败， 错误码为%s", resp.Status)
			return nil, errors.New("请求失败, 错误码为" + resp.Status)
		}
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("[Req]-读取响应失败")
			return nil, err
		}
		return bodyBytes, nil
	}
}
