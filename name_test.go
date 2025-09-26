// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package olog_test

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/logtest"

	"github.com/pellared/olog"
)

func TestLogger_DefaultName(t *testing.T) {
	recorder := logtest.NewRecorder()
	logger := olog.New(olog.Options{
		Provider: recorder,
	})

	ctx := t.Context()

	// Test Info logging
	logger.Info(ctx, "test info message", "key1", "value1", "key2", 42)

	// Verify using logtest.AssertEqual with Recording
	want := logtest.Recording{
		logtest.Scope{
			Name: "github.com/pellared/olog_test",
		}: {
			logtest.Record{
				Context:  ctx,
				Severity: log.SeverityInfo,
				Body:     log.StringValue("test info message"),
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
