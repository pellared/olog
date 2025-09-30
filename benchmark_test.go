package olog

import (
	"testing"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
)

func BenchmarkLogger_Info(b *testing.B) {
	logger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	ctx := b.Context()

	for i := 0; b.Loop(); i++ {
		logger.Info(ctx, "benchmark message", "iteration", i, "data", "test")
	}
}

func BenchmarkLogger_InfoWithEnabled(b *testing.B) {
	logger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	ctx := b.Context()

	for i := 0; b.Loop(); i++ {
		if logger.InfoEnabled(ctx) {
			logger.Info(ctx, "benchmark message", "iteration", i, "data", "test")
		}
	}
}

func BenchmarkLogger_With(b *testing.B) {
	baseLogger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	logger := baseLogger.With("service", "test", "version", "1.0.0")
	ctx := b.Context()

	for i := 0; b.Loop(); i++ {
		logger.Info(ctx, "benchmark message", "iteration", i)
	}
}

func BenchmarkLogger_Event(b *testing.B) {
	logger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	ctx := b.Context()

	for i := 0; b.Loop(); i++ {
		logger.Event(ctx, log.SeverityInfo, "test.event", "iteration", i, "data", "test")
	}
}

// Benchmark the new attribute-based methods

func BenchmarkLogger_InfoAttr(b *testing.B) {
	logger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	ctx := b.Context()

	for i := 0; b.Loop(); i++ {
		logger.InfoAttr(ctx, "benchmark message", log.Int64("iteration", int64(i)), log.String("data", "test"))
	}
}

func BenchmarkLogger_InfoAttrWithEnabled(b *testing.B) {
	logger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	ctx := b.Context()

	for i := 0; b.Loop(); i++ {
		if logger.InfoEnabled(ctx) {
			logger.InfoAttr(ctx, "benchmark message", log.Int64("iteration", int64(i)), log.String("data", "test"))
		}
	}
}

func BenchmarkLogger_WithAttr(b *testing.B) {
	baseLogger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	logger := baseLogger.WithAttr(log.String("service", "test"), log.String("version", "1.0.0"))
	ctx := b.Context()

	for i := 0; b.Loop(); i++ {
		logger.InfoAttr(ctx, "benchmark message", log.Int64("iteration", int64(i)))
	}
}

func BenchmarkLogger_EventAttr(b *testing.B) {
	logger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	ctx := b.Context()

	for i := 0; b.Loop(); i++ {
		logger.EventAttr(ctx, log.SeverityInfo, "test.event", log.Int64("iteration", int64(i)), log.String("data", "test"))
	}
}

// Comparison benchmarks between args and attr methods

func BenchmarkLogger_InfoComparison(b *testing.B) {
	logger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	ctx := b.Context()

	b.Run("Args", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			logger.Info(ctx, "benchmark message", "iteration", i, "data", "test", "bool_flag", true, "score", 98.5)
		}
	})

	b.Run("Attr", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			logger.InfoAttr(ctx, "benchmark message",
				log.Int64("iteration", int64(i)),
				log.String("data", "test"),
				log.Bool("bool_flag", true),
				log.Float64("score", 98.5))
		}
	})
}

func BenchmarkLogger_WithComparison(b *testing.B) {
	baseLogger := New(Options{Provider: noop.NewLoggerProvider(), Name: "bench"})
	ctx := b.Context()

	b.Run("WithArgs", func(b *testing.B) {
		logger := baseLogger.With("service", "test", "version", "1.0.0", "environment", "prod")
		for i := 0; b.Loop(); i++ {
			logger.Info(ctx, "benchmark message", "iteration", i)
		}
	})

	b.Run("WithAttr", func(b *testing.B) {
		logger := baseLogger.WithAttr(
			log.String("service", "test"),
			log.String("version", "1.0.0"),
			log.String("environment", "prod"))
		for i := 0; b.Loop(); i++ {
			logger.InfoAttr(ctx, "benchmark message", log.Int64("iteration", int64(i)))
		}
	})
}
