// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"git.willowglen.ca/sq/third-party/mainflux.git/pkg/messaging"
)

type publishReq struct {
	msg   messaging.Message
	token string
}
