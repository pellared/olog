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
		"go.opentelemetry.io/otel/log"
		"go.opentelemetry.io/otel/log/global"
		"github.com/pellared/olog"
	)

	ctx := context.Background()
	logger := olog.New(olog.Options{
		Provider: global.GetLoggerProvider(),
		Name:     "myapp",
	})

	logger.TraceAttr(ctx, "detailed tracing", log.String("trace_id", "abc123"))
	logger.InfoAttr(ctx, "application started",
		log.String("version", "1.0.0"),
		log.Int("port", 8080))
	logger.WarnAttr(ctx, "deprecated feature used", log.String("feature", "old-api"))
	logger.ErrorAttr(ctx, "failed to connect", log.String("host", "db.example.com"))
	logger.DebugAttr(ctx, "processing request",
		log.String("method", "GET"),
		log.String("path", "/api/users"))

	// Check if logging is enabled before expensive operations
	if logger.DebugEnabled(ctx) {
		expensiveData := computeExpensiveDebugInfo()
		logger.DebugAttr(ctx, "debug info", log.String("data", expensiveData))
	}

# Logger Composition

Use WithAttr to create loggers with common attributes:

	serviceLogger := logger.WithAttr(
		log.String("service", "user-service"),
		log.String("version", "2.1.0"))
	serviceLogger.InfoAttr(ctx, "user created", log.Int("user_id", 12345))

	requestLogger := serviceLogger.WithAttr(log.String("request_id", "req-789"))
	requestLogger.InfoAttr(ctx, "processing request", log.String("endpoint", "/api/users"))

Use structured attributes to organize your logs:

	httpLogger := logger.WithAttr(log.String("component", "http"))
	httpLogger.InfoAttr(ctx, "request",
		log.String("method", "POST"),
		log.Int("status", 201))
	// Logs with component="http" and the specified attributes

# Event Logging

Log structured events at different severity levels following semantic conventions:

	// Log events at different levels
	logger.InfoEventAttr(ctx, "user.login",
		log.String("user.id", "12345"),
		log.String("user.email", "user@example.com"),
		log.String("session.id", "sess-abc123"))

	logger.WarnEventAttr(ctx, "rate.limit.exceeded",
		log.String("client.ip", "192.168.1.100"),
		log.Int("requests_per_minute", 150))

		logger.ErrorEventAttr(ctx, \"payment.failed\",\n\t\tlog.String(\"payment.id\", \"pay-123\"),\n\t\tlog.String(\"error\", \"insufficient_funds\"))\n\n\t// Use the generic Event method with custom severity\n\tlogger.Event(ctx, log.SeverityWarn2, \"custom.event\", \"key\", \"value\")

# Performance

olog is designed with performance in mind:

  - Use TraceEnabled, DebugEnabled, InfoEnabled, WarnEnabled, ErrorEnabled, and EventEnabled checks to avoid expensive operations when logging is disabled
  - Logger composition with WithAttr pre-processes common attributes
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
