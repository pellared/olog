package olog_test

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"

	"github.com/pellared/olog"
)

func ExampleLogger_basic() {
	ctx := context.Background()

	// Create a logger instance
	logger := olog.New(olog.Options{})

	// Use the logger for basic logging
	logger.Info(ctx, "application started", "version", "1.0.0", "port", 8080)
	logger.Warn(ctx, "deprecated feature used", "feature", "old-api")
	logger.Error(ctx, "failed to connect", "host", "db.example.com", "error", "connection timeout")
	logger.Debug(ctx, "processing request", "method", "GET", "path", "/api/users")

	// Check if logging is enabled before expensive operations
	if logger.DebugEnabled(ctx) {
		expensiveData := computeExpensiveDebugInfo()
		logger.Debug(ctx, "debug info", "data", expensiveData)
	}
}

func ExampleLogger_with() {
	ctx := context.Background()

	// Create a logger and add common attributes
	baseLogger := olog.New(olog.Options{Provider: global.GetLoggerProvider(), Name: "example"})
	logger := baseLogger.With("service", "user-service", "version", "2.1.0")

	// All log records from this logger will include the common attributes
	logger.Info(ctx, "user created", "user_id", 12345, "email", "user@example.com")
	logger.Warn(ctx, "user login failed", "user_id", 12345, "reason", "invalid password")

	// Chain with additional attributes
	requestLogger := logger.With("request_id", "req-789", "ip", "192.168.1.100")
	requestLogger.Info(ctx, "processing request", "endpoint", "/api/users/12345")
	requestLogger.Error(ctx, "request failed", "status", 500, "duration_ms", 1234)
}

func ExampleLogger_events() {
	ctx := context.Background()

	// Create a logger for events
	logger := olog.New(olog.Options{Provider: global.GetLoggerProvider(), Name: "example"})

	// Log structured events following semantic conventions
	logger.Event(ctx, "user.login",
		"user.id", "12345",
		"user.email", "user@example.com",
		"session.id", "sess-abc123",
		"client.ip", "192.168.1.100")

	logger.Event(ctx, "payment.processed",
		"payment.id", "pay-xyz789",
		"payment.amount", 99.99,
		"payment.currency", "USD",
		"user.id", "12345")

	// Events with additional attributes
	logger.Event(ctx, "order.created",
		"order.id", "order-456",
		"order.total", 149.99,
		"customer.id", "cust-789",
		"domain", "ecommerce")
}

func ExampleLogger_performance() {
	ctx := context.Background()

	// Create a base logger
	logger := olog.New(olog.Options{Provider: global.GetLoggerProvider(), Name: "example"})

	// Check if logging is enabled to avoid expensive operations
	if logger.DebugEnabled(ctx) {
		// Only compute expensive debug information if debug logging is enabled
		debugData := computeExpensiveDebugInfo()
		logger.Debug(ctx, "detailed debug information", "data", debugData)
	}

	// Pre-configure logger with common attributes for better performance
	requestLogger := logger.With(
		"service", "api-server",
		"version", "1.2.3",
		"request_id", generateRequestID(),
	)

	// Use the pre-configured logger for all request-scoped logging
	requestLogger.Info(ctx, "request started", "method", "GET", "path", "/api/users")
	requestLogger.Info(ctx, "database query", "table", "users", "duration_ms", 23)
	requestLogger.Info(ctx, "request completed", "status", 200, "total_duration_ms", 145)
}

func computeExpensiveDebugInfo() string {
	// Simulate expensive computation
	return "expensive debug data"
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

	// All log records will include the pre-configured attributes
	logger.Info(ctx, "service started", "status", "ready")
	logger.Warn(ctx, "high memory usage", "memory_percent", 85.5)
	logger.Error(ctx, "database connection failed", "retry_count", 3)
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

	logger.Info(ctx, "using global logger provider", "initialized", true)
}

func ExampleNew_minimal() {
	ctx := context.Background()

	// Minimal configuration - only name is required
	logger := olog.New(olog.Options{
		Name: "minimal-logger",
	})

	logger.Info(ctx, "minimal logger example")
}

func ExampleLogger_withAttributes() {
	ctx := context.Background()

	// Import the log package for KeyValue construction
	// import "go.opentelemetry.io/otel/log"

	// Create a logger instance
	logger := olog.New(olog.Options{Provider: global.GetLoggerProvider(), Name: "example"})

	// Using the new attribute-based methods for type-safe logging
	logger.InfoAttr(ctx, "user logged in",
		log.String("user.id", "12345"),
		log.String("user.email", "user@example.com"),
		log.Int64("session.duration", 3600),
		log.Bool("first_login", false))

	logger.WarnAttr(ctx, "rate limit exceeded",
		log.String("client.ip", "192.168.1.100"),
		log.Int64("requests_per_minute", 150),
		log.Int64("limit", 100))

	logger.ErrorAttr(ctx, "database connection failed",
		log.String("database.host", "db.example.com"),
		log.Int64("database.port", 5432),
		log.String("error.type", "connection_timeout"))

	// Use LogAttr for custom severity levels
	logger.LogAttr(ctx, log.SeverityWarn2, "custom warning",
		log.String("component", "cache"),
		log.Float64("memory_usage_percent", 85.5))
}

func ExampleLogger_eventAttr() {
	ctx := context.Background()

	// Create a logger for structured events
	logger := olog.New(olog.Options{Provider: global.GetLoggerProvider(), Name: "example"})

	// Log events using the attribute-based method for better type safety
	logger.EventAttr(ctx, "user.signup",
		log.String("user.id", "user-789"),
		log.String("user.email", "newuser@example.com"),
		log.String("signup.method", "email"),
		log.Bool("email.verified", false))

	logger.EventAttr(ctx, "payment.processed",
		log.String("payment.id", "pay-abc123"),
		log.Float64("payment.amount", 49.99),
		log.String("payment.currency", "USD"),
		log.String("payment.method", "credit_card"))

	logger.EventAttr(ctx, "file.uploaded",
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

	// All subsequent logs will include the service attributes
	serviceLogger.InfoAttr(ctx, "service started",
		log.Int64("port", 8080),
		log.String("build", "abc1234"))

	// Chain additional attributes for request-scoped logging
	requestLogger := serviceLogger.WithAttr(
		log.String("request.id", "req-789"),
		log.String("user.id", "user-456"))

	requestLogger.InfoAttr(ctx, "processing request",
		log.String("http.method", "POST"),
		log.String("http.route", "/api/users"))

	requestLogger.ErrorAttr(ctx, "validation failed",
		log.String("field", "email"),
		log.String("error", "invalid format"))

	// Mix WithAttr and With methods
	mixedLogger := requestLogger.With("trace.id", "trace-xyz").WithAttr(log.Bool("debug.enabled", true))
	mixedLogger.DebugAttr(ctx, "detailed processing info",
		log.Int64("processing.step", 3),
		log.Float64("processing.duration_ms", 12.5))
}
