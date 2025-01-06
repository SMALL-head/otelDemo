package main

import (
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
)

const name = "go.opentelemetry.io/otel/example/dice"

var (
	tracer = otel.Tracer(name)
	logger = otelslog.NewLogger(name)
	meter  = otel.Meter(name)
	// rollCnt metric.Int64Counter
)

func init() {
	//var err error
	//rollCnt, err = meter.Int64Counter("dice.rolls",
	//	metric.WithDescription("The number of rolls by roll value"),
	//	metric.WithUnit("{roll}"))
	//if err != nil {
	//	panic(err)
	//}
}

func rolldice(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "[svc1] Get /rolldice")
	defer span.End()
	roll := 1 + rand.Intn(6)

	var msg string
	if player := r.PathValue("player"); player != "" {
		msg = fmt.Sprintf("%s is rolling the dice", player)
	} else {
		msg = "Anonymous player is rolling the dice"
	}
	logger.InfoContext(ctx, msg, "result", roll)

	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)
	//span.SetTag("roll.value", roll)
	// rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	resp := strconv.Itoa(roll) + "\n"

	// request for svc2
	//_, span2 := tracer.Start(ctx, "svc2")
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8082/svc2", nil)
	if err != nil {
		log.Printf("request for svc2 err: %v", err)
		return
	}
	// 等价于otelhttp.DefaultClient
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	respSvc2, err := client.Do(req)
	//span2.End()

	if err != nil {
		log.Printf("request for svc2 err: %v", err)
		return
	}
	defer respSvc2.Body.Close()
	var svc2Resp = make([]byte, 1024)
	respSvc2.Body.Read(svc2Resp)
	log.Printf("svc2 response: %v", string(svc2Resp))

	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}
