// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0
package http

import (
	"net/http"

	"git.willowglen.ca/sq/third-party/mainflux.git"
	"git.willowglen.ca/sq/third-party/mainflux.git/auth"
	"git.willowglen.ca/sq/third-party/mainflux.git/auth/api/http/groups"
	"git.willowglen.ca/sq/third-party/mainflux.git/auth/api/http/keys"
	"github.com/go-zoo/bone"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MakeHandler(svc auth.Service, tracer opentracing.Tracer) http.Handler {
	mux := bone.New()
	mux = keys.MakeHandler(svc, mux, tracer)
	mux = groups.MakeHandler(svc, mux, tracer)
	mux.GetFunc("/version", mainflux.Version("auth"))
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
