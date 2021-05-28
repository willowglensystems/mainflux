// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package factory_test

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	dockertest "github.com/ory/dockertest/v3"
	"github.com/mainflux/mainflux/pkg/messaging/factory"
	"github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Create nats container
	natsContainer, err := pool.Run("nats", "1.3.0", []string{})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	handleInterrupt(pool, natsContainer)

	// Create rabbitmq container
	rabbitContainer, err := pool.BuildAndRun("rabbitmq-amqp1.0", "./Dockerfile", []string{})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	handleInterrupt(pool, rabbitContainer)

	// Set environment variables for the factory to use
	os.Setenv("MF_NATS_URL", fmt.Sprintf("%s:%s", "localhost", natsContainer.GetPort("4222/tcp")))
	os.Setenv("MF_RABBITMQ_URL", fmt.Sprintf("%s:%s", "amqp://guest:guest@localhost", rabbitContainer.GetPort("5672/tcp")))
	os.Setenv("MF_RABBITMQ_TLS_CERTIFICATE", "rootCA.crt")
	os.Setenv("MF_RABBITMQ_TLS_KEY", "client.key")
	os.Setenv("MF_RABBITMQ_TLS_CA", "client.crt")
	os.Setenv("MF_NATS_TLS_CERTIFICATE", "rootCA.crt")
	os.Setenv("MF_NATS_TLS_KEY", "client.key")
	os.Setenv("MF_NATS_TLS_CA", "client.crt")
	
	// Wait for nats server to start
	if err := pool.Retry(func() error {
		os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.NatsMessagingSystem)
		_, err = factory.NewPublisher()
		return err
	}); err != nil {
		log.Fatalf("Nats container not started: %s", err)
	}

	// Wait for rabbitmq server to start
	if err := pool.Retry(func() error {
		os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.RabbitmqMessagingSystem)
		_, err = factory.NewPublisher()
		return err
	}); err != nil {
		log.Fatalf("Rabbitmq container not started: %s", err)
	}

	// Run tests
	code := m.Run()

	// Clean up nats container
	if err := pool.Purge(natsContainer); err != nil {
		log.Fatalf("Could not purge nats container: %s", err)
	}

	// Clean up rabbitmq container
	if err := pool.Purge(rabbitContainer); err != nil {
		log.Fatalf("Could not purge rabbit container: %s", err)
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
