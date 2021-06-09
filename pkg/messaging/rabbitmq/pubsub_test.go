// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq_test

import (
	"fmt"
	"testing"

	"git.willowglen.ca/sq/third-party/mainflux.git/pkg/messaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	topic       = "topic"
	chansPrefix = "channels"
	channel     = "9b7b1b3f-b1b0-46a8-a717-b8213f9eda3b"
	subtopic    = "queue"
)

var (
	msgChan = make(chan messaging.Message)
	data    = []byte("payload")
)

func TestPubsub(t *testing.T) {
	err := pubsub.Subscribe(topic, handler)

	cases := []struct {
		desc     string
		channel  string
		subtopic string
		payload  []byte
	}{
		{
			desc:     "publish message with string payload",
			payload:  data,
			subtopic: subtopic,
		},
	}

	for _, tc := range cases {
		expectedMsg := messaging.Message{
			Payload: tc.payload,
		}
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		err = pubsub.Publish(topic, expectedMsg)

		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		receivedMsg := <-msgChan
		assert.Equal(t, expectedMsg, receivedMsg, fmt.Sprintf("%s: expected %+v got %+v\n", tc.desc, expectedMsg, receivedMsg))
	}
}

func handler(msg messaging.Message) error {
	msgChan <- msg
	return nil
}
