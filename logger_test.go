// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package olog

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/logtest"
)

func TestLogger_BasicOperations(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	// Test basic logging methods
	logger.Trace(ctx, "trace message", "key", "value")
	logger.Debug(ctx, "debug message", "key", "value")
	logger.Info(ctx, "info message", "key", "value")
	logger.Warn(ctx, "warn message", "key", "value")
	logger.Error(ctx, "error message", "key", "value")
	logger.Log(ctx, log.SeverityInfo, "log message", "key", "value")
	logger.Event(ctx, log.SeverityInfo, "test.event", "key", "value")

	// Verify all log records were captured
	want := logtest.Recording{
		logtest.Scope{
			Name: "test",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityTrace,
				Body:     log.StringValue("trace message"),
				Attributes: []log.KeyValue{
					log.String("key", "value"),
				},
			},
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityDebug,
				Body:     log.StringValue("debug message"),
				Attributes: []log.KeyValue{
					log.String("key", "value"),
				},
			},
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("info message"),
				Attributes: []log.KeyValue{
					log.String("key", "value"),
				},
			},
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityWarn,
				Body:     log.StringValue("warn message"),
				Attributes: []log.KeyValue{
					log.String("key", "value"),
				},
			},
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityError,
				Body:     log.StringValue("error message"),
				Attributes: []log.KeyValue{
					log.String("key", "value"),
				},
			},
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("log message"),
				Attributes: []log.KeyValue{
					log.String("key", "value"),
				},
			},
			logtest.Record{
				Context:   ctx,
				Severity:  log.SeverityInfo,
				EventName: "test.event",
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
}

func TestLogger_EnabledBasic(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	// Test enabled methods - should work with all levels
	if !logger.TraceEnabled(ctx) {
		t.Error("expected trace level to be enabled")
	}
	if !logger.DebugEnabled(ctx) {
		t.Error("expected debug level to be enabled")
	}
	if !logger.InfoEnabled(ctx) {
		t.Error("expected info level to be enabled")
	}
	if !logger.WarnEnabled(ctx) {
		t.Error("expected warn level to be enabled")
	}
	if !logger.ErrorEnabled(ctx) {
		t.Error("expected error level to be enabled")
	}
}

