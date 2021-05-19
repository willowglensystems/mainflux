// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package rabbitmq_test

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/rabbitmq"
	dockertest "github.com/ory/dockertest/v3"
)

var (
	publisher messaging.Publisher
	pubsub    messaging.PubSub
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	container, err := pool.BuildAndRun("rabbitmq-amqp1.0", "./Dockerfile", []string{})

	handleInterrupt(pool, container)

	if err != nil {
		log.Fatalf("Container creation failed: %s", err)
	}

	address := fmt.Sprintf("%s:%s", "amqp://guest:guest@localhost", container.GetPort("5672/tcp"))

	if err := pool.Retry(func() error {
		publisher, err = rabbitmq.NewPublisher(address)
		return err
	}); err != nil {
		log.Fatalf("Failed to create publisher: %s", err)
	}

	logger, err := logger.New(os.Stdout, "error")
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err := pool.Retry(func() error {
		pubsub, err = rabbitmq.NewPubSub(address, "queue", logger)
		return err
	}); err != nil {
		log.Fatalf("Failed to create pubsub: %s", err)
	}

	code := m.Run()
	if err := pool.Purge(container); err != nil {
		log.Fatalf("Could not purge container: %s", err)
	}

	os.Exit(code)
}

func handleInterrupt(pool *dockertest.Pool, container *dockertest.Resource) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if err := pool.Purge(container); err != nil {
			log.Fatalf("Could not purge container: %s", err)
		}
		os.Exit(0)
	}()
}
