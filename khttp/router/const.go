// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package router

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kopexa-grc/common/khttp/cors"
	"github.com/kopexa-grc/common/khttp/metric"
	"github.com/kopexa-grc/common/khttp/security"
)

var defaultMiddlewareRegistry MiddlewareRegistry

var (
	defaultRecovererMiddleware      = middleware.Recoverer
	defaultCorsMiddleware           = cors.Middleware
	defaultSecurityHeaderMiddleware = security.Headers
	defaultRequestIDMiddleware      = middleware.RequestID
	defaultMetricMiddleware         = metric.Middleware
)
