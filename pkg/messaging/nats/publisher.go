// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package nats

import (
	"crypto/tls"
	"fmt"

	"git.willowglen.ca/sq/third-party/mainflux/pkg/messaging"
	"github.com/gogo/protobuf/proto"
	broker "github.com/nats-io/nats.go"
)

var _ messaging.Publisher = (*publisher)(nil)

type PublisherOption func(*publisher)

type publisher struct {
	conn      *broker.Conn
	tlsConfig *tls.Config
}

// Publisher wraps messaging Publisher exposing
// Close() method for NATS connection.
type Publisher interface {
	messaging.Publisher
	Close()
}

// NewPublisher returns NATS message Publisher.
func NewPublisher(url string, tlsConfig *tls.Config) (Publisher, error) {
	var conn *broker.Conn
	var err error

	if tlsConfig != nil {
		cfg := broker.Secure(tlsConfig)
		conn, err = broker.Connect(url, cfg)
	} else {
		conn, err = broker.Connect(url)
	}

	if err != nil {
		return nil, err
	}

	ret := &publisher{
		conn: conn,
	}

	return ret, nil
}

func (pub *publisher) Publish(topic string, msg messaging.Message) error {
	data, err := proto.Marshal(&msg)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("%s.%s", chansPrefix, topic)
	if msg.Subtopic != "" {
		subject = fmt.Sprintf("%s.%s", subject, msg.Subtopic)
	}
	if err := pub.conn.Publish(subject, data); err != nil {
		return err
	}

	return nil
}

func (pub *publisher) Close() {
	pub.conn.Close()
}
