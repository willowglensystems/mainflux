// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mainflux/mainflux/queue-configuration"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/Azure/go-amqp"
	"github.com/google/uuid"
)

const (
	publishTimeout         = 5
)

var _ messaging.Publisher = (*publisher)(nil)

type publisher struct {
	conn *amqp.Client
	session *amqp.Session
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

	ret := &publisher{
		conn: conn,
		session: session,
	}

	return ret, nil
}

func (pub *publisher) Publish(topic string, msg messaging.Message) error {
	ctx := context.Background()

	sender, err := pub.session.NewSender(amqp.LinkTargetAddress(topic))
	if err != nil {
		fmt.Println("Creating sender link:", err)
	}

	ctx, cancel := context.WithTimeout(ctx, publishTimeout * time.Second)

	// Send message
	err = sender.Send(ctx, pub.createMessage(topic,msg))
	if err != nil {
		fmt.Println("Sending message:", err)
	}

	cancel()
	sender.Close(ctx)

	return nil
}

func (pub *publisher) Close() {
	pub.conn.Close()
}

func (pub *publisher) createMessage(topic string, msg messaging.Message) *amqp.Message {
	configs, queues, _ := queueConfiguration.GetConfig()

	durableValue, err := strconv.ParseBool(configs[queueConfiguration.EnvRabbitmqDurable])

	if err != nil {
		fmt.Println("Unable to parse Durable configuration, defaulting to false")
		durableValue = false
	}

	priorityValue, err := strconv.ParseUint(configs[queueConfiguration.EnvRabbitmqPriority], 10, 64)

	if err != nil {
		fmt.Println("Unable to parse Priority configuration, defaulting to 1")
		priorityValue = 1
	}


	ttlValue, err := strconv.ParseUint(configs[queueConfiguration.EnvRabbitmqTTL], 10, 64)

	if err != nil {
		fmt.Println("Unable to parse TTL configuration, defaulting to 3600000 milliseconds")
		ttlValue = 3600000
	}

	message := amqp.NewMessage([]byte(msg.Payload))
	message.Header = &amqp.MessageHeader {
		Durable: durableValue,
		Priority: uint8(priorityValue),
		TTL: time.Duration( ttlValue ) * time.Millisecond,
	}
	message.Properties = &amqp.MessageProperties {
		ReplyTo: queues[topic],
		CorrelationID: uuid.New(),
		ContentType: configs[queueConfiguration.EnvRabbitmqContentType],
	}

	return message
}
