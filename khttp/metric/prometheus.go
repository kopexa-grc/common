// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package metric

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kopexa-grc/common/wellknown"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Registry is a shallow wrapper around a prometheus Registry with some helper functions go get an HTTP Request wrapper and an Handler for this registry.
type Registry struct {
	*prometheus.Registry
}

// NewRegistry returns a new registry with some default collectors registered
func NewRegistry() *Registry {
	r := prometheus.NewRegistry()
	r.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewBuildInfoCollector(),
	)

	return &Registry{Registry: r}
}

// Handler returns a HTTP handler for this registry. Should be registered at "/metrics" with:
// router.Use("/metrics", registry.Handler())
func (r *Registry) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(r, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
}

// PrometheusMiddleware creates a Chi middleware and registers the metrics to the given prometheus registerer
func (r *Registry) PrometheusMiddleware() func(next http.Handler) http.Handler {
	namespace := wellknown.PrometheusNamespaceKopexa
	subsystem := "http"

	requestsInFlight := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      MetricInFlightRequests,
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      HelpInFlightRequests,
	})

	totalRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      MetricAPITotalRequests,
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      HelpAPITotalRequests,
		},
		[]string{LabelCode, LabelMethod},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      MetricRequestDuration,
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      HelpRequestDuration,
			Buckets:   DurationBuckets,
		},
		[]string{LabelHandler, LabelCode, LabelMethod},
	)

	timeToFirstHeader := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      MetricTimeToFirstHeader,
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      HelpTimeToFirstHeader,
			Buckets:   DurationBuckets,
		},
		[]string{LabelCode, LabelMethod},
	)

	responseSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      MetricResponseSize,
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      HelpResponseSize,
			Buckets:   SizeBuckets,
		},
		[]string{LabelHandler, LabelMethod, LabelCode},
	)

	requestSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      MetricRequestSize,
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      HelpRequestSize,
			Buckets:   SizeBuckets,
		},
		[]string{LabelHandler, LabelMethod, LabelCode},
	)

	r.MustRegister(
		requestsInFlight,
		totalRequests,
		requestDuration,
		timeToFirstHeader,
		responseSize,
		requestSize,
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			reqDuration := &delayingObserver{requestDuration.MustCurryWith(prometheus.Labels{LabelHandler: DefaultHandler})}
			reqSize := &delayingObserver{requestSize.MustCurryWith(prometheus.Labels{LabelHandler: DefaultHandler})}
			respSize := &delayingObserver{responseSize.MustCurryWith(prometheus.Labels{LabelHandler: DefaultHandler})}

			chain := promhttp.InstrumentHandlerInFlight(requestsInFlight,
				promhttp.InstrumentHandlerDuration(reqDuration,
					promhttp.InstrumentHandlerCounter(totalRequests,
						promhttp.InstrumentHandlerRequestSize(reqSize,
							promhttp.InstrumentHandlerResponseSize(respSize,
								promhttp.InstrumentHandlerTimeToWriteHeader(timeToFirstHeader,
									http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
										next.ServeHTTP(w, r)

										if req.URL != nil {
											path := requestPath(req)
											reqDuration.ObserverVec = requestDuration.MustCurryWith(prometheus.Labels{LabelHandler: path})
											reqSize.ObserverVec = requestSize.MustCurryWith(prometheus.Labels{LabelHandler: path})
											respSize.ObserverVec = responseSize.MustCurryWith(prometheus.Labels{LabelHandler: path})
										}
									})),
							),
						),
					),
				),
			)

			chain.ServeHTTP(w, req)
		})
	}
}

// delayingObserver is used to set a label after the upstream handler has run, but the InstrumentationHandler has not yet called Observe.
// This solves this issue https://github.com/go-chi/chi/issues/200 without using an external library like httpsnoop.
// TL;DR of that issue: chi.RouteContext(r.Context()).RoutePattern() returns the RoutePattern only after the request has
// traversed all routers and middlewares have run.
type delayingObserver struct {
	prometheus.ObserverVec
}

// requestPath returns the matched request path and removes the trailing '/'.
// If no match is found it returns an empty string.
func requestPath(r *http.Request) string {
	path := chi.RouteContext(r.Context()).RoutePattern()
	return strings.TrimRight(path, "/")
}
