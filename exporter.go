package edgeexporter

import (
	"reflect"
	"sync"
	"time"

	"go.opencensus.io/trace"
	"golang.org/x/time/rate"
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
	tail      *trace.SpanData
	exportMu  sync.Mutex

	interval time.Duration
	limiter  *rate.Limiter
}

const (
	DefaultExportInterval = time.Second
)

type Option func(*EdgeExporter)

func WithExportInterval(t time.Duration) Option {
	return func(e *EdgeExporter) {
		e.interval = t
	}
}

func New(opts ...Option) *EdgeExporter {
	e := &EdgeExporter{
		interval: DefaultExportInterval,
	}

	for _, opt := range opts {
		opt(e)
	}

	limit := rate.Every(e.interval)
	e.limiter = rate.NewLimiter(limit, 1)

	return e
}

func (e *EdgeExporter) RegisterExporter(exp trace.Exporter) {
	e.exporters.Store(exp)
}

func (e *EdgeExporter) UnregisterExporter(exp trace.Exporter) {
	e.exporters.Delete(exp)
}

func (e *EdgeExporter) ExportSpan(sd *trace.SpanData) {
	e.storeTailLatencySpan(sd)
	if !e.limiter.Allow() {
		return
	}

	e.exportMu.Lock()
	for _, exp := range e.exporters.Load() {
		exp.ExportSpan(sd)
	}
	e.exportMu.Unlock()
}

func (e *EdgeExporter) storeTailLatencySpan(sd *trace.SpanData) {
	if e.tail == nil {
		e.tail = sd
		return
	}

	oldLatency := e.tail.EndTime.Sub(e.tail.StartTime)
	latency := sd.EndTime.Sub(sd.StartTime)
	if oldLatency < latency {
		e.tail = sd
	}
}
