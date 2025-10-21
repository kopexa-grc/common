// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fgamock

//go:generate go run -mod=mod go.uber.org/mock/mockgen -destination=./fga.go -package=fgamock github.com/openfga/go-sdk/client SdkClient,SdkClientCheckRequestInterface,SdkClientWriteRequestInterface,SdkClientReadRequestInterface,SdkClientListObjectsRequestInterface,SdkClientListStoresRequestInterface,SdkClientCreateStoreRequestInterface,SdkClientReadAuthorizationModelsRequestInterface,SdkClientWriteAuthorizationModelRequestInterface
