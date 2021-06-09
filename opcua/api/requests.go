// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import "git.willowglen.ca/sq/third-party/mainflux.git/opcua"

type browseReq struct {
	ServerURI  string
	Namespace  string
	Identifier string
}

func (req *browseReq) validate() error {
	if req.ServerURI == "" {
		return opcua.ErrMalformedEntity
	}

	return nil
}
