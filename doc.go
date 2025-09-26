/*
Package olog provides an ergonomic OpenTelemetry Logging Facade.

This package addresses the usability concerns with the OpenTelemetry Logs API
by providing a user-friendly frontend interface while using the OpenTelemetry
Logs API as the backend. It offers simple methods similar to popular logging
libraries while maintaining full compatibility with OpenTelemetry's structured
logging capabilities.

# Basic Usage

The simplest way to use olog is by creating a logger instance:

	import (
		"context"
		"go.opentelemetry.io/otel/log/global"
		"github.com/pellared/olog"
	)

	ctx := context.Background()
	logger := olog.New(olog.Options{
		Provider: global.GetLoggerProvider(),
		Name:     "myapp",
	})

	logger.Trace(ctx, "detailed tracing", "trace_id", "abc123")
	logger.Info(ctx, "application started", "version", "1.0.0", "port", 8080)
	logger.Warn(ctx, "deprecated feature used", "feature", "old-api")
	logger.Error(ctx, "failed to connect", "host", "db.example.com")
	logger.Debug(ctx, "processing request", "method", "GET", "path", "/api/users")

	// Check if logging is enabled before expensive operations
	if logger.DebugEnabled(ctx) {
		expensiveData := computeExpensiveDebugInfo()
		logger.Debug(ctx, "debug info", "data", expensiveData)
	}

# Logger Composition

Use With to create loggers with common attributes:

	serviceLogger := logger.With("service", "user-service", "version", "2.1.0")
	serviceLogger.Info(ctx, "user created", "user_id", 12345)

	requestLogger := serviceLogger.With("request_id", "req-789")
	requestLogger.Info(ctx, "processing request", "endpoint", "/api/users")

Use structured attributes to organize your logs:

	httpLogger := logger.With("component", "http")
	httpLogger.Info(ctx, "request", "method", "POST", "status", 201)
	// Logs with component="http" and the specified attributes

# Event Logging

Log structured events following semantic conventions:

	logger.Event(ctx, "user.login",
		"user.id", "12345",
		"user.email", "user@example.com",
		"session.id", "sess-abc123")

# Performance

olog is designed with performance in mind:

  - Use TraceEnabled, DebugEnabled, InfoEnabled, WarnEnabled, and ErrorEnabled checks to avoid expensive operations when logging is disabled
  - Logger composition with With pre-processes common attributes
  - Direct integration with OpenTelemetry Logs API avoids unnecessary conversions

# Design Goals

This package is designed to provide:

 1. Simple, ergonomic API similar to popular logging libraries
 2. Performance-oriented design with efficient enabled checks
 3. Full compatibility with OpenTelemetry Logs API and ecosystem
 4. Support for structured logging with key-value pairs
 5. Logger composition for better code organization and performance
 6. Event logging capabilities for semantic events
*/
package olog // import "github.com/pellared/olog"
