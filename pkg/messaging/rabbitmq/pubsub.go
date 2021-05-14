// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/streadway/amqp"
)

const SubjectAllChannels = "channels.*"

var _ messaging.PubSub = (*pubsub)(nil)

type PubSub interface {
	messaging.PubSub
}

type pubsub struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	logger        log.Logger
	queue         string
}

func NewPubSub(url, queue string, logger log.Logger) (PubSub, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	ret := &pubsub{
		conn:          conn,
		channel:       channel,
		queue:         queue,
		logger:        logger,
	}
	return ret, nil
}

func (pubsub *pubsub) Publish(topic string, msg messaging.Message) error {
	data, err := proto.Marshal(&msg)

	if err != nil {
		return err
	}

	exchangeName := topic
	exchangeType := "topic"
	queueName := msg.Subtopic
	routingKey := msg.Subtopic

	if err := pubsub.channel.ExchangeDeclare(
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

	_, err = pubsub.channel.QueueDeclare( 
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

	if err := pubsub.channel.QueueBind( 
		queueName,    // name
		routingKey,   // key
		exchangeName, // exchange
		false,        // noWait
		nil,          // args
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	if err = pubsub.channel.Publish(
		exchangeName, // publish to an exchange
		routingKey,   // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
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

func (pubsub *pubsub) Subscribe(topic string, handler messaging.MessageHandler) error {

	exchangeName := topic
	exchangeType := "topic"
	queueName := pubsub.queue
	routingKey := pubsub.queue

	channel, err := pubsub.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	if err := channel.ExchangeDeclare(
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

	_, err = channel.QueueDeclare( 
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

	if err := channel.QueueBind( 
		queueName,    // name
		routingKey,   // key
		exchangeName, // exchange
		false,        // noWait
		nil,          // args
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	messages, err := channel.Consume(
		queueName,  // name
		"",         // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(messages, handler)

	return nil
}

func handle(messages <-chan amqp.Delivery, handler messaging.MessageHandler) {
	for message := range messages {
		//fmt.Printf( "got %dB delivery: [%v] %q", len(message.Body), message.DeliveryTag, message.Body )
		message.Ack(false)

		var msg messaging.Message
		if err := proto.Unmarshal(message.Body, &msg); err != nil {
			fmt.Println(fmt.Sprintf("Failed to unmarshal received message: %s", err))
			return
		}

		handler( msg )
	}
}

func (pubsub *pubsub) Unsubscribe(topic string) error {
	if err := pubsub.channel.Cancel("", true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}
	return nil
}

func (pubsub *pubsub) Close() {
	pubsub.conn.Close()
}

