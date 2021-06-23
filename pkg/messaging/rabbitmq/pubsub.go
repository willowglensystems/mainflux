// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	queueConfiguration "github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
	"github.com/Azure/go-amqp"
	"github.com/gogo/protobuf/proto"
)

const (
	maxMessages    = 10
	pubTimeout     = 5
	receiveTimeout = 1
)

const SubjectAllChannels = "channels.*"

var (
	errAlreadySubscribed = errors.New("already subscribed to topic")
	errNotSubscribed     = errors.New("not subscribed")
	errEmptyTopic        = errors.New("empty topic")
)

var _ messaging.PubSub = (*pubsub)(nil)

type PubSubOption func(*pubsub)

type PubSub interface {
	messaging.PubSub
}

type pubsub struct {
	conn          *amqp.Client
	session       *amqp.Session
	configs       *queueConfiguration.Config
	mu            sync.Mutex
	logger        log.Logger
	queue         string
	subscriptions map[string]bool
	tlsConfig     *tls.Config
}

func NewPubSub(url, queue string, logger log.Logger, tlsConfig *tls.Config) (PubSub, error) {
	var conn *amqp.Client
	var err error

	if tlsConfig != nil {
		cfg := amqp.ConnTLSConfig(tlsConfig)
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

	ret := &pubsub{
		conn:          conn,
		session:       session,
		configs:       configs,
		queue:         queue,
		logger:        logger,
		subscriptions: make(map[string]bool),
	}

	return ret, nil
}

func (pubsub *pubsub) Publish(topic string, msg messaging.Message) error {
	ctx := context.Background()

	subject := topic

	if msg.Subtopic != "" {
		subject = fmt.Sprintf("/exchange/%s/%s", subject, msg.Subtopic)
	}

	sender, err := pubsub.session.NewSender(amqp.LinkTargetAddress(subject))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, pubTimeout*time.Second)

	message, err := createMessage(topic, &msg, pubsub.configs)

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

// Subscribe to a topic (RabbitMQ queue) with a message handler.
// Note: the message handler is called from a goroutine so concurrency must be taken into account
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

	messages := make(chan *amqp.Message)

	go pubsub.handleMessages(messages, handler)
	go pubsub.receiveMessages(topic, messages, receiver)

	return nil
}

func (pubsub *pubsub) handleMessages(messages chan *amqp.Message, handler messaging.MessageHandler) {
	for message := range messages {
		var msg messaging.Message
		if err := proto.Unmarshal(message.GetData(), &msg); err != nil {
			pubsub.logger.Warn(fmt.Sprintf("Failed to unmarshal received message: %s", err))
		} else {
			handler(msg)
		}
	}
}

func (pubsub *pubsub) receiveMessages(topic string, messages chan *amqp.Message, receiver *amqp.Receiver) {
	ctx := context.Background()

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, receiveTimeout*time.Second)
		receiver.Close(ctx)
		cancel()
	}()

	for pubsub.subscriptions[topic] {
		msg, err := receiver.Receive(ctx)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error receiving message from AMQP: %s", err))
		} else {
			// Accept message
			msg.Accept(context.Background())

			messages <- msg
		}
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
	ctx, cancel := context.WithTimeout(ctx, receiveTimeout*time.Second)
	if err := pubsub.session.Close(ctx); err != nil {
		fmt.Sprintf("Consumer cancel failed: %s", err)
	}
	cancel()

	pubsub.conn.Close()
}
