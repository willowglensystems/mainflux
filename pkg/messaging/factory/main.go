package factory

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/messaging"
	"github.com/mainflux/mainflux/pkg/messaging/nats"
	"github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
	"github.com/mainflux/mainflux/pkg/messaging/rabbitmq"
)

// Generates a NATS or RabbitMQ publisher based on the messaging queue configuration
func NewPublisher() (messaging.Publisher, error) {
	systemType := queueConfiguration.GetSystem()
	configs, _, _ := queueConfiguration.GetConfig()

	if systemType == queueConfiguration.RabbitmqMessagingSystem {
		return rabbitmq.NewPublisher(configs.RabbitmqURL, nil)
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		return nats.NewPublisher(configs.NatsURL, nil)
	} else {
		fmt.Println("Invalid messaging system type for creating a publisher:", systemType)
		return nil, errors.New("Invalid queue type")
	}
}

// Generates a NATS or RabbitMQ publisher with TLS based on the messaging queue configuration
func NewPublisherTLS() (messaging.Publisher, error) {
	systemType := queueConfiguration.GetSystem()
	configs, _, _ := queueConfiguration.GetConfig()

	if !isTLSConfigured(systemType, configs) {
		fmt.Println("Invalid TLS configuration for creating a TLS publisher:", systemType)
		return nil, errors.New("Invalid TLS configuration")
	}

	if systemType == queueConfiguration.RabbitmqMessagingSystem {
		return rabbitmq.NewPublisher(
			configs.RabbitmqURL, 
			generateTLSConfig(
				configs.RabbitmqTLSCA,
				configs.RabbitmqTLSCertificate,
				configs.RabbitmqTLSKey,
			),
		)
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		return nats.NewPublisher(
			configs.NatsURL, 
			generateTLSConfig(
				configs.NatsTLSCA,
				configs.NatsTLSCertificate,
				configs.NatsTLSKey,
			),
		)
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
		return rabbitmq.NewPubSub(configs.RabbitmqURL, queue, logger, nil)
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		return nats.NewPubSub(configs.NatsURL, queue, logger, nil)
	} else {
		fmt.Println("Invalid messaging system type for creating a pubsub:", systemType)
		return nil, errors.New("Invalid queue type")
	}
}

// Generates a NATS or RabbitMQ publisher/subscriber with TLS based on the messaging queue configuration
func NewPubSubTLS(queue string, logger log.Logger) (messaging.PubSub, error) {
	systemType := queueConfiguration.GetSystem()
	configs, _, _ := queueConfiguration.GetConfig()

	if !isTLSConfigured(systemType, configs) {
		fmt.Println("Invalid TLS configuration for creating a TLS publisher:", systemType)
		return nil, errors.New("Invalid TLS configuration")
	}

	if systemType == queueConfiguration.RabbitmqMessagingSystem {
		return rabbitmq.NewPubSub(
			configs.RabbitmqURL, 
			queue,
			logger,
			generateTLSConfig(
				configs.RabbitmqTLSCA,
				configs.RabbitmqTLSCertificate,
				configs.RabbitmqTLSKey,
			),
		)
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		return nats.NewPubSub(
			configs.NatsURL, 
			queue, 
			logger,
			generateTLSConfig(
				configs.NatsTLSCA,
				configs.NatsTLSCertificate,
				configs.NatsTLSKey,
			),
		)
	} else {
		fmt.Println("Invalid messaging system type for creating a pubsub:", systemType)
		return nil, errors.New("Invalid queue type")
	}
}

// Returns the string representing all channels
func GetAllChannels() (string) {
	systemType := queueConfiguration.GetSystem()

	if systemType == queueConfiguration.RabbitmqMessagingSystem {
		return rabbitmq.SubjectAllChannels;
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		return nats.SubjectAllChannels;
	} else {
		fmt.Println("Invalid messaging system type for creating a pubsub:", systemType)
		return "";
	}
}

func isTLSConfigured(systemType string, configs *queueConfiguration.Config) bool {
	if systemType == queueConfiguration.RabbitmqMessagingSystem {
		return configs.RabbitmqTLSCA != "" && configs.RabbitmqTLSCertificate != "" && configs.RabbitmqTLSKey != "";
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		return configs.NatsTLSCA != "" && configs.NatsTLSCertificate != "" && configs.NatsTLSKey != "";
	} else {
		fmt.Println("Invalid messaging system type for checking TLS configuration:", systemType)
		return false
	}
}

// Takes in the file paths of a TLS Certificate Authority and a TLS Client Certificate and Key to generate a tls.Config object.
func generateTLSConfig(caFile, certFile, keyFile string) (*tls.Config) {
	roots := x509.NewCertPool()
	data, err := ioutil.ReadFile(caFile)

	if err != nil {
		fmt.Errorf("Error reading in root CA: %s", err)
		return nil
	}

	roots.AppendCertsFromPEM(data)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		fmt.Errorf("Error loading certificate: %s", err)
		return nil
	}

	return &tls.Config {
		Certificates: []tls.Certificate{cert},
		RootCAs: roots,
		InsecureSkipVerify: false,
	}
}