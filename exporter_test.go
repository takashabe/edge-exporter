package edgeexporter_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	edgeexporter "github.com/takashabe/edge-exporter"
	"go.opencensus.io/trace"
)

func TestBundling(t *testing.T) {
	exporter := &edgeexporter.EdgeExporter{
		OutStream: os.Stdout,
	}
	trace.RegisterExporter(exporter)

	for i := 0; i < 10; i++ {
		_, span := trace.StartSpan(
			context.Background(),
			fmt.Sprintf("span_%d", i),
			trace.WithSampler(trace.AlwaysSample()),
		)
		span.End()
	}
}
