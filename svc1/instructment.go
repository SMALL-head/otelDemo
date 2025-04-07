package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
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

	// 设置span只会在当前的span中展示这个标签，但是如果我希望这个attribute能够透传到下一个服务调用方，我应该怎么办？
	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)
	//span.SetTag("roll.value", roll)
	//rollCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	// 书接上回，我们可以通过使用Baggage来达到透传span attribute的效果。如果想要通过Baggage进行透传，需要保证otelsdk中设置了对应的propagator

	baggageMember, _ := baggage.NewMember("roll.value", strconv.Itoa(roll))
	bag, _ := baggage.New(baggageMember)
	ctx = baggage.ContextWithBaggage(ctx, bag)

	resp := strconv.Itoa(roll) + "\n"

	// request for svc2
	//_, span2 := tracer.Start(ctx, "svc2")
	reqForSvc2(ctx)

	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

func reqForSvc2(ctx context.Context) {
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
}
