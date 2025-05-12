// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob

//go:generate  go run -mod=mod go.uber.org/mock/mockgen -destination=./driver_mock_test.go -package=blob_test -source=./driver/driver.go Bucket
