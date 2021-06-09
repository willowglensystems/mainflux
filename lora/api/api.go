// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"

	"git.willowglen.ca/sq/third-party/mainflux"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler() http.Handler {
	r := bone.New()
	r.GetFunc("/version", mainflux.Version("lora-adapter"))
	r.Handle("/metrics", promhttp.Handler())

	return r
}
