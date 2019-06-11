package edgeexporter

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

type (
	emptyExporter  struct{}
	emptyExporter2 struct{}
)

func (*emptyExporter) ExportSpan(_ *trace.SpanData) {}

func (*emptyExporter2) ExportSpan(_ *trace.SpanData) {}

// TODO: internal test
func TestExportersList(t *testing.T) {
	edge := EdgeExporter{}
	edge.RegisterExporter(&emptyExporter{})
	edge.RegisterExporter(&emptyExporter{})
	edge.RegisterExporter(&emptyExporter2{})
	assert.Equal(t, len(edge.exporters.Load()), 2)

  edge.UnregisterExporter(&emptyExporter{})
	assert.Equal(t, len(edge.exporters.Load()), 1)
}

func TestBundling(t *testing.T) {
	exporter := &EdgeExporter{
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
