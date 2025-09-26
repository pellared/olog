// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package olog

import (
	"testing"

	"go.opentelemetry.io/otel/log/noop"
)

func TestLogger_StructBehavior(t *testing.T) {
	logger := New(Options{Provider: noop.NewLoggerProvider(), Name: "test"})

	ctx := t.Context()

	// Test that the original logger is not modified when creating a new one with With()
	originalLogger := logger
	withLogger := logger.With("key1", "value1")

	// Should be different instances
	if originalLogger == withLogger {
		t.Error("With() should return a new logger instance")
	}

	// Original logger should not have attributes
	if len(logger.attrs) != 0 {
		t.Errorf("Original logger should have no attrs, got %d", len(logger.attrs))
	}

	// With logger should have attributes
	if len(withLogger.attrs) != 1 {
		t.Errorf("With logger should have 1 KeyValue attr, got %d", len(withLogger.attrs))
	}

	// Test chaining With calls
	chainedLogger := withLogger.With("key2", "value2")
	if len(chainedLogger.attrs) != 2 {
		t.Errorf("Chained logger should have 2 KeyValue attrs, got %d", len(chainedLogger.attrs))
	}

	// Test logging doesn't panic
	logger.Info(ctx, "test")
	withLogger.Info(ctx, "test with attrs")
	chainedLogger.Info(ctx, "test chained")
}

func TestLogger_AttributeHandling(t *testing.T) {
	logger := New(Options{Provider: noop.NewLoggerProvider(), Name: "test"})

	// Test that attributes are properly stored
	withLogger := logger.With("service", "api", "version", "1.0")
	if len(withLogger.attrs) != 2 {
		t.Errorf("Expected 2 KeyValue attrs, got %d", len(withLogger.attrs))
	}

	// Check the key-value pairs
	expectedKeys := []string{"service", "version"}
	expectedValues := []string{"api", "1.0"}

	if len(withLogger.attrs) != len(expectedKeys) {
		t.Fatalf("Attr length mismatch: expected %d, got %d", len(expectedKeys), len(withLogger.attrs))
	}

	for i, expectedKey := range expectedKeys {
		if withLogger.attrs[i].Key != expectedKey {
			t.Errorf("Attr[%d] key: expected %s, got %s", i, expectedKey, withLogger.attrs[i].Key)
		}
		if withLogger.attrs[i].Value.AsString() != expectedValues[i] {
			t.Errorf("Attr[%d] value: expected %s, got %s", i, expectedValues[i], withLogger.attrs[i].Value.AsString())
		}
	}
}
