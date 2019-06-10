package edgeexporter

import (
	"io"

	"github.com/k0kubun/pp"
	"go.opencensus.io/trace"
)

type EdgeExporter struct {
	OutStream io.Writer
}

func (e *EdgeExporter) ExportSpan(sd *trace.SpanData) {
	pp.Fprintln(e.OutStream, sd)
}
