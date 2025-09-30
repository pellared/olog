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

func TestLogger_LevelSpecificEventEnabled(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test all level-specific EventEnabled methods
	eventNames := []string{"test.event", "user.login", "payment.processed", "system.startup"}

	tests := []struct {
		name        string
		enabledFunc func(context.Context, string) bool
		severity    string
	}{
		{
			name:        "TraceEventEnabled",
			enabledFunc: logger.TraceEventEnabled,
			severity:    "trace",
		},
		{
			name:        "DebugEventEnabled",
			enabledFunc: logger.DebugEventEnabled,
			severity:    "debug",
		},
		{
			name:        "InfoEventEnabled",
			enabledFunc: logger.InfoEventEnabled,
			severity:    "info",
		},
		{
			name:        "WarnEventEnabled",
			enabledFunc: logger.WarnEventEnabled,
			severity:    "warn",
		},
		{
			name:        "ErrorEventEnabled",
			enabledFunc: logger.ErrorEventEnabled,
			severity:    "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, eventName := range eventNames {
				if !tt.enabledFunc(ctx, eventName) {
					t.Errorf("expected %s event logging to be enabled for event: %s", tt.severity, eventName)
				}
			}
		})
	}

	// Test with a recorder that's disabled for info level and below
	disabledRecorder := logtest.NewRecorder(
		logtest.WithEnabledFunc(func(_ context.Context, param log.EnabledParameters) bool {
			return param.Severity > log.SeverityInfo
		}),
	)
	disabledLogger := New(Options{
		Provider: disabledRecorder,
		Name:     "disabled-logger",
	})

	// Test that lower levels are disabled
	for _, eventName := range eventNames {
		if disabledLogger.TraceEventEnabled(ctx, eventName) {
			t.Errorf("expected trace event logging to be disabled for event: %s", eventName)
		}
		if disabledLogger.DebugEventEnabled(ctx, eventName) {
			t.Errorf("expected debug event logging to be disabled for event: %s", eventName)
		}
		if disabledLogger.InfoEventEnabled(ctx, eventName) {
			t.Errorf("expected info event logging to be disabled for event: %s", eventName)
		}

		// But higher levels should be enabled
		if !disabledLogger.WarnEventEnabled(ctx, eventName) {
			t.Errorf("expected warn event logging to be enabled for event: %s", eventName)
		}
		if !disabledLogger.ErrorEventEnabled(ctx, eventName) {
			t.Errorf("expected error event logging to be enabled for event: %s", eventName)
		}
	}
}

func TestLogger_EventWithSeverity(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test that the Event method works with explicit severity level
	logger.Event(ctx, log.SeverityInfo, "user.login", "user_id", "12345")

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

func TestLogger_EventAttrWithSeverity(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test that the EventAttr method works with explicit severity level
	logger.EventAttr(ctx, log.SeverityInfo, "user.login", log.String("user_id", "12345"))

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

func TestLogger_GenericEventMethods(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test generic Event method with custom severity levels
	tests := []struct {
		name      string
		logFunc   func()
		severity  log.Severity
		eventName string
	}{
		{
			name:      "Event with Trace severity",
			logFunc:   func() { logger.Event(ctx, log.SeverityTrace, "custom.trace.event", "key", "value") },
			severity:  log.SeverityTrace,
			eventName: "custom.trace.event",
		},
		{
			name:      "Event with Warn2 severity",
			logFunc:   func() { logger.Event(ctx, log.SeverityWarn2, "custom.warn2.event", "key", "value") },
			severity:  log.SeverityWarn2,
			eventName: "custom.warn2.event",
		},
		{
			name:      "EventAttr with Error2 severity",
			logFunc:   func() { logger.EventAttr(ctx, log.SeverityError2, "custom.error2.event", log.String("key", "value")) },
			severity:  log.SeverityError2,
			eventName: "custom.error2.event",
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
