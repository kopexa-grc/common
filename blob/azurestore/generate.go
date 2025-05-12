// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package azurestore

//go:generate go run -mod=mod go.uber.org/mock/mockgen -destination=./store_mock_test.go -package=azurestore_test -source=./service.go AzService,AzBlob
