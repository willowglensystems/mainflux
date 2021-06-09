// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"git.willowglen.ca/sq/third-party/mainflux.git/opcua"
	"github.com/go-kit/kit/endpoint"
)

func browseEndpoint(svc opcua.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(browseReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		nodes, err := svc.Browse(req.ServerURI, req.Namespace, req.Identifier)
		if err != nil {
			return nil, err
		}

		res := browseRes{
			Nodes: nodes,
		}

		return res, nil
	}
}
