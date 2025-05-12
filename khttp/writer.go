// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package khttp

import (
	"encoding/json"
	"net/http"

	kerr "github.com/kopexa-grc/common/errors"
)

// WriteJSON encodes the given object to JSON and writes it to the
// http.ResponseWriter with the given statusCode.
func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Add("Content-Type", ContentTypeJSON)
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

// WriteErr writes the paas error to the response using json
func WriteErr(w http.ResponseWriter, err *kerr.Error) {
	_ = WriteJSON(w, err.Status, err)
}

// WriteSimpleErr writes the kerr error to the response using json based on the err.
// This will always be understood as unhandled error (HTTP 500). If the error was
// detected to ba a kerr.KError it will automatically use WriteErr and behave the same.
func WriteSimpleErr(w http.ResponseWriter, err error) {
	pe, ok := (err).(*kerr.Error) // nolint: errorlint
	if ok {
		WriteErr(w, pe)
	} else {
		WriteErr(w, kerr.NewUnexpectedFailure("request failed"))
	}
}

// WriteNoContent writes a 204 No Content response to the http.ResponseWriter.
func WriteNoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
