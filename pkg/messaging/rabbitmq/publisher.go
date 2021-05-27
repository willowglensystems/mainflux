// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
	"github.com/Azure/go-amqp"
)

const (
	publishTimeout         = 5
)

var _ messaging.Publisher = (*publisher)(nil)

type PublisherOption func(*publisher)

type publisher struct {
	conn      *amqp.Client
	session   *amqp.Session
	configs   *queueConfiguration.Config
	tlsConfig *tls.Config
}

type Publisher interface {
	messaging.Publisher
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

func NewPublisher(url string, opts ...PublisherOption) (Publisher, error) {
	newPub := &publisher{
		tlsConfig: nil,
	}
	for _, opt := range opts {
		opt(newPub)
	}

	var conn *amqp.Client
	var err error

	if newPub.tlsConfig != nil {
		cfg := amqp.ConnTLSConfig(newPub.tlsConfig)
		conn, err = amqp.Dial(url, cfg)
	} else {
		conn, err = amqp.Dial(url)
	}
	
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	configs, _, _ := queueConfiguration.GetConfig()

	newPub.conn = conn
	newPub.session = session
	newPub.configs = configs

	return newPub, nil
}

func (pub *publisher) Publish(topic string, msg messaging.Message) error {
	ctx := context.Background()

	sender, err := pub.session.NewSender(amqp.LinkTargetAddress(topic))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, publishTimeout * time.Second)

	message, err := createMessage(topic, &msg, pub.configs)

	if err != nil {
		return err
	}

	// Send message
	err = sender.Send(ctx, message)
	if err != nil {
		return err
	}

	cancel()
	sender.Close(ctx)

	return nil
}

func (pub *publisher) Close() {
	pub.conn.Close()
}

func createMessage(topic string, msg *messaging.Message, configs *queueConfiguration.Config) (*amqp.Message, error) {
	data, err := proto.Marshal(msg)

	if err != nil {
		return nil, err
	}

	message := amqp.NewMessage(data)
	message.Header = &amqp.MessageHeader {
		Durable: configs.RabbitmqDurable,
		Priority: configs.RabbitmqPriority,
		TTL: configs.RabbitmqTTL,
	}
	message.Properties = &amqp.MessageProperties {
		CorrelationID: string(msg.Metadata["CorrelationID"]),
		ContentType: configs.RabbitmqContentType,
		CreationTime: time.Unix(msg.Created, 0),
		ReplyTo: string(msg.Metadata["ReplyTo"]),
		Subject: string(msg.Metadata["Type"]),
	}

	return message, nil
}
