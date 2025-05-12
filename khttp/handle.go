// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package khttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	kerr "github.com/kopexa-grc/common/errors"
	"github.com/kopexa-grc/common/khttp/metric"
	"github.com/kopexa-grc/common/wellknown"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var (
	totalOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "operation_total",
			Namespace: wellknown.PrometheusNamespaceKopexa,
			Subsystem: "http",
			Help:      "Operations total per method and if they had an error or not",
		},
		[]string{"err", "name"},
	)

	operationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "operation_duration_seconds",
			Namespace: wellknown.PrometheusNamespaceKopexa,
			Subsystem: "http",
			Help:      "A histogram of latencies for operations.",
			Buckets:   []float64{0.001, 0.01, 0.05, 0.1, .25, .5, 1, 1.5, 2, 5, 10, 30},
		},
		[]string{"err", "name"},
	)
)

func init() {
	metric.GlobalRegistry.MustRegister(totalOperations)
	metric.GlobalRegistry.MustRegister(operationDuration)
}

const StatusClientClosedRequest = 499

func Handle(name string, w http.ResponseWriter, r *http.Request, fn func(ctx context.Context) error) {
	ctx := r.Context()

	// we make sure no body is ever leaking
	defer r.Body.Close() // nolint: errcheck

	span, ctx := opentracing.StartSpanFromContext(ctx, name)
	labels := prometheus.Labels{
		"name": name,
		"err":  "",
	}
	start := time.Now()

	err := fn(ctx)
	if err != nil {
		span.LogKV("error", err)
		labels["err"] = err.Error()
	}

	dur := time.Since(start)

	// Note: in case of a panic these will not be completed
	span.Finish()
	totalOperations.With(labels).Inc()
	operationDuration.With(labels).Observe(dur.Seconds())

	if err != nil {
		switch {
		case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
			select {
			case <-ctx.Done():
				if ctx.Err() != nil {
					log.Error().Err(err).Msg("request failed, report err to client")
					w.WriteHeader(StatusClientClosedRequest)

					return
				}
			default:
			}
		case kerr.IsError(err):
			WriteSimpleErr(w, err)
		default:
			log.Error().Err(err).Msg("request failed")
			WriteSimpleErr(w, err)
		}
	}
}

func ReadJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
