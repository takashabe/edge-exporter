# edge-exporter

_Currently, this project is in alpha state_

edge-exporter provides tracing the tail latency span.

## Example

```
package main

import (
  "context"
  "net/http"
  "os"
  "sync/atomic"
  "time"

  "contrib.go.opencensus.io/exporter/stackdriver"
  edgeexporter "github.com/takashabe/edge-exporter"
  "go.opencensus.io/trace"
)

func main() {
  sd, err := stackdriver.NewExporter(stackdriver.Options{
    ProjectID: os.Getenv("PROJECT_ID"),
  })
  if err != nil {
    panic(err)
  }

  edge := edgeexporter.New(edgeexporter.WithExportInterval(10*time.Second))
  edge.RegisterExporter(sd)

  trace.RegisterExporter(edge)
  trace.ApplyConfig(trace.Config{
    DefaultSampler: trace.AlwaysSample(),
  })

  http.HandleFunc("/", handler)

  http.ListenAndServe(":8080", nil)
}

var cnt int64

func handler(w http.ResponseWriter, req *http.Request) {
  _, span := trace.StartSpan(context.Background(), "handler")
  c := atomic.LoadInt64(&cnt)
  if c%5 == 0 {
    time.Sleep(50 * time.Millisecond)
  }
  span.End()
  atomic.AddInt64(&cnt, 1)
}
```