func TestLogger_WithBasic(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test",
	})

	ctx := t.Context()

	// Test With method
	withLogger := logger.With("service", "test", "version", "1.0.0")
	withLogger.Info(ctx, "test message", "additional", "attr")

	// Test chaining With calls
	chainedLogger := withLogger.With("request_id", "123")
	chainedLogger.Info(ctx, "chained message")

	// Verify all log records were captured
	want := logtest.Recording{
		logtest.Scope{
			Name: "test",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("test message"),
				Attributes: []log.KeyValue{
					log.String("service", "test"),
					log.String("version", "1.0.0"),
					log.String("additional", "attr"),
				},
			},
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("chained message"),
				Attributes: []log.KeyValue{
					log.String("service", "test"),
					log.String("version", "1.0.0"),
					log.String("request_id", "123"),
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

func TestLogger_AllLevels(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test all log levels
	tests := []struct {
		name     string
		logFunc  func()
		severity log.Severity
		message  string
	}{
		{
			name:     "Debug",
			logFunc:  func() { logger.Debug(ctx, "debug message", "level", "debug") },
			severity: log.SeverityDebug,
			message:  "debug message",
		},
		{
			name:     "Info",
			logFunc:  func() { logger.Info(ctx, "info message", "level", "info") },
			severity: log.SeverityInfo,
			message:  "info message",
		},
		{
			name:     "Warn",
			logFunc:  func() { logger.Warn(ctx, "warn message", "level", "warn") },
			severity: log.SeverityWarn,
			message:  "warn message",
		},
		{
			name:     "Error",
			logFunc:  func() { logger.Error(ctx, "error message", "level", "error") },
			severity: log.SeverityError,
			message:  "error message",
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
						Context:  ctx,
						Severity: tt.severity,
						Body:     log.StringValue(tt.message),
						Attributes: []log.KeyValue{
							log.String("level", strings.ToLower(tt.name)),
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

func TestLogger_LogMethod(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test Log method with custom level
	logger.Log(ctx, log.SeverityTrace2, "custom log message", "custom", "attribute")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityTrace2,
				Body:     log.StringValue("custom log message"),
				Attributes: []log.KeyValue{
					log.String("custom", "attribute"),
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

func TestLogger_Event(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test Event method
	logger.Event(ctx, log.SeverityInfo, "user.login", "user_id", "12345", "session", "abc123")

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
					log.String("session", "abc123"),
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

func TestLogger_WithAttributes(t *testing.T) {
	recorder := logtest.NewRecorder()
	baseLogger := New(Options{Provider: recorder, Name: "test-logger"})

	ctx := t.Context()

	// Test With method
	logger := baseLogger.With("service", "user-service", "version", "1.0.0")
	logger.Info(ctx, "test message", "request_id", "req-123")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("test message"),
				Attributes: []log.KeyValue{
					log.String("service", "user-service"),
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
	baseLogger := New(Options{Provider: recorder, Name: "test-logger"})

	ctx := t.Context()

	// Test chained With calls
	logger := baseLogger.With("service", "api").With("version", "2.0").With("env", "test")
	logger.Info(ctx, "chained attributes")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("chained attributes"),
				Attributes: []log.KeyValue{
					log.String("service", "api"),
					log.String("version", "2.0"),
					log.String("env", "test"),
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

func TestLogger_ComplexAttributes(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()
	testErr := errors.New("test error")

	// Test various attribute types
	logger.Info(ctx, "complex attributes",
		"string", "test",
		"int", 42,
		"int64", int64(64),
		"float64", 3.14,
		"bool", true,
		"error", testErr,
	)

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("complex attributes"),
				Attributes: []log.KeyValue{
					log.String("string", "test"),
					log.Int64("int", 42),
					log.Int64("int64", 64),
					log.Float64("float64", 3.14),
					log.Bool("bool", true),
					log.String("error", "test error"),
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

func TestLogger_EnabledMethod(t *testing.T) {
	// Test with a recorder that's disabled for debug level
	recorder := logtest.NewRecorder(
		logtest.WithEnabledFunc(func(_ context.Context, param log.EnabledParameters) bool {
			return param.Severity >= log.SeverityInfo
		}),
	)
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test Enabled methods
	if logger.DebugEnabled(ctx) {
		t.Error("expected debug level to be disabled")
	}

	if !logger.InfoEnabled(ctx) {
		t.Error("expected info level to be enabled")
	}

	if !logger.WarnEnabled(ctx) {
		t.Error("expected warn level to be enabled")
	}

	if !logger.ErrorEnabled(ctx) {
		t.Error("expected error level to be enabled")
	}
}

func TestLogger_OddNumberOfArgs(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test with odd number of args - should handle gracefully
	logger.Info(ctx, "test message", "key1", "value1", "key2")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("test message"),
				Attributes: []log.KeyValue{
					log.String("key1", "value1"),
					log.String("key2", ""), // odd arg should get empty value
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

func TestLogger_EmptyMessage(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test with empty message
	logger.Info(ctx, "", "key", "value")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue(""),
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
}

func TestLogger_ContextPropagation(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	// Create a context with a value
	type contextKey string
	key := contextKey("test-key")
	ctx := context.WithValue(t.Context(), key, "test-value")

	logger.Info(ctx, "context test")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("context test"),
			},
		},
	}

	got := recorder.Result()
	logtest.AssertEqual(t, want, got, logtest.Transform(func(r logtest.Record) logtest.Record {
		r.Timestamp = time.Time{}
		r.ObservedTimestamp = time.Time{}
		return r
	}))

	// Additionally verify the context value is preserved
	records := got[logtest.Scope{Name: "test-logger"}]
	if len(records) > 0 {
		if val := records[0].Context.Value(key); val != "test-value" {
			t.Errorf("expected context value 'test-value', got %v", val)
		}
	}
}

func TestLogger_AssertEqual(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Log a message
	logger.Info(ctx, "test message", "key1", "value1", "key2", 42)

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("test message"),
				Attributes: []log.KeyValue{
					log.String("key1", "value1"),
					log.Int64("key2", 42),
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

func TestLogger_AttrMethods(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test all attr-based log levels
	tests := []struct {
		name     string
		logFunc  func()
		severity log.Severity
		message  string
		number   int64
	}{
		{
			name: "TraceAttr",
			logFunc: func() {
				logger.TraceAttr(ctx, "trace message", log.String("level", "trace"), log.Int64("number", 0))
			},
			severity: log.SeverityTrace,
			message:  "trace message",
			number:   0,
		},
		{
			name: "DebugAttr",
			logFunc: func() {
				logger.DebugAttr(ctx, "debug message", log.String("level", "debug"), log.Int64("number", 1))
			},
			severity: log.SeverityDebug,
			message:  "debug message",
			number:   1,
		},
		{
			name: "InfoAttr",
			logFunc: func() {
				logger.InfoAttr(ctx, "info message", log.String("level", "info"), log.Int64("number", 2))
			},
			severity: log.SeverityInfo,
			message:  "info message",
			number:   2,
		},
		{
			name: "WarnAttr",
			logFunc: func() {
				logger.WarnAttr(ctx, "warn message", log.String("level", "warn"), log.Int64("number", 3))
			},
			severity: log.SeverityWarn,
			message:  "warn message",
			number:   3,
		},
		{
			name: "ErrorAttr",
			logFunc: func() {
				logger.ErrorAttr(ctx, "error message", log.String("level", "error"), log.Int64("number", 4))
			},
			severity: log.SeverityError,
			message:  "error message",
			number:   4,
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
						Context:  ctx,
						Severity: tt.severity,
						Body:     log.StringValue(tt.message),
						Attributes: []log.KeyValue{
							log.String("level", strings.ToLower(tt.name[:len(tt.name)-4])), // Remove "Attr" suffix
							log.Int64("number", tt.number),
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

func TestLogger_LogAttr(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test LogAttr method with custom level
	logger.LogAttr(ctx, log.SeverityWarn3, "custom log message", log.String("custom", "attribute"))

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityWarn3,
				Body:     log.StringValue("custom log message"),
				Attributes: []log.KeyValue{
					log.String("custom", "attribute"),
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

func TestLogger_EventAttr(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "test-logger",
	})

	ctx := t.Context()

	// Test EventAttr method
	logger.EventAttr(ctx, log.SeverityInfo, "user.login", log.String("user_id", "12345"), log.String("session", "abc123"))

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
					log.String("session", "abc123"),
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
	baseLogger := New(Options{Provider: recorder, Name: "test-logger"})

	ctx := t.Context()

	// Test WithAttr method
	logger := baseLogger.WithAttr(log.String("service", "user-service"), log.String("version", "1.0.0"))
	logger.InfoAttr(ctx, "test message", log.String("request_id", "req-123"))

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("test message"),
				Attributes: []log.KeyValue{
					log.String("service", "user-service"),
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

func TestLogger_ChainedWithAttr(t *testing.T) {
	recorder := logtest.NewRecorder()
	baseLogger := New(Options{Provider: recorder, Name: "test-logger"})

	ctx := t.Context()

	// Test chained WithAttr calls
	logger := baseLogger.WithAttr(log.String("service", "api")).
		WithAttr(log.String("version", "2.0")).
		WithAttr(log.String("env", "test"))
	logger.InfoAttr(ctx, "chained attributes")

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("chained attributes"),
				Attributes: []log.KeyValue{
					log.String("service", "api"),
					log.String("version", "2.0"),
					log.String("env", "test"),
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

func TestLogger_MixedWithAttrAndWith(t *testing.T) {
	recorder := logtest.NewRecorder()
	baseLogger := New(Options{Provider: recorder, Name: "test-logger"})

	ctx := t.Context()

	// Test mixing WithAttr and With calls
	logger := baseLogger.WithAttr(log.String("service", "api")).
		With("version", "2.0").
		WithAttr(log.String("env", "test"))
	logger.InfoAttr(ctx, "mixed attributes", log.String("request_id", "req-456"))

	want := logtest.Recording{
		logtest.Scope{
			Name: "test-logger",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("mixed attributes"),
				Attributes: []log.KeyValue{
					log.String("service", "api"),
					log.String("version", "2.0"),
					log.String("env", "test"),
					log.String("request_id", "req-456"),
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

func TestNew_WithOptions(t *testing.T) {
	// Test with all options using recorder to verify logger name and functionality
	recorder := logtest.NewRecorder()
	attrs := attribute.NewSet(
		attribute.String("service.name", "test-service"),
		attribute.String("service.version", "1.0.0"),
		attribute.Int("port", 8080),
	)

	logger := New(Options{
		Provider:   recorder,
		Name:       "test-logger",
		Version:    "2.1.0",
		Attributes: attrs,
	})

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	ctx := t.Context()
	logger.Info(ctx, "test message", "key", "value")

	// Verify the logger was created with correct name and records logs
	want := logtest.Recording{
		logtest.Scope{
			Name:    "test-logger",
			Version: "2.1.0",
			Attributes: attribute.NewSet(
				attribute.String("service.name", "test-service"),
				attribute.String("service.version", "1.0.0"),
				attribute.Int("port", 8080),
			),
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("test message"),
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
}

func TestNew_WithGlobalProvider(t *testing.T) {
	// Test with nil provider (should use global)
	logger := New(Options{
		Name: "test-global",
	})

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	// Should not panic and should work with global provider
	ctx := t.Context()
	logger.Info(ctx, "test with global provider")

	// Note: We can't easily test the global provider behavior with logtest.Recording
	// since the global provider is typically noop or configured externally.
	// This test primarily ensures the logger creation doesn't panic.
}

func TestNew_MinimalOptions(t *testing.T) {
	// Test with minimal options using recorder to verify logger name
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "minimal-test",
	})

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	ctx := t.Context()
	logger.Info(ctx, "minimal test")

	// Verify the logger was created with correct name
	want := logtest.Recording{
		logtest.Scope{
			Name: "minimal-test",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("minimal test"),
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

func TestNew_WithAttributesOnly(t *testing.T) {
	recorder := logtest.NewRecorder()
	attrs := attribute.NewSet(
		attribute.String("component", "database"),
		attribute.Bool("enabled", true),
	)

	logger := New(Options{
		Provider:   recorder,
		Name:       "attr-test",
		Attributes: attrs,
	})

	if logger == nil {
		t.Fatal("expected non-nil logger")
	}

	ctx := t.Context()
	logger.Info(ctx, "test message", "key", "value")

	// Verify the logger was created with correct name and instrumentation attributes
	want := logtest.Recording{
		logtest.Scope{
			Name: "attr-test",
			Attributes: attribute.NewSet(
				attribute.String("component", "database"),
				attribute.Bool("enabled", true),
			),
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("test message"),
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
}

func TestLogger_EmbeddedLogger(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := New(Options{
		Provider: recorder,
		Name:     "embedded-test",
	})

	ctx := context.Background()

	// Test that the embedded log.Logger methods are accessible
	// Test Enabled method from embedded logger
	enabled := logger.Enabled(ctx, log.EnabledParameters{
		Severity: log.SeverityInfo,
	})
	if !enabled {
		t.Error("Expected logger to be enabled for Info severity")
	}

	// Test direct Emit method from embedded logger
	var record log.Record
	record.SetBody(log.StringValue("direct emit test"))
	record.SetTimestamp(time.Now())
	record.SetSeverity(log.SeverityInfo)
	record.AddAttributes(log.String("test", "embedding"))

	logger.Emit(ctx, record)

	// Verify the record was captured
	recordings := recorder.Result()
	scope := logtest.Scope{Name: "embedded-test"}
	records := recordings[scope]

	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}

	gotRecord := records[0]
	if gotRecord.Body.AsString() != "direct emit test" {
		t.Errorf("Expected body 'direct emit test', got %q", gotRecord.Body.AsString())
	}

	if gotRecord.Severity != log.SeverityInfo {
		t.Errorf("Expected severity Info, got %v", gotRecord.Severity)
	}

	// Check that the test attribute is present
	found := false
	for _, attr := range gotRecord.Attributes {
		if attr.Key == "test" && attr.Value.AsString() == "embedding" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'test=embedding' attribute")
	} // Test that we can assign Logger to log.Logger interface
	var _ log.Logger = logger
}
