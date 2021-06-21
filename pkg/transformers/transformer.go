// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package transformers

import "git.willowglen.ca/sq/third-party/mainflux.git/pkg/messaging"

// Transformer specifies API form Message transformer.
type Transformer interface {
	// Transform Mainflux message to any other format.
	Transform(msg messaging.Message) (interface{}, error)
}
