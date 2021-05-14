// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/streadway/amqp"
)

var _ messaging.Publisher = (*publisher)(nil)

type publisher struct {
	conn *amqp.Connection
	channel *amqp.Channel
}

type Publisher interface {
	messaging.Publisher
}


func NewPublisher(url string) (Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	ret := &publisher{
		conn: conn,
		channel: channel,
	}

	return ret, nil
}

func (pub *publisher) Publish(topic string, msg messaging.Message) error {
	data, err := proto.Marshal(&msg)

	if err != nil {
		return err
	}

	exchangeName := topic
	exchangeType := "topic"
	queueName := msg.Subtopic
	routingKey := msg.Subtopic

	fmt.Println( exchangeName )
	fmt.Println( exchangeType )
	fmt.Println( queueName )
	fmt.Println( routingKey )

	if err = pub.channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	_, err = pub.channel.QueueDeclare( 
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	);
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	if err := pub.channel.QueueBind( 
		queueName,    // name
		routingKey,   // key
		exchangeName, // exchange
		false,        // noWait
		nil,          // args
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	if err := pub.channel.Publish(
		topic,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(data),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

func (pub *publisher) Close() {
	pub.conn.Close()
}
