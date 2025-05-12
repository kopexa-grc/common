// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package metric

// GlobalRegistry used to register prometheus metrics
var GlobalRegistry = NewRegistry()

// Middleware to use to collect the metrics
var Middleware = GlobalRegistry.PrometheusMiddleware()

// Handler to use to expose the metrics
var Handler = GlobalRegistry.Handler()
