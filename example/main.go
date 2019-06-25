package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	edgeexporter "github.com/takashabe/edge-exporter"
	"go.opencensus.io/trace"
)

var cnt int32

func main() {
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("PROJECT_ID"),
	})
	if err != nil {
		panic(err)
	}

	edge := edgeexporter.New(edgeexporter.WithExportInterval(10 * time.Second))
	edge.RegisterExporter(sd)

	trace.RegisterExporter(edge)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		_, span := trace.StartSpan(context.Background(), "handler")

		c := atomic.LoadInt32(&cnt)
		if c%5 == 0 {
			time.Sleep(time.Second)
		}

		latency := startTime.Sub(time.Now())
		span.AddAttributes(trace.Int64Attribute("latency", latency.Nanoseconds()))
		span.AddAttributes(trace.Int64Attribute("count", int64(c)))
		span.End()

		atomic.AddInt32(&cnt, 1)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
