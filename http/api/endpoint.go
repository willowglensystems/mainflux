// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"git.willowglen.ca/sq/third-party/mainflux/http"
	"github.com/go-kit/kit/endpoint"
)

func sendMessageEndpoint(svc http.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(publishReq)
		err := svc.Publish(ctx, req.token, req.msg)
		return nil, err
	}
}
