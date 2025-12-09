package olog_test

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"

	"github.com/pellared/olog"
)

func ExampleLogger_events() {
	ctx := context.Background()

	// Create a logger for events
	logger := olog.New(olog.Options{Provider: global.GetLoggerProvider(), Name: "example"})

	// Log structured events at different severity levels
	logger.InfoEvent(ctx, "user.login",
		"user.id", "12345",
		"user.email", "user@example.com",
		"session.id", "sess-abc123",
		"client.ip", "192.168.1.100")

	logger.WarnEvent(ctx, "rate.limit.approached",
		"client.ip", "192.168.1.100",
		"requests_per_minute", 85,
		"limit", 100)

	logger.ErrorEvent(ctx, "payment.failed",
		"payment.id", "pay-xyz789",
		"payment.amount", 99.99,
		"payment.currency", "USD",
		"user.id", "12345",
		"error", "insufficient_funds")

	// Check if debug event logging is enabled before expensive operations
	if logger.DebugEventEnabled(ctx, "debug.session.details") {
		logger.DebugEvent(ctx, "debug.session.details",
			"session.data", computeSessionDebugInfo(),
			"trace.id", "trace-abc123")
	}
}

func ExampleLogger_eventsWithContext() {
	ctx := context.Background()

	// Create a base logger
	logger := olog.New(olog.Options{Provider: global.GetLoggerProvider(), Name: "example"})

	// Pre-configure logger with common attributes for better performance
	requestLogger := logger.With(
		"service", "api-server",
		"version", "1.2.3",
		"request_id", generateRequestID(),
	)

	// Use the pre-configured logger for all request-scoped events
	requestLogger.InfoEvent(ctx, "request.started", "method", "GET", "path", "/api/users")
	requestLogger.InfoEvent(ctx, "database.query", "table", "users", "duration_ms", 23)
	requestLogger.InfoEvent(ctx, "request.completed", "status", 200, "total_duration_ms", 145)
}

func computeSessionDebugInfo() string {
	// Simulate expensive session debug data computation
	return "session debug information"
}

func generateRequestID() string {
	return "req-12345"
}

func ExampleNew_withOptions() {
	ctx := context.Background()

	// Create a logger using the new Options API
	logger := olog.New(olog.Options{
		Name:    "my-service",
		Version: "1.2.3",
		Attributes: attribute.NewSet(
			attribute.String("service.name", "user-service"),
			attribute.String("deployment.environment", "production"),
			attribute.Int("service.port", 8080),
		),
	})

	// All event records will include the pre-configured attributes
	logger.InfoEvent(ctx, "service.started", "status", "ready")
	logger.WarnEvent(ctx, "resource.high_usage", "memory_percent", 85.5)
	logger.ErrorEvent(ctx, "database.connection_failed", "retry_count", 3)
}

func ExampleNew_withGlobalProvider() {
	ctx := context.Background()

	// Create a logger that uses the global provider (provider is nil)
	logger := olog.New(olog.Options{
		Name: "global-logger",
		Attributes: attribute.NewSet(
			attribute.String("component", "authentication"),
			attribute.String("version", "2.0.0"),
		),
	})

	logger.InfoEvent(ctx, "logger.initialized", "uses_global_provider", true)
}

func ExampleNew_minimal() {
	ctx := context.Background()

	// Minimal configuration.
	// The logger name is the caller's full package name.
	logger := olog.New(olog.Options{})

	logger.InfoEvent(ctx, "minimal.example")
}

func ExampleLogger_eventAttr() {
	ctx := context.Background()

	// Create a logger for structured events
	logger := olog.New(olog.Options{Provider: global.GetLoggerProvider(), Name: "example"})

	// Log events at different severity levels using the attribute-based methods
	logger.InfoEventAttr(ctx, "user.signup",
		log.String("user.id", "user-789"),
		log.String("user.email", "newuser@example.com"),
		log.String("signup.method", "email"),
		log.Bool("email.verified", false))

	logger.WarnEventAttr(ctx, "api.deprecated.usage",
		log.String("api.endpoint", "/v1/users"),
		log.String("client.id", "client-123"),
		log.String("replacement", "/v2/users"))

	logger.ErrorEventAttr(ctx, "payment.failed",
		log.String("payment.id", "pay-abc123"),
		log.Float64("payment.amount", 49.99),
		log.String("payment.currency", "USD"),
		log.String("payment.method", "credit_card"),
		log.String("error.code", "card_declined"))

	logger.DebugEventAttr(ctx, "file.uploaded",
		log.String("file.id", "file-456"),
		log.String("file.name", "document.pdf"),
		log.Int64("file.size_bytes", 2048576),
		log.String("user.id", "user-123"))
}

func ExampleLogger_withAttr() {
	ctx := context.Background()

	// Create a base logger
	baseLogger := olog.New(olog.Options{Provider: global.GetLoggerProvider(), Name: "example"})

	// Create a logger with common attributes using WithAttr
	serviceLogger := baseLogger.WithAttr(
		log.String("service.name", "user-service"),
		log.String("service.version", "2.1.0"),
		log.String("deployment.environment", "production"))

	// All subsequent events will include the service attributes
	serviceLogger.InfoEventAttr(ctx, "service.started",
		log.Int64("port", 8080),
		log.String("build", "abc1234"))

	// Chain additional attributes for request-scoped logging
	requestLogger := serviceLogger.WithAttr(
		log.String("request.id", "req-789"),
		log.String("user.id", "user-456"))

	requestLogger.InfoEventAttr(ctx, "request.processing",
		log.String("http.method", "POST"),
		log.String("http.route", "/api/users"))

	requestLogger.ErrorEventAttr(ctx, "validation.failed",
		log.String("field", "email"),
		log.String("error", "invalid format"))

	// Mix WithAttr and With methods
	mixedLogger := requestLogger.With("trace.id", "trace-xyz").WithAttr(log.Bool("debug.enabled", true))
	mixedLogger.DebugEventAttr(ctx, "processing.detail",
		log.Int64("processing.step", 3),
		log.Float64("processing.duration_ms", 12.5))
}
