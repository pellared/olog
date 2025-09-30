# Changelog

All notable changes to this library are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this library adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html),
as well as to [Module version numbering](https://go.dev/doc/modules/version-numbers).

## [Unreleased](https://github.com/pellared/olog/compare/v0.0.2...HEAD)

### Added

- `Logger.TraceEvent(ctx context.Context, name string, args ...any)` that logs a trace-level event with the specified name and optional key-value pairs.
- `Logger.DebugEvent(ctx context.Context, name string, args ...any)` that logs a debug-level event with the specified name and optional key-value pairs.
- `Logger.InfoEvent(ctx context.Context, name string, args ...any)` that logs an info-level event with the specified name and optional key-value pairs.
- `Logger.WarnEvent(ctx context.Context, name string, args ...any)` that logs a warn-level event with the specified name and optional key-value pairs.
- `Logger.ErrorEvent(ctx context.Context, name string, args ...any)` that logs an error-level event with the specified name and optional key-value pairs.
- `Logger.TraceEventAttr(ctx context.Context, name string, attrs ...log.KeyValue)` that logs a trace-level event with the specified name and the provided attributes.
- `Logger.DebugEventAttr(ctx context.Context, name string, attrs ...log.KeyValue)` that logs a debug-level event with the specified name and the provided attributes.
- `Logger.InfoEventAttr(ctx context.Context, name string, attrs ...log.KeyValue)` that logs an info-level event with the specified name and the provided attributes.
- `Logger.WarnEventAttr(ctx context.Context, name string, attrs ...log.KeyValue)` that logs a warn-level event with the specified name and the provided attributes.
- `Logger.ErrorEventAttr(ctx context.Context, name string, attrs ...log.KeyValue)` that logs an error-level event with the specified name and the provided attributes.
- `Logger.EventEnabled(ctx context.Context) bool` that reports whether the logger emits event log records.

### Changed

- **BREAKING:** `Logger.Event(ctx context.Context, name string, args ...any)` now requires a `level log.Severity` parameter: `Logger.Event(ctx context.Context, level log.Severity, name string, args ...any)`.
- **BREAKING:** `Logger.EventAttr(ctx context.Context, name string, attrs ...log.KeyValue)` now requires a `level log.Severity` parameter: `Logger.EventAttr(ctx context.Context, level log.Severity, name string, attrs ...log.KeyValue)`.

## [0.0.2](https://github.com/pellared/olog/releases/tag/v0.0.2) - 2025-09-26

### Changed

- Minimum required Go version is now Go 1.24 instead of Go 1.25.1.

## [0.0.1](https://github.com/pellared/olog/releases/tag/v0.0.1) - 2025-09-26

### Added

- `New(options Options) *Logger` that creates a new Logger with the provided options.
- `Logger` type that provides an ergonomic frontend API for OpenTelemetry structured logging with embedded `log.Logger`.
- `Options` type that contains configuration options for creating a Logger.
- `Logger.Trace(ctx context.Context, msg string, args ...any)` that logs a trace message with optional key-value pairs.
- `Logger.TraceAttr(ctx context.Context, msg string, attrs ...log.KeyValue)` that logs a trace message with the provided attributes.
- `Logger.TraceEnabled(ctx context.Context) bool` that reports whether the logger emits trace-level log records.
- `Logger.Debug(ctx context.Context, msg string, args ...any)` that logs a debug message with optional key-value pairs.
- `Logger.DebugAttr(ctx context.Context, msg string, attrs ...log.KeyValue)` that logs a debug message with the provided attributes.
- `Logger.DebugEnabled(ctx context.Context) bool` that reports whether the logger emits debug-level log records.
- `Logger.Info(ctx context.Context, msg string, args ...any)` that logs an info message with optional key-value pairs.
- `Logger.InfoAttr(ctx context.Context, msg string, attrs ...log.KeyValue)` that logs an info message with the provided attributes.
- `Logger.InfoEnabled(ctx context.Context) bool` that reports whether the logger emits info-level log records.
- `Logger.Warn(ctx context.Context, msg string, args ...any)` that logs a warning message with optional key-value pairs.
- `Logger.WarnAttr(ctx context.Context, msg string, attrs ...log.KeyValue)` that logs a warning message with the provided attributes.
- `Logger.WarnEnabled(ctx context.Context) bool` that reports whether the logger emits warn-level log records.
- `Logger.Error(ctx context.Context, msg string, args ...any)` that logs an error message with optional key-value pairs.
- `Logger.ErrorAttr(ctx context.Context, msg string, attrs ...log.KeyValue)` that logs an error message with the provided attributes.
- `Logger.ErrorEnabled(ctx context.Context) bool` that reports whether the logger emits error-level log records.
- `Logger.Log(ctx context.Context, level log.Severity, msg string, args ...any)` that logs a message at the specified level with optional key-value pairs.
- `Logger.LogAttr(ctx context.Context, level log.Severity, msg string, attrs ...log.KeyValue)` that logs a message at the specified level with the provided attributes.
- `Logger.Event(ctx context.Context, name string, args ...any)` that logs an event with the specified name and optional key-value pairs.
- `Logger.EventAttr(ctx context.Context, name string, attrs ...log.KeyValue)` that logs an event with the specified name and the provided attributes.
- `Logger.With(args ...any) *Logger` that returns a new Logger that includes the given attributes in all log records.
- `Logger.WithAttr(attrs ...log.KeyValue) *Logger` that returns a new Logger that includes the given attributes in all log records.

<!-- markdownlint-configure-file
{
  "MD024": {
    "siblings_only": true
  }
}
-->
