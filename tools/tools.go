// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

//go:build tools

// Package tools is used to track build tool dependencies in go.mod
package tools

import (
	_ "gotest.tools/gotestsum"
	_ "mvdan.cc/gofumpt"
)
