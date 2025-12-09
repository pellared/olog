// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package olog

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/logtest"
)

func TestLogger_Basic(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	logger.Trace(ctx, "trace.event", "key", "value")
	logger.Debug(ctx, "debug.event", "key", "value")
	logger.Info(ctx, "info.event", "key", "value")
	logger.Warn(ctx, "warn.event", "key", "value")
	logger.Error(ctx, "error.event", "key", "value")
	logger.Log(ctx, log.SeverityInfo, "test.event", "key", "value")

	got := recorder.Result()

	// Verify correct number of records
	var count int
	for _, records := range got {
		count += len(records)
	}

	if count != 6 {
		t.Errorf("expected 6 records, got %d", count)
	}
}

func TestLogger_AttrBasic(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	logger.TraceAttr(ctx, "trace.event", log.String("key", "value"))
	logger.DebugAttr(ctx, "debug.event", log.String("key", "value"))
	logger.InfoAttr(ctx, "info.event", log.String("key", "value"))
	logger.WarnAttr(ctx, "warn.event", log.String("key", "value"))
	logger.ErrorAttr(ctx, "error.event", log.String("key", "value"))
	logger.LogAttr(ctx, log.SeverityInfo, "test.event", log.String("key", "value"))

	got := recorder.Result()

	// Verify correct number of records
	var count int
	for _, records := range got {
		count += len(records)
	}

	if count != 6 {
		t.Errorf("expected 6 records, got %d", count)
	}
}

func TestLogger_With(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	withLogger := logger.With("service", "api", "version", "1.0.0")
	withLogger.Info(ctx, "test.event", "additional", "attr")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test",
		}: {
			logtest.Record{
				Context:   ctx,
				Severity:  log.SeverityInfo,
				EventName: "test.event",
				Attributes: []log.KeyValue{
					log.String("service", "api"),
					log.String("version", "1.0.0"),
					log.String("additional", "attr"),
				},
			},
		},
	}

	got := recorder.Result()
	logtest.AssertEqual(t, want, got, logtest.Transform(func(r logtest.Record) logtest.Record {
		r.Timestamp = time.Time{}
		r.ObservedTimestamp = time.Time{}
		return r
	}))
}

func TestLogger_WithAttr(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	withLogger := logger.WithAttr(log.String("service", "api"), log.String("version", "1.0.0"))
	withLogger.InfoAttr(ctx, "test.event", log.String("request_id", "req-123"))

	want := logtest.Recording{
		logtest.Scope{
			Name: "test",
		}: {
			logtest.Record{
				Context:   ctx,
				Severity:  log.SeverityInfo,
				EventName: "test.event",
				Attributes: []log.KeyValue{
					log.String("service", "api"),
					log.String("version", "1.0.0"),
					log.String("request_id", "req-123"),
				},
			},
		},
	}

	got := recorder.Result()
	logtest.AssertEqual(t, want, got, logtest.Transform(func(r logtest.Record) logtest.Record {
		r.Timestamp = time.Time{}
		r.ObservedTimestamp = time.Time{}
		return r
	}))
}

func TestLogger_ChainedWith(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	// Chain multiple With calls
	logger1 := logger.With("service", "api")
	logger2 := logger1.With("version", "1.0.0")
	logger2.Info(ctx, "test.event")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test",
		}: {
			logtest.Record{
				Context:   ctx,
				Severity:  log.SeverityInfo,
				EventName: "test.event",
				Attributes: []log.KeyValue{
					log.String("service", "api"),
					log.String("version", "1.0.0"),
				},
			},
		},
	}

	got := recorder.Result()
	logtest.AssertEqual(t, want, got, logtest.Transform(func(r logtest.Record) logtest.Record {
		r.Timestamp = time.Time{}
		r.ObservedTimestamp = time.Time{}
		return r
	}))
}

func TestLogger_ChainedWithAttr(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	// Chain multiple WithAttr calls
	logger1 := logger.WithAttr(log.String("service", "api"))
	logger2 := logger1.WithAttr(log.String("version", "1.0.0"))
	logger2.InfoAttr(ctx, "test.event")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test",
		}: {
			logtest.Record{
				Context:   ctx,
				Severity:  log.SeverityInfo,
				EventName: "test.event",
				Attributes: []log.KeyValue{
					log.String("service", "api"),
					log.String("version", "1.0.0"),
				},
			},
		},
	}

	got := recorder.Result()
	logtest.AssertEqual(t, want, got, logtest.Transform(func(r logtest.Record) logtest.Record {
		r.Timestamp = time.Time{}
		r.ObservedTimestamp = time.Time{}
		return r
	}))
}

func TestLogger_Enabled(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	if !logger.TraceEnabled(ctx, "trace.event") {
		t.Error("expected trace event to be enabled")
	}
	if !logger.DebugEnabled(ctx, "debug.event") {
		t.Error("expected debug event to be enabled")
	}
	if !logger.InfoEnabled(ctx, "info.event") {
		t.Error("expected info event to be enabled")
	}
	if !logger.WarnEnabled(ctx, "warn.event") {
		t.Error("expected warn event to be enabled")
	}
	if !logger.ErrorEnabled(ctx, "error.event") {
		t.Error("expected error event to be enabled")
	}
}

func TestNew_WithOptions(t *testing.T) {
	recorder := logtest.NewRecorder()

	attrs := attribute.NewSet(
		attribute.String("service", "test-service"),
		attribute.String("version", "1.0.0"),
	)

	logger := New(Options{
		Provider:   recorder,
		Name:       "custom-logger",
		Version:    "2.0.0",
		Attributes: attrs,
	})

	ctx := t.Context()
	logger.Info(ctx, "test.event", "key", "value")

	got := recorder.Result()

	if len(got) != 1 {
		t.Fatalf("expected 1 scope, got %d", len(got))
	}

	for scope := range got {
		if scope.Name != "custom-logger" {
			t.Errorf("expected scope name 'custom-logger', got %q", scope.Name)
		}
		if scope.Version != "2.0.0" {
			t.Errorf("expected scope version '2.0.0', got %q", scope.Version)
		}
		if !scope.Attributes.Equals(&attrs) {
			t.Errorf("expected scope attributes %v, got %v", attrs, scope.Attributes)
		}
	}
}
