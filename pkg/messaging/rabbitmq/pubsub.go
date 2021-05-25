// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	log "github.com/mainflux/mainflux/logger"
	"github.com/gogo/protobuf/proto"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
	"github.com/Azure/go-amqp"
)

const (
	maxMessages            = 10
	pubTimeout             = 5
	receiveTimeout         = 1
)

const SubjectAllChannels = "channels.*"

var (
	errAlreadySubscribed = errors.New("already subscribed to topic")
	errNotSubscribed     = errors.New("not subscribed")
	errEmptyTopic        = errors.New("empty topic")
	messages             = make(chan *amqp.Message)
)

var _ messaging.PubSub = (*pubsub)(nil)


type PubSub interface {
	messaging.PubSub
}

type pubsub struct {
	conn          *amqp.Client
	session       *amqp.Session
	mu            sync.Mutex
	logger        log.Logger
	queue         string
	subscriptions map[string]bool
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
		subscriptions: make(map[string]bool),
	}
	return ret, nil
}

func (pubsub *pubsub) Publish(topic string, msg messaging.Message) error {
	data, err := proto.Marshal(&msg)
	ctx := context.Background()

	sender, err := pubsub.session.NewSender(amqp.LinkTargetAddress(topic))
	if err != nil {
		pubsub.logger.Error( fmt.Sprintf( "Creating sender link: %s ", err) )
	}

	ctx, cancel := context.WithTimeout(ctx, pubTimeout * time.Second)

	message, err := pubsub.createMessage(topic, data)

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
	if topic == "" {
                return errEmptyTopic
        }
        pubsub.mu.Lock()
        defer pubsub.mu.Unlock()

	if pubsub.subscriptions[topic] {
                return errAlreadySubscribed
        }

	pubsub.subscriptions[topic] = true

	receiver, err := pubsub.session.NewReceiver(
		amqp.LinkSourceAddress(topic),
		amqp.LinkCredit(maxMessages),
	)
	if err != nil {
		pubsub.logger.Error(fmt.Sprintf("Creating receiver link: %s", err))
	}

	go pubsub.handleMessages(messages, handler)
	go pubsub.receiveMessages(topic, messages, receiver)

	return nil
}

func (pubsub *pubsub) handleMessages(messages chan *amqp.Message, handler messaging.MessageHandler) {
	for message := range messages {
		var msg messaging.Message
                if err := proto.Unmarshal(message.GetData(), &msg); err != nil {
                        pubsub.logger.Warn(fmt.Sprintf("Failed to unmarshal received message: %s", err))
                        return
                }
		handler( msg )
	}
}

func (pubsub *pubsub) receiveMessages(topic string, messages chan *amqp.Message, receiver *amqp.Receiver) {
	ctx := context.Background()

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, receiveTimeout * time.Second)
		receiver.Close(ctx)
		cancel()
	}()


	for pubsub.subscriptions[topic] {
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
	if topic == "" {
                return errEmptyTopic
        }
        pubsub.mu.Lock()
        defer pubsub.mu.Unlock()

        subscribed, ok := pubsub.subscriptions[topic]
        if !ok || !subscribed {
                return errNotSubscribed
        }

	pubsub.subscriptions[topic] = false

	return nil
}

func (pubsub *pubsub) Close() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, receiveTimeout * time.Second)
	if err := pubsub.session.Close(ctx); err != nil {
		fmt.Sprintf("Consumer cancel failed: %s", err)
	}
	cancel()

	pubsub.conn.Close()
}

func (pubsub *pubsub) createMessage(topic string, data []byte) (*amqp.Message, error) {
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

	message := amqp.NewMessage(data)
	message.Header = &amqp.MessageHeader {
		Durable: durableValue,
		Priority: uint8(priorityValue),
		TTL: time.Duration( ttlValue ) * time.Millisecond,
	}
	message.Properties = &amqp.MessageProperties {
		ReplyTo: queues[topic],
		//CorrelationID: correlationID, TODO add when metadata changes are in
		ContentType: configs[queueConfiguration.EnvRabbitmqContentType],
	}

	return message, nil
}
