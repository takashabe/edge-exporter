package edgeexporter

import (
	"reflect"
	"sync"
	"time"

	"go.opencensus.io/trace"
)

type exportersList struct {
	mu   sync.RWMutex
	list []trace.Exporter
}

func (l *exportersList) Store(e trace.Exporter) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, exp := range l.list {
		if reflect.DeepEqual(e, exp) {
			return false
		}
	}
	l.list = append(l.list, e)
	return true
}

func (l *exportersList) Load() []trace.Exporter {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.list
}

func (l *exportersList) Delete(e trace.Exporter) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, exp := range l.list {
		if reflect.DeepEqual(e, exp) {
			l.list = append(l.list[:i], l.list[i+1:]...)
			return true
		}
	}
	return false
}

type EdgeExporter struct {
	exporters exportersList

	Interval time.Duration
}

func (e *EdgeExporter) RegisterExporter(exp trace.Exporter) {
	e.exporters.Store(exp)
}

func (e *EdgeExporter) UnregisterExporter(exp trace.Exporter) {
	e.exporters.Delete(exp)
}

func (e *EdgeExporter) ExportSpan(sd *trace.SpanData) {
	// TODO: Store spans and choose a tail latency span
	//       emit the spans each exporter every interval times.

	for _, exp := range e.exporters.Load() {
		exp.ExportSpan(sd)
	}
}
