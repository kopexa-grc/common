// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package router

func DefaultRecovererMiddleware() Middleware {
	return defaultRecovererMiddleware
}

func DefaultCorsMiddleware() Middleware {
	return defaultCorsMiddleware
}

func DefaultSecurityHeaderMiddleware() Middleware {
	return defaultSecurityHeaderMiddleware
}

func DefaultRequestIDMiddleware() Middleware {
	return defaultRequestIDMiddleware
}

func DefaultMetricMiddleware() Middleware {
	return defaultMetricMiddleware
}
