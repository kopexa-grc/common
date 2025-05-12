// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package metric

// Metric Names
const (
	MetricInFlightRequests  = "in_flight_requests"
	MetricAPITotalRequests  = "api_requests_total"
	MetricRequestDuration   = "request_duration_seconds"
	MetricTimeToFirstHeader = "time_to_first_header_duration_seconds"
	MetricResponseSize      = "response_size_bytes"
	MetricRequestSize       = "request_size_bytes"
)

// Metric Labels
const (
	LabelCode    = "code"
	LabelMethod  = "method"
	LabelHandler = "handler"
)

// Metric Help Messages
const (
	HelpInFlightRequests  = "Number of HTTP requests in flight"
	HelpAPITotalRequests  = "Requests total per method and code"
	HelpRequestDuration   = "A histogram of latencies for requests."
	HelpTimeToFirstHeader = "A histogram of time to first header latencies"
	HelpResponseSize      = "A histogram of response sizes for requests."
	HelpRequestSize       = "A histogram of request sizes for requests."
)

// Default Values
const (
	DefaultHandler = "undefined"
)

// Histogram Buckets
var (
	DurationBuckets = []float64{0.001, 0.01, 0.015, .25, .5, 1, 1.5, 2.5, 5, 10}
	SizeBuckets     = []float64{200, 500, 900, 1500, 2000, 10000}
)
