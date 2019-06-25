package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	edgeexporter "github.com/takashabe/edge-exporter"
	"go.opencensus.io/trace"
	"golang.org/x/net/netutil"
)

func main() {
	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("PROJECT_ID"),
	})
	if err != nil {
		panic(err)
	}

	edge := edgeexporter.New(edgeexporter.WithExportInterval(time.Second))
	edge.RegisterExporter(sd)

	trace.RegisterExporter(edge)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	http.HandleFunc("/", handler)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	log.Fatal(http.Serve(netutil.LimitListener(l, 100), nil))
}

var cnt int64

func handler(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	_, span := trace.StartSpan(context.Background(), "handler")

	c := atomic.LoadInt64(&cnt)
	if c%5 == 0 {
		time.Sleep(100 * time.Millisecond)
	}

	latency := time.Now().Sub(startTime)
	span.AddAttributes(trace.Int64Attribute("latency", latency.Nanoseconds()))
	span.AddAttributes(trace.Int64Attribute("count", c))
	span.End()

	log.Printf("count=%d, latency=%v", c, latency.Seconds())

	atomic.AddInt64(&cnt, 1)
}
