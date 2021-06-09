// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"git.willowglen.ca/sq/third-party/mainflux.git/pkg/errors"
	"git.willowglen.ca/sq/third-party/mainflux.git/pkg/messaging"
)

var _ messaging.Publisher = (*mockBroker)(nil)

type mockBroker struct {
	subscriptions map[string]string
}

// NewBroker returns mock message publisher.
func NewBroker(sub map[string]string) messaging.Publisher {
	return &mockBroker{
		subscriptions: sub,
	}
}

func (mb mockBroker) Publish(topic string, msg messaging.Message) error {
	if len(msg.Payload) == 0 {
		return errors.New("failed to publish")
	}
	return nil
}

func (mb mockBroker) Close() {

}
