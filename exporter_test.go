package edgeexporter

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

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

type fakeCountExporter struct {
	cnt int32
}

func (f *fakeCountExporter) ExportSpan(_ *trace.SpanData) {
	fmt.Println("fakeCounterExporter.ExportSpan")
	atomic.AddInt32(&f.cnt, 1)
}

func TestExportSpan(t *testing.T) {
	tests := []struct {
		interval time.Duration
		wantCnt  int
	}{
		{
			interval: 100 * time.Millisecond,
			wantCnt:  1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			fakeExporter := &fakeCountExporter{}
			exporter := New(WithExportInterval(tt.interval))
			exporter.RegisterExporter(fakeExporter)

			trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
			trace.RegisterExporter(exporter)

			for i := 0; i < 10; i++ {
				_, span := trace.StartSpan(context.Background(), "span")
				span.End()
			}

			assert.Equal(t, int32(tt.wantCnt), fakeExporter.cnt)
		})
	}
}
