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
	queueConfiguration "github.com/mainflux/mainflux/pkg/messaging/queue-configuration"
	"github.com/mainflux/mainflux/pkg/messaging/rabbitmq"
)

// Generates a NATS or RabbitMQ publisher based on the messaging queue configuration.
// If TLS is configured then a TLS connection is made.
func NewPublisher() (messaging.Publisher, error) {
	systemType := queueConfiguration.GetSystem()
	configs, _, _ := queueConfiguration.GetConfig()

	if configs.EnableTLS {
		tlsConfig := generateTLSConfig(systemType, configs)

		if tlsConfig == nil {
			fmt.Println("Invalid TLS configuration for creating a TLS publisher:", systemType)
			return nil, errors.New("Invalid TLS configuration")
		}

		if systemType == queueConfiguration.RabbitmqMessagingSystem {
			return rabbitmq.NewPublisher(configs.RabbitmqURL, tlsConfig)
		} else if systemType == queueConfiguration.NatsMessagingSystem {
			return nats.NewPublisher(configs.NatsURL, tlsConfig)
		} else {
			fmt.Println("Invalid messaging system type for creating a publisher:", systemType)
			return nil, errors.New("Invalid queue type")
		}
	} else {
		if systemType == queueConfiguration.RabbitmqMessagingSystem {
			return rabbitmq.NewPublisher(configs.RabbitmqURL, nil)
		} else if systemType == queueConfiguration.NatsMessagingSystem {
			return nats.NewPublisher(configs.NatsURL, nil)
		} else {
			fmt.Println("Invalid messaging system type for creating a publisher:", systemType)
			return nil, errors.New("Invalid queue type")
		}
	}
}

// Generates a NATS or RabbitMQ publisher/subscriber based on the messaging queue configuration.
// If TLS is configured then a TLS connection is made.
func NewPubSub(queue string, logger log.Logger) (messaging.PubSub, error) {
	systemType := queueConfiguration.GetSystem()
	configs, _, _ := queueConfiguration.GetConfig()

	if configs.EnableTLS {
		tlsConfig := generateTLSConfig(systemType, configs)

		if tlsConfig == nil {
			fmt.Println("Invalid TLS configuration for creating a TLS publisher:", systemType)
			return nil, errors.New("Invalid TLS configuration")
		}

		if systemType == queueConfiguration.RabbitmqMessagingSystem {
			return rabbitmq.NewPubSub(configs.RabbitmqURL, queue, logger, tlsConfig)
		} else if systemType == queueConfiguration.NatsMessagingSystem {
			return nats.NewPubSub(configs.NatsURL, queue, logger, tlsConfig)
		} else {
			fmt.Println("Invalid messaging system type for creating a pubsub:", systemType)
			return nil, errors.New("Invalid queue type")
		}
	} else {
		if systemType == queueConfiguration.RabbitmqMessagingSystem {
			return rabbitmq.NewPubSub(configs.RabbitmqURL, queue, logger, nil)
		} else if systemType == queueConfiguration.NatsMessagingSystem {
			return nats.NewPubSub(configs.NatsURL, queue, logger, nil)
		} else {
			fmt.Println("Invalid messaging system type for creating a pubsub:", systemType)
			return nil, errors.New("Invalid queue type")
		}
	}
}

// Returns the string representing all channels
func GetAllChannels() string {
	systemType := queueConfiguration.GetSystem()

	if systemType == queueConfiguration.RabbitmqMessagingSystem {
		return rabbitmq.SubjectAllChannels
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		return nats.SubjectAllChannels
	} else {
		fmt.Println("Invalid messaging system type for creating a pubsub:", systemType)
		return ""
	}
}

// Takes in the file paths of a TLS Certificate Authority and a TLS Client Certificate and Key to generate a tls.Config object.
func generateTLSConfig(systemType string, configs *queueConfiguration.Config) *tls.Config {
	var caFile, certFile, keyFile string

	if systemType == queueConfiguration.RabbitmqMessagingSystem {
		caFile, certFile, keyFile = configs.RabbitmqTLSCA, configs.RabbitmqTLSCertificate, configs.RabbitmqTLSKey
	} else if systemType == queueConfiguration.NatsMessagingSystem {
		caFile, certFile, keyFile = configs.NatsTLSCA, configs.NatsTLSCertificate, configs.NatsTLSKey
	} else {
		fmt.Println("Invalid messaging system type for checking TLS configuration:", systemType)
		return nil
	}

	if caFile != "" && certFile != "" && keyFile != "" {
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

		return &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            roots,
			InsecureSkipVerify: false,
		}
	} else {
		if caFile == "" {
			fmt.Errorf("No Certificate Authority specified")
		}

		if certFile == "" {
			fmt.Errorf("No client certificate specified")
		}

		if keyFile == "" {
			fmt.Errorf("No client key file specified")
		}

		return nil
	}
}
