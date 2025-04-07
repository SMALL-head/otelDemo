package main

import (
	"encoding/json"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
)

const name = "go.opentelemetry.io/otel/example/dice"

var (
	tracer = otel.Tracer(name)
	logger = otelslog.NewLogger(name)
)

func init() {
}

func svc2(w http.ResponseWriter, r *http.Request) {
	//span := opentracing.SpanFromContext(r.Context()).SetTag("svc2", "svc2")
	_, span := tracer.Start(r.Context(), "[svc2] Get /svc2")
	//spanAttribute1 := attribute.String("svc2-kv", "lalala")
	//span.SetAttributes(spanAttribute1)
	defer span.End()

	bag := baggage.FromContext(r.Context())

	value := bag.Member("roll.value").Value()
	valueInt, err := strconv.ParseInt(value, 10, 32)
	span.SetAttributes(attribute.Int("roll.value", int(valueInt)))

	// do sth really slow
	time.Sleep(1 * time.Second)
	resp := "return from svc2\n"

	indent, err := json.MarshalIndent(r.Header, "", "  ")
	if err != nil {
		log.Printf("Marshal failed: %v\n", err)
	}
	log.Printf("Request headers: %s\n", indent)

	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}
