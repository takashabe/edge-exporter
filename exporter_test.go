package edgeexporter

import (
	"context"
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
	atomic.AddInt32(&f.cnt, 1)
}

func TestExportSpanWithLimiter(t *testing.T) {
	tests := []struct {
		exportInterval  time.Duration
		spanEndInterval time.Duration
		spanEndCount    int
		wantCnt         int
	}{
		{
			exportInterval:  100 * time.Millisecond,
			spanEndInterval: 10 * time.Millisecond,
			spanEndCount:    5,
			wantCnt:         1,
		},
		{
			exportInterval:  100 * time.Millisecond,
			spanEndInterval: 10 * time.Millisecond,
			spanEndCount:    15,
			wantCnt:         2,
		},
		{
			exportInterval:  200 * time.Millisecond,
			spanEndInterval: 10 * time.Millisecond,
			spanEndCount:    15,
			wantCnt:         1,
		},
		{
			exportInterval:  0,
			spanEndInterval: 0,
			spanEndCount:    5,
			wantCnt:         5,
		},
	}
	for _, tt := range tests {
		fakeExporter := &fakeCountExporter{}
		exporter := New(WithExportInterval(tt.exportInterval))
		exporter.RegisterExporter(fakeExporter)

		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		trace.RegisterExporter(exporter)

		for i := 0; i < tt.spanEndCount; i++ {
			_, span := trace.StartSpan(context.Background(), "span")
			span.End()
			time.Sleep(tt.spanEndInterval)
		}

		assert.Equal(t, int32(tt.wantCnt), fakeExporter.cnt)
	}
}

type fakeTailExporter struct {
	sd *trace.SpanData
}

func (f *fakeTailExporter) ExportSpan(sd *trace.SpanData) {
	f.sd = sd
}

func TestExportSpanWithTailLantecy(t *testing.T) {
	fakeTailExporter := &fakeTailExporter{}
	exporter := New(WithExportInterval(50 * time.Millisecond))
	exporter.RegisterExporter(fakeTailExporter)

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	trace.RegisterExporter(exporter)

	_, span := trace.StartSpan(context.Background(), "span")
	span.End()
	// consume a first limiter token
	time.Sleep(50 * time.Millisecond)

	_, span2 := trace.StartSpan(context.Background(), "span2")
	time.Sleep(20 * time.Millisecond)
	span2.End()

	_, span3 := trace.StartSpan(context.Background(), "span3")
	span3.End()

	assert.Equal(t, "span2", fakeTailExporter.sd.Name)
}
