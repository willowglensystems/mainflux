// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
	"github.com/mainflux/mainflux/pkg/uuid"
	"github.com/Azure/go-amqp"
)

const (
	maxMessages            = 10
	pubTimeout             = 5
	receiveTimeout         = 1
)

const SubjectAllChannels = "channels.*"

var _ messaging.PubSub = (*pubsub)(nil)

type PubSub interface {
	messaging.PubSub
}

type pubsub struct {
	conn          *amqp.Client
	session       *amqp.Session
	logger        log.Logger
	queue         string
}

func NewPubSub(url, queue string, logger log.Logger) (PubSub, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	ret := &pubsub{
		conn:          conn,
		session:       session,
		queue:         queue,
		logger:        logger,
	}
	return ret, nil
}

func (pubsub *pubsub) Publish(topic string, msg messaging.Message) error {
	ctx := context.Background()

	sender, err := pubsub.session.NewSender(amqp.LinkTargetAddress(topic))
	if err != nil {
		pubsub.logger.Error( fmt.Sprintf( "Creating sender link: %s ", err) )
	}

	ctx, cancel := context.WithTimeout(ctx, pubTimeout * time.Second)

	message, err := pubsub.createMessage(topic, msg)

	if err != nil {
		pubsub.logger.Error( fmt.Sprintf( "Error creating message: %s", err ) )
	}

	// Send message
	err = sender.Send(ctx, message)
	if err != nil {
		pubsub.logger.Error( fmt.Sprintf( "Sending message: %s", err ) )
	}

	cancel()
	sender.Close(ctx)

	return nil
}

func (pubsub *pubsub) Subscribe(topic string, handler messaging.MessageHandler) error {
	receiver, err := pubsub.session.NewReceiver(
		amqp.LinkSourceAddress(topic),
		amqp.LinkCredit(maxMessages),
	)
	if err != nil {
		pubsub.logger.Error(fmt.Sprintf("Creating receiver link: %s", err))
	}

	messages := make(chan *amqp.Message)
	go handleMessages(messages, handler)
	go receiveMessages(messages, receiver)

	return nil
}

func handleMessages(messages chan *amqp.Message, handler messaging.MessageHandler) {
	for message := range messages {
		var msg = messaging.Message {
			Payload: message.GetData(),
		}
		handler( msg )
	}
}

func receiveMessages(messages chan *amqp.Message, receiver *amqp.Receiver) {
	ctx := context.Background()

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, receiveTimeout * time.Second)
		receiver.Close(ctx)
		cancel()
	}()


	for {
		msg, err := receiver.Receive(ctx)
		if err != nil {
			fmt.Println(fmt.Sprintf("Reading message from AMQP: %s", err))
			return
		}

		// Accept message
		msg.Accept(context.Background())

		messages <- msg
	}
}

func (pubsub *pubsub) Unsubscribe(topic string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, receiveTimeout * time.Second)
	if err := pubsub.session.Close(ctx); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}
	cancel()
	return nil
}

func (pubsub *pubsub) Close() {
	pubsub.conn.Close()
}

func (pubsub *pubsub) createMessage(topic string, msg messaging.Message) (*amqp.Message, error) {
	configs, queues, _ := queueConfiguration.GetConfig()

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

	correlationID, err := uuid.New().ID()

	if err != nil {
		return nil, errors.New("Unable to parse generate uuid")
	}

	message := amqp.NewMessage([]byte(msg.Payload))
	message.Header = &amqp.MessageHeader {
		Durable: durableValue,
		Priority: uint8(priorityValue),
		TTL: time.Duration( ttlValue ) * time.Millisecond,
	}
	message.Properties = &amqp.MessageProperties {
		ReplyTo: queues[topic],
		CorrelationID: correlationID,
		ContentType: configs[queueConfiguration.EnvRabbitmqContentType],
	}

	return message, nil
}
