// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package olog // import "github.com/pellared/olog"

import (
	"context"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
)

// Options contains configuration options for creating a Logger.
type Options struct {
	// Provider is the LoggerProvider to use. If nil, the global LoggerProvider is used.
	Provider log.LoggerProvider

	// Name is the name of the logger, typically the package or component name.
	// If empty, the caller's full package name is automatically detected.
	Name string

	// Version is the version of the logger, typically the package or component version.
	Version string

	// Attributes are pre-configured attributes that will be included in all log records.
	Attributes attribute.Set
}

// Logger provides an ergonomic frontend API for OpenTelemetry structured logging.
// It provides convenience methods for common logging patterns while using the
// OpenTelemetry Logs API as the backend.
//
// The Logger offers two styles of API:
//   - Argument-based methods (Debug, Info, Warn, Error, Log, Event, With) that accept
//     alternating key-value pairs as ...any arguments
//   - Attribute-based methods (DebugAttr, InfoAttr, WarnAttr, ErrorAttr, LogAttr,
//     EventAttr, WithAttr) that accept strongly-typed log.KeyValue attributes
//
// The attribute-based methods provide better type safety and can offer better
// performance in some scenarios, particularly when used with WithAttr for
// pre-configured loggers.
type Logger struct {
	log.Logger
	attrs []log.KeyValue
}

// getCallerPackage returns the full package name of the caller.
// It walks the call stack to find the first caller outside of this package.
func getCallerPackage() string {
	// Start from frame 2 to skip getCallerPackage itself and New function.
	for i := 2; ; i++ {
		pc, _, _, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		name := fn.Name()
		pkg := extractPackageFromFuncName(name)

		// Skip empty packages.
		if pkg != "" {
			return pkg
		}
	}

	return "unknown"
}

// extractPackageFromFuncName extracts the package name from a full function name.
// Function names look like: "package/path.Function" or "package/path.(*Type).Method".
func extractPackageFromFuncName(funcName string) string {
	// Strategy: find the last dot before any parentheses
	// For "pkg.Function" -> "pkg"
	// For "pkg.(*Type).Method" -> "pkg" (dot before parentheses)

	// First, find the position of the first opening parenthesis (if any)
	parenPos := -1
	for i, r := range funcName {
		if r == '(' {
			parenPos = i
			break
		}
	}

	// Look for the last dot before the parenthesis (or in the entire string if no parenthesis)
	searchEnd := len(funcName)
	if parenPos >= 0 {
		searchEnd = parenPos
	}

	lastDot := -1
	for i := 0; i < searchEnd; i++ {
		if funcName[i] == '.' {
			lastDot = i
		}
	}

	if lastDot >= 0 {
		return funcName[:lastDot]
	}

	// No dot found
	return ""
}

// New creates a new Logger with the provided options.
// If options.Provider is nil, the global LoggerProvider is used.
// If options.Name is empty, the caller's full package name is automatically detected.
func New(options Options) *Logger {
	provider := options.Provider
	if provider == nil {
		provider = global.GetLoggerProvider()
	}

	// Use caller's package name if Name is not provided
	name := options.Name
	if name == "" {
		name = getCallerPackage()
	}

	// Create logger options
	var loggerOptions []log.LoggerOption
	if options.Version != "" {
		loggerOptions = append(loggerOptions, log.WithInstrumentationVersion(options.Version))
	}
	if options.Attributes.Len() > 0 {
		// TODO: Replace log.WithInstrumentationAttributes with log.WithInstrumentationAttributesSet when available
		loggerOptions = append(loggerOptions, log.WithInstrumentationAttributes(options.Attributes.ToSlice()...))
	}

	// Create the underlying log.Logger
	otelLogger := provider.Logger(name, loggerOptions...)
	return &Logger{
		Logger: otelLogger,
	}
}

// DebugEnabled reports whether the logger emits debug-level log records.
func (l *Logger) DebugEnabled(ctx context.Context) bool {
	return l.Enabled(ctx, log.EnabledParameters{
		Severity: log.SeverityDebug,
	})
}

// InfoEnabled reports whether the logger emits info-level log records.
func (l *Logger) InfoEnabled(ctx context.Context) bool {
	return l.Enabled(ctx, log.EnabledParameters{
		Severity: log.SeverityInfo,
	})
}

// WarnEnabled reports whether the logger emits warn-level log records.
func (l *Logger) WarnEnabled(ctx context.Context) bool {
	return l.Enabled(ctx, log.EnabledParameters{
		Severity: log.SeverityWarn,
	})
}

// ErrorEnabled reports whether the logger emits error-level log records.
func (l *Logger) ErrorEnabled(ctx context.Context) bool {
	return l.Enabled(ctx, log.EnabledParameters{
		Severity: log.SeverityError,
	})
}

// Debug logs a debug message with optional key-value pairs.
func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, log.SeverityDebug, msg, args...)
}

// Info logs an info message with optional key-value pairs.
func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, log.SeverityInfo, msg, args...)
}

// Warn logs a warning message with optional key-value pairs.
func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, log.SeverityWarn, msg, args...)
}

// Error logs an error message with optional key-value pairs.
func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, log.SeverityError, msg, args...)
}

// Log logs a message at the specified level with optional key-value pairs.
func (l *Logger) Log(ctx context.Context, level log.Severity, msg string, args ...any) {
	l.log(ctx, level, msg, args...)
}

