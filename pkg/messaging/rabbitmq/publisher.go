// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"context"
	"fmt"
	"strconv"
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
	data, err := proto.Marshal(&msg)

	ctx := context.Background()

	sender, err := pub.session.NewSender(amqp.LinkTargetAddress(topic))
	if err != nil {
		fmt.Println("Creating sender link:", err)
	}

	ctx, cancel := context.WithTimeout(ctx, publishTimeout * time.Second)

	message, err := createMessage(topic, data)

	if err != nil {
		fmt.Println( fmt.Sprintf( "Error creating message: %s", err ) )
	}

	// Send message
	err = sender.Send(ctx, message)
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

func createMessage(topic string, data []byte) (*amqp.Message, error) {
	configs, _, _ := queueConfiguration.GetConfig()

	durableValue, err := strconv.ParseBool(configs[queueConfiguration.EnvRabbitmqDurable])

	if err != nil {
		fmt.Println("Unable to parse Durable configuration, defaulting to false")
		durableValue, _ = strconv.ParseBool(queueConfiguration.DefRabbitmqDurable)
	}

	priorityValue, err := strconv.ParseUint(configs[queueConfiguration.EnvRabbitmqPriority], 10, 64)

	if err != nil {
		fmt.Println("Unable to parse Priority configuration, defaulting to 1")
		priorityValue, _ = strconv.ParseUint(queueConfiguration.DefRabbitmqPriority, 10, 64)
	}

	ttlValue, err := strconv.ParseUint(configs[queueConfiguration.EnvRabbitmqTTL], 10, 64)

	if err != nil {
		fmt.Println("Unable to parse TTL configuration, defaulting to 3600000 milliseconds")
		ttlValue, _ = strconv.ParseUint(queueConfiguration.DefRabbitmqTTL, 10, 64)
	}

	message := amqp.NewMessage(data)
	message.Header = &amqp.MessageHeader {
		Durable: durableValue,
		Priority: uint8(priorityValue),
		TTL: time.Duration( ttlValue ) * time.Millisecond,
	}
	message.Properties = &amqp.MessageProperties {
		//CorrelationID: correlationID, TODO add when metadata changes are in
		ContentType: configs[queueConfiguration.EnvRabbitmqContentType],
	}

	return message, nil
}
