// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/Azure/go-amqp"
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

	if err != nil {
		return err
	}

	ctx := context.Background()

	sender, err := pub.session.NewSender(amqp.LinkTargetAddress(topic))
	if err != nil {
		fmt.Println("Creating sender link:", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)

	// Send message
	err = sender.Send(ctx, amqp.NewMessage([]byte(data)))
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
