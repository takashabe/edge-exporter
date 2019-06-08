package edgeexporeter

import (
	"github.com/k0kubun/pp"
	"go.opencensus.io/trace"
)

type EdgeExporter struct{}

func (e *EdgeExporter) ExportSpan(sd *trace.SpanData) {
	pp.Println(sd)
}
