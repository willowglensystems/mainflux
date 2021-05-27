package factory

import (
	"errors"
	"fmt"

	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
	"github.com/mainflux/mainflux/pkg/messaging/rabbitmq"
)

const (
	rabbitmqSystem = "rabbitmq"
)

// Generates a NATS or RabbitMQ publisher based on the messaging queue configuration
func NewPublisher() (messaging.Publisher, error) {
	systemType := queueConfiguration.GetSystem()
	configs, _, _ := queueConfiguration.GetConfig()

	if systemType == queueConfiguration.RabbitmqMessagingSystem {
		if configs.RabbitmqTLSCA != "" &&
		   configs.RabbitmqTLSCertificate != "" &&
		   configs.RabbitmqTLSKey != "" {
			return rabbitmq.NewPublisher(
				configs.RabbitmqURL, 
				rabbitmq.WithPublisherTLS(
					configs.RabbitmqTLSCA,
					configs.RabbitmqTLSCertificate,
					configs.RabbitmqTLSKey,
				),
			)
		} else {
			return rabbitmq.NewPublisher(configs.RabbitmqURL)
		}
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		if configs.NatsTLSCA != "" &&
		   configs.NatsTLSCertificate != "" &&
		   configs.NatsTLSKey != "" {
			return nats.NewPublisher(
				configs.NatsURL, 
				nats.WithPublisherTLS(
					configs.NatsTLSCA,
					configs.NatsTLSCertificate,
					configs.NatsTLSKey,
				),
			)
		} else {
			return nats.NewPublisher(configs.NatsURL)
		}
	} else {
		fmt.Println("Invalid messaging system type for creating a publisher:", systemType)
		return nil, errors.New("Invalid queue type")
	}
}

// Generates a NATS or RabbitMQ publisher/subscriber based on the messaging queue configuration
func NewPubSub(queue string, logger log.Logger) (messaging.PubSub, error) {
	systemType := queueConfiguration.GetSystem()
	configs, _, _ := queueConfiguration.GetConfig()

	if systemType == queueConfiguration.RabbitmqMessagingSystem {
		if configs.RabbitmqTLSCA != "" &&
		   configs.RabbitmqTLSCertificate != "" &&
		   configs.RabbitmqTLSKey != "" {
			return rabbitmq.NewPubSub(
				configs.RabbitmqURL, 
				queue,
				logger,
				rabbitmq.WithPubSubTLS(
					configs.RabbitmqTLSCA,
					configs.RabbitmqTLSCertificate,
					configs.RabbitmqTLSKey,
				),
			)
		} else {
			return rabbitmq.NewPubSub(configs.RabbitmqURL, queue, logger)
		}
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		if configs.NatsTLSCA != "" &&
		   configs.NatsTLSCertificate != "" &&
		   configs.NatsTLSKey != "" {
			return nats.NewPubSub(
				configs.NatsURL, 
				queue, 
				logger,
				nats.WithPubSubTLS(
					configs.NatsTLSCA,
					configs.NatsTLSCertificate,
					configs.NatsTLSKey,
				),
			)
		} else {
			return nats.NewPubSub(configs.NatsURL, queue, logger)
		}
	} else {
		fmt.Println("Invalid messaging system type for creating a pubsub:", systemType)
		return nil, errors.New("Invalid queue type")
	}
}

// Returns the string representing all channels
func GetAllChannels() (string) {
	systemType := queueConfiguration.GetSystem()

	if systemType == rabbitmqSystem {
		return rabbitmq.SubjectAllChannels;
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		return nats.SubjectAllChannels;
	} else {
		fmt.Println("Invalid messaging system type for creating a pubsub:", systemType)
		return "";
	}
}
