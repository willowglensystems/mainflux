// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"context"
	"fmt"
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

type publisher struct {
	conn    *amqp.Client
	session *amqp.Session
	configs *queueConfiguration.Config
}

type Publisher interface {
	messaging.Publisher
}


func NewPublisher(url string) (Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	configs, _, _ := queueConfiguration.GetConfig()

	ret := &publisher{
		conn: conn,
		session: session,
		configs: configs,
	}

	return ret, nil
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
		Priority: uint8(configs.RabbitmqPriority),
		TTL: time.Duration(configs.RabbitmqTTL) * time.Millisecond,
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
