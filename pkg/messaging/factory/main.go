package factory

import (
	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
	"github.com/mainflux/mainflux/pkg/messaging/rabbitmq"
)

const (
	rabbitmqSystem = "rabbitmq"
)

func NewPublisher() (messaging.Publisher, error) {
	systemType := queueConfiguration.GetSystem()
	configs, _, _ := queueConfiguration.GetConfig()

	if systemType == rabbitmqSystem {
		return rabbitmq.NewPublisher(configs[queueConfiguration.EnvRabbitmqURL])
	} else {
		return nats.NewPublisher(configs[queueConfiguration.EnvNatsURL])
	}
}

func NewPubSub(queue string, logger log.Logger) (messaging.PubSub, error) {
	systemType := queueConfiguration.GetSystem()
	configs, _, _ := queueConfiguration.GetConfig()

	if systemType == rabbitmqSystem {
		return rabbitmq.NewPubSub(configs[queueConfiguration.EnvRabbitmqURL], queue, logger)
	} else {
		return nats.NewPubSub(configs[queueConfiguration.EnvNatsURL], queue, logger)
	}
}

func GetAllChannels() (string) {
	systemType := queueConfiguration.GetSystem()

	if systemType == rabbitmqSystem {
		return rabbitmq.SubjectAllChannels;
	} else {
		return nats.SubjectAllChannels;
	}
}
