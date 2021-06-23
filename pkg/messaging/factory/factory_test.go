package factory_test

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	mfLogger "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging/factory"
	queueConfiguration "github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
	dockertest "github.com/ory/dockertest/v3"
)

func TestUnsecuredNats(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	natsContainer, err := pool.Run("nats", "1.3.0", []string{})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	handleInterrupt(pool, natsContainer)

	// Set environment variables for the factory to use
	os.Setenv("MF_ENABLE_TLS", "false")
	os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.NatsMessagingSystem)
	os.Setenv("MF_NATS_URL", fmt.Sprintf("%s:%s", "localhost", natsContainer.GetPort("4222/tcp")))

	// Wait for nats server to start
	if err := pool.Retry(func() error {
		_, err = factory.NewPublisher()
		return err
	}); err != nil {
		log.Fatalf("Nats container not started: %s", err)
	}

	t.Run("T=NewPublisher", createPublisher)

	t.Run("T=NewPubSub", createPubSub)

	// Clean up nats container
	if err := pool.Purge(natsContainer); err != nil {
		log.Fatalf("Could not purge nats container: %s", err)
	}
}

func TestSecuredNats(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory: %s", err)
	}

	// Create nats container
	natsContainer, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "nats",
			Tag:        "1.3.0",
			Mounts: []string{
				workingDirectory + "/test_files:/config",
			},
			Cmd: []string{"-c", "/config/nats.conf"},
		},
	)

	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	handleInterrupt(pool, natsContainer)

	// Set environment variables for the factory to use
	os.Setenv("MF_ENABLE_TLS", "true")
	os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.NatsMessagingSystem)
	os.Setenv("MF_NATS_URL", fmt.Sprintf("%s:%s", "localhost", natsContainer.GetPort("4222/tcp")))
	os.Setenv("MF_NATS_TLS_CERTIFICATE", "test_files/client.crt")
	os.Setenv("MF_NATS_TLS_KEY", "test_files/client.key")
	os.Setenv("MF_NATS_TLS_CA", "test_files/rootCA.crt")

	// Wait for nats server to start
	if err := pool.Retry(func() error {
		_, err = factory.NewPublisher()
		return err
	}); err != nil {
		log.Fatalf("Nats container not started: %s", err)
	}

	t.Run("T=NewPublisher", createPublisher)

	t.Run("T=NewPubSub", createPubSub)

	// Clean up nats container
	if err := pool.Purge(natsContainer); err != nil {
		log.Fatalf("Could not purge nats container: %s", err)
	}
}

func TestUnsecuredRabbitmq(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Create rabbitmq container
	rabbitContainer, err := pool.BuildAndRun("rabbitmq-amqp1.0", "./Dockerfile", []string{})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	handleInterrupt(pool, rabbitContainer)

	// Set environment variables for the factory to use
	os.Setenv("MF_ENABLE_TLS", "false")
	os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.RabbitmqMessagingSystem)
	os.Setenv("MF_RABBITMQ_URL", fmt.Sprintf("%s:%s", "amqp://guest:guest@localhost", rabbitContainer.GetPort("5672/tcp")))

	// Wait for rabbitmq server to start
	if err := pool.Retry(func() error {
		_, err = factory.NewPublisher()
		return err
	}); err != nil {
		log.Fatalf("Rabbitmq container not started: %s", err)
	}

	t.Run("T=NewPublisher", createPublisher)

	t.Run("T=NewPubSub", createPubSub)

	// Clean up rabbitmq container
	if err := pool.Purge(rabbitContainer); err != nil {
		log.Fatalf("Could not purge rabbit container: %s", err)
	}
}

func TestSecuredRabbitmq(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	workingDirectory, err := os.Getwd()

	// Create rabbitmq container
	rabbitContainer, err := pool.BuildAndRunWithOptions(
		"./Dockerfile",
		&dockertest.RunOptions{
			Name: "rabbitmq-amqp1-0",
			Mounts: []string{
				workingDirectory + "/test_files:/etc/ssl",
			},
			Env: []string{
				"RABBITMQ_SSL_CACERTFILE=/etc/ssl/ca.crt",
				"RABBITMQ_SSL_CERTFILE=/etc/ssl/server.crt",
				"RABBITMQ_SSL_KEYFILE=/etc/ssl/server.key",
			},
		},
	)
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}
	handleInterrupt(pool, rabbitContainer)

	// Set environment variables for the factory to use
	os.Setenv("MF_ENABLE_TLS", "true")
	os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.RabbitmqMessagingSystem)
	os.Setenv("MF_RABBITMQ_URL", fmt.Sprintf("%s:%s", "amqps://guest:guest@localhost", rabbitContainer.GetPort("5671/tcp")))
	os.Setenv("MF_RABBITMQ_TLS_CERTIFICATE", "test_files/client.crt")
	os.Setenv("MF_RABBITMQ_TLS_KEY", "test_files/client.key")
	os.Setenv("MF_RABBITMQ_TLS_CA", "test_files/rootCA.crt")

	// Wait for rabbitmq server to start
	if err := pool.Retry(func() error {
		_, err = factory.NewPublisher()
		return err
	}); err != nil {
		log.Fatalf("Rabbitmq container not started: %s", err)
	}

	t.Run("T=NewPublisher", createPublisher)

	t.Run("T=NewPubSub", createPubSub)

	// Clean up rabbitmq container
	if err := pool.Purge(rabbitContainer); err != nil {
		log.Fatalf("Could not purge rabbit container: %s", err)
	}
}

func createPublisher(t *testing.T) {
	publisher, err := factory.NewPublisher()
	if publisher == nil || err != nil {
		t.Errorf(err.Error())
	}
}

func createPubSub(t *testing.T) {
	logger, err := mfLogger.New(os.Stdout, "debug")

	if err != nil {
		t.Errorf(err.Error())
	}

	pubsub, err := factory.NewPubSub("", logger)

	if pubsub == nil || err != nil {
		t.Errorf(err.Error())
	}
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
