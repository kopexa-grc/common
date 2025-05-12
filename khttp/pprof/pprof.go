// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package pprof

import (
	"net/http"
	"net/http/pprof"
)

const PathPrefixPProf = "/debug/pprof"

// Handler returns a preconfigured pprof handler equal to the handler registered to the default HTTP server by importing the pprof package itself
// This should be mounted at the location defined by the paashttp.PPathPrefixPProf constant.
// Please mount this
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(PathPrefixPProf+"/cmdline", pprof.Cmdline)
	mux.HandleFunc(PathPrefixPProf+"/profile", pprof.Profile)
	mux.HandleFunc(PathPrefixPProf+"/symbol", pprof.Symbol)
	mux.HandleFunc(PathPrefixPProf+"/trace", pprof.Trace)
	mux.HandleFunc(PathPrefixPProf+"/", pprof.Index)

	return mux
}
