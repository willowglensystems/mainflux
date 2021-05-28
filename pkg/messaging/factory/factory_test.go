package factory_test

import (
	"os"
	"testing"

	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging/factory"
	"github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
)

func TestNatsPublisher(t *testing.T) {
	os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.NatsMessagingSystem)

	publisher, err := factory.NewPublisher()
	if publisher == nil || err != nil {
		t.Errorf(err.Error())
	}
}

func TestNatsPubsub(t *testing.T) {
	os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.NatsMessagingSystem)

	logger, err := log.New(os.Stdout, "debug")

	if err != nil {
		t.Errorf(err.Error())
	}

	pubsub, err := factory.NewPubSub("", logger)

	if pubsub == nil || err != nil {
		t.Errorf(err.Error())
	}
}

func TestRabbitPublisher(t *testing.T) {
	os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.RabbitmqMessagingSystem)
	
	publisher, err := factory.NewPublisher()
	if publisher == nil || err != nil {
		t.Errorf(err.Error())
	}
}

func TestRabbitPubsub(t *testing.T) {
	os.Setenv("MF_QUEUE_SYSTEM", queueConfiguration.RabbitmqMessagingSystem)

	logger, err := log.New(os.Stdout, "debug")

	if err != nil {
		t.Errorf(err.Error())
	}

	pubsub, err := factory.NewPubSub("", logger)

	if pubsub == nil || err != nil {
		t.Errorf(err.Error())
	}
}
