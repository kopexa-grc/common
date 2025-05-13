// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fgamock

//go:generate go run -mod=mod go.uber.org/mock/mockgen -destination=./fga.go -package=fgamock github.com/openfga/go-sdk/client SdkClient,SdkClientCheckRequestInterface,SdkClientWriteRequestInterface,SdkClientReadRequestInterface
