// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package olog

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/logtest"
)

func TestLogger_LevelSpecificEvents(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test all level-specific event methods
	tests := []struct {
		name      string
		logFunc   func()
		severity  log.Severity
		eventName string
	}{
		{
			name:      "TraceEvent",
			logFunc:   func() { logger.TraceEvent(ctx, "trace.event", "key", "value") },
			severity:  log.SeverityTrace,
			eventName: "trace.event",
		},
		{
			name:      "DebugEvent",
			logFunc:   func() { logger.DebugEvent(ctx, "debug.event", "key", "value") },
			severity:  log.SeverityDebug,
			eventName: "debug.event",
		},
		{
			name:      "InfoEvent",
			logFunc:   func() { logger.InfoEvent(ctx, "info.event", "key", "value") },
			severity:  log.SeverityInfo,
			eventName: "info.event",
		},
		{
			name:      "WarnEvent",
			logFunc:   func() { logger.WarnEvent(ctx, "warn.event", "key", "value") },
			severity:  log.SeverityWarn,
			eventName: "warn.event",
		},
		{
			name:      "ErrorEvent",
			logFunc:   func() { logger.ErrorEvent(ctx, "error.event", "key", "value") },
			severity:  log.SeverityError,
			eventName: "error.event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder.Reset()
			tt.logFunc()

			want := logtest.Recording{
				logtest.Scope{
					Name: "test-logger",
				}: {
					logtest.Record{
						Context:   ctx,
						Severity:  tt.severity,
						EventName: tt.eventName,
						Attributes: []log.KeyValue{
							log.String("key", "value"),
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
		})
	}
}

func TestLogger_LevelSpecificEventsAttr(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test all level-specific event attribute methods
	tests := []struct {
		name      string
		logFunc   func()
		severity  log.Severity
		eventName string
	}{
		{
			name:      "TraceEventAttr",
			logFunc:   func() { logger.TraceEventAttr(ctx, "trace.event", log.String("key", "value")) },
			severity:  log.SeverityTrace,
			eventName: "trace.event",
		},
		{
			name:      "DebugEventAttr",
			logFunc:   func() { logger.DebugEventAttr(ctx, "debug.event", log.String("key", "value")) },
			severity:  log.SeverityDebug,
			eventName: "debug.event",
		},
		{
			name:      "InfoEventAttr",
			logFunc:   func() { logger.InfoEventAttr(ctx, "info.event", log.String("key", "value")) },
			severity:  log.SeverityInfo,
			eventName: "info.event",
		},
		{
			name:      "WarnEventAttr",
			logFunc:   func() { logger.WarnEventAttr(ctx, "warn.event", log.String("key", "value")) },
			severity:  log.SeverityWarn,
			eventName: "warn.event",
		},
		{
			name:      "ErrorEventAttr",
			logFunc:   func() { logger.ErrorEventAttr(ctx, "error.event", log.String("key", "value")) },
			severity:  log.SeverityError,
			eventName: "error.event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder.Reset()
			tt.logFunc()

			want := logtest.Recording{
				logtest.Scope{
					Name: "test-logger",
				}: {
					logtest.Record{
						Context:   ctx,
						Severity:  tt.severity,
						EventName: tt.eventName,
						Attributes: []log.KeyValue{
							log.String("key", "value"),
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
		})
	}
}

func TestLogger_EventEnabled(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test EventEnabled method - should return true for Info level
	if !logger.EventEnabled(ctx) {
		t.Error("expected event logging to be enabled")
	}

	// Test with a recorder that's disabled for info level
	disabledRecorder := logtest.NewRecorder(
		logtest.WithEnabledFunc(func(_ context.Context, param log.EnabledParameters) bool {
			return param.Severity > log.SeverityInfo
		}),
	)
	disabledLogger := New(Options{
		Provider: disabledRecorder,
		Name:     "disabled-logger",
	})

	if disabledLogger.EventEnabled(ctx) {
		t.Error("expected event logging to be disabled")
	}
}

func TestLogger_BackwardCompatibleEvent(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test that the original Event method still works and logs at Info level
	logger.Event(ctx, "user.login", "user_id", "12345")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:   ctx,
				Severity:  log.SeverityInfo,
				EventName: "user.login",
				Attributes: []log.KeyValue{
					log.String("user_id", "12345"),
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

func TestLogger_BackwardCompatibleEventAttr(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test that the original EventAttr method still works and logs at Info level
	logger.EventAttr(ctx, "user.login", log.String("user_id", "12345"))

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:   ctx,
				Severity:  log.SeverityInfo,
				EventName: "user.login",
				Attributes: []log.KeyValue{
					log.String("user_id", "12345"),
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

func TestLogger_EventWithPreConfiguredAttributes(t *testing.T) {
	recorder := logtest.NewRecorder()
	baseLogger := New(Options{Provider: recorder, Name: "test-logger"})

	ctx := t.Context()

	// Test event methods with pre-configured attributes
	logger := baseLogger.With("service", "auth", "version", "1.0.0")
	logger.WarnEvent(ctx, "rate.limit.exceeded", "client_ip", "192.168.1.100")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:   ctx,
				Severity:  log.SeverityWarn,
				EventName: "rate.limit.exceeded",
				Attributes: []log.KeyValue{
					log.String("service", "auth"),
					log.String("version", "1.0.0"),
					log.String("client_ip", "192.168.1.100"),
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

func TestLogger_EventAttrWithPreConfiguredAttributes(t *testing.T) {
	recorder := logtest.NewRecorder()
	baseLogger := New(Options{Provider: recorder, Name: "test-logger"})

	ctx := t.Context()

	// Test event attr methods with pre-configured attributes
	logger := baseLogger.WithAttr(log.String("service", "auth"), log.String("version", "1.0.0"))
	logger.WarnEventAttr(ctx, "rate.limit.exceeded", log.String("client_ip", "192.168.1.100"))

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:   ctx,
				Severity:  log.SeverityWarn,
				EventName: "rate.limit.exceeded",
				Attributes: []log.KeyValue{
					log.String("service", "auth"),
					log.String("version", "1.0.0"),
					log.String("client_ip", "192.168.1.100"),
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
