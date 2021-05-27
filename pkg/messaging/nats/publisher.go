// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package nats

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
	"github.com/mainflux/mainflux/pkg/messaging"
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

// Constructor option to create a publisher with TLS access
func WithPublisherTLS(caFile, certFile, keyFile string) PublisherOption {
	return func(pub *publisher) {
		roots := x509.NewCertPool()

		data, err := ioutil.ReadFile(caFile)

		if err != nil {
			fmt.Errorf("Error reading in root CA: %s", err)
			return
		}

		roots.AppendCertsFromPEM(data)

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			fmt.Errorf("Error loading certificate: %s", err)
			return
		}

		pub.tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs: roots,
			InsecureSkipVerify: true,
		}
	}
}

// NewPublisher returns NATS message Publisher.
func NewPublisher(url string, opts ...PublisherOption) (Publisher, error) {
	newPub := &publisher{
		tlsConfig: nil,
	}
	for _, opt := range opts {
		opt(newPub)
	}

	var conn *broker.Conn
	var err error

	if newPub.tlsConfig != nil {
		cfg := broker.Secure(newPub.tlsConfig)
		conn, err = broker.Connect(url, cfg)
	} else {
		conn, err = broker.Connect(url)
	}

	if err != nil {
		return nil, err
	}

	newPub.conn = conn

	return newPub, nil
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