// Event logs an event with the specified name and optional key-value pairs.
func (l *Logger) Event(ctx context.Context, name string, args ...any) {
	var record log.Record
	record.SetEventName(name)
	record.SetTimestamp(time.Now())
	record.SetSeverity(log.SeverityInfo)

	l.addAttributes(&record, args)
	l.Emit(ctx, record)
}

// DebugAttr logs a debug message with the provided attributes.
func (l *Logger) DebugAttr(ctx context.Context, msg string, attrs ...log.KeyValue) {
	l.logAttr(ctx, log.SeverityDebug, msg, attrs...)
}

// InfoAttr logs an info message with the provided attributes.
func (l *Logger) InfoAttr(ctx context.Context, msg string, attrs ...log.KeyValue) {
	l.logAttr(ctx, log.SeverityInfo, msg, attrs...)
}

// WarnAttr logs a warning message with the provided attributes.
func (l *Logger) WarnAttr(ctx context.Context, msg string, attrs ...log.KeyValue) {
	l.logAttr(ctx, log.SeverityWarn, msg, attrs...)
}

// ErrorAttr logs an error message with the provided attributes.
func (l *Logger) ErrorAttr(ctx context.Context, msg string, attrs ...log.KeyValue) {
	l.logAttr(ctx, log.SeverityError, msg, attrs...)
}

// LogAttr logs a message at the specified level with the provided attributes.
func (l *Logger) LogAttr(ctx context.Context, level log.Severity, msg string, attrs ...log.KeyValue) {
	l.logAttr(ctx, level, msg, attrs...)
}

// EventAttr logs an event with the specified name and the provided attributes.
func (l *Logger) EventAttr(ctx context.Context, name string, attrs ...log.KeyValue) {
	var record log.Record
	record.SetEventName(name)
	record.SetTimestamp(time.Now())
	record.SetSeverity(log.SeverityInfo)

	l.addKeyValueAttributes(&record, attrs)
	l.Emit(ctx, record)
}

// WithAttr returns a new Logger that includes the given attributes in all log records.
func (l *Logger) WithAttr(attrs ...log.KeyValue) *Logger {
	// Combine existing attrs with new attrs
	combinedAttrs := make([]log.KeyValue, 0, len(l.attrs)+len(attrs))
	combinedAttrs = append(combinedAttrs, l.attrs...)
	combinedAttrs = append(combinedAttrs, attrs...)

	return &Logger{
		Logger: l.Logger,
		attrs:  combinedAttrs,
	}
}

// With returns a new Logger that includes the given attributes in all log records.
func (l *Logger) With(args ...any) *Logger {
	// Convert args to KeyValue attributes
	newAttrs := convertArgsToKeyValues(args)

	// Combine existing attrs with new attrs
	combinedAttrs := make([]log.KeyValue, 0, len(l.attrs)+len(newAttrs))
	combinedAttrs = append(combinedAttrs, l.attrs...)
	combinedAttrs = append(combinedAttrs, newAttrs...)

	return &Logger{
		Logger: l.Logger,
		attrs:  combinedAttrs,
	}
}

// log is the internal logging method that handles the common logging logic.
func (l *Logger) log(ctx context.Context, level log.Severity, msg string, args ...any) {
	var record log.Record
	record.SetBody(log.StringValue(msg))
	record.SetTimestamp(time.Now())
	record.SetSeverity(level)

	l.addAttributes(&record, args)
	l.Emit(ctx, record)
}

// addAttributes adds key-value pairs to the record.
// It supports the alternating key-value syntax like slog.
func (l *Logger) addAttributes(record *log.Record, args []any) {
	// Add pre-configured attributes first
	record.AddAttributes(l.attrs...)
	// Then add call-specific attributes
	addArgsAsAttributes(record, args)
}

// convertArgsToKeyValues converts alternating key-value arguments to log.KeyValue slice.
func convertArgsToKeyValues(args []any) []log.KeyValue {
	keyValues := make([]log.KeyValue, 0, len(args)/2+1)
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			// Odd number of arguments, add the key with empty value
			if key, ok := args[i].(string); ok {
				keyValues = append(keyValues, log.String(key, ""))
			}
			break
		}

		key, ok := args[i].(string)
		if !ok {
			continue
		}

		value := args[i+1]
		kv := log.KeyValue{
			Key:   key,
			Value: convertValue(value),
		}
		keyValues = append(keyValues, kv)
	}
	return keyValues
}

// addArgsAsAttributes processes alternating key-value arguments and adds them to the record.
func addArgsAsAttributes(record *log.Record, args []any) {
	keyValues := convertArgsToKeyValues(args)
	record.AddAttributes(keyValues...)
}

// logAttr is the internal logging method that handles logging with log.KeyValue attributes.
func (l *Logger) logAttr(ctx context.Context, level log.Severity, msg string, attrs ...log.KeyValue) {
	var record log.Record
	record.SetBody(log.StringValue(msg))
	record.SetTimestamp(time.Now())
	record.SetSeverity(level)

	l.addKeyValueAttributes(&record, attrs)
	l.Emit(ctx, record)
}

// addKeyValueAttributes adds log.KeyValue attributes to the record.
func (l *Logger) addKeyValueAttributes(record *log.Record, attrs []log.KeyValue) {
	// Add pre-configured attributes first
	record.AddAttributes(l.attrs...)
	// Then add call-specific attributes
	record.AddAttributes(attrs...)
}
