package edgeexporter

import (
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

func TestRegisterAndUnregisterExporter(t *testing.T) {
	edge := EdgeExporter{}
	edge.RegisterExporter(&emptyExporter{})
	edge.RegisterExporter(&emptyExporter{})
	edge.RegisterExporter(&emptyExporter2{})
	// TODO: internal test
	assert.Equal(t, len(edge.exporters.Load()), 2)

	edge.UnregisterExporter(&emptyExporter{})
	// TODO: internal test
	assert.Equal(t, len(edge.exporters.Load()), 1)
}
