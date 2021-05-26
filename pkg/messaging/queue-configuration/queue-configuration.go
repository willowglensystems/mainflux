package queueConfiguration

import (
	"errors"
	"strings"

	"github.com/mainflux/mainflux"
)

const (
	DefQueueSystem            = "nats"
	DefRabbitmqURL            = "amqp://guest:guest@rabbitmq/"
	DefRabbitmqTLSCertificate = ""
	DefRabbitmqTLSKey         = ""
	DefRabbitmqTLSCA          = ""
	DefRabbitmqDurable        = "false"
	DefRabbitmqTTL            = "3600000" // Milliseconds
	DefRabbitmqPriority       = "1"
	DefRabbitmqContentType    = "application/json"
	DefRabbitmqSubSystem      = "smc"
	DefRabbitmqSeverity       = "6"
	DefRabbitmqQueues         = ""
	DefNatsURL                = "nats://nats:4222"
	DefNatsTLSCertificate     = ""
	DefNatsTLSKey             = ""
	DefNatsTLSCA              = ""

	EnvQueueSystem            = "MF_QUEUE_SYSTEM"
	EnvRabbitmqURL            = "MF_RABBITMQ_URL"
	EnvRabbitmqTLSCertificate = "MF_RABBITMQ_TLS_CERTIFICATE"
	EnvRabbitmqTLSKey         = "MF_RABBITMQ_TLS_KEY"
	EnvRabbitmqTLSCA          = "MF_RABBITMQ_TLS_CA"
	EnvRabbitmqDurable        = "MF_RABBITMQ_DURABLE"
	EnvRabbitmqTTL            = "MF_RABBITMQ_TTL"
	EnvRabbitmqPriority       = "MF_RABBITMQ_PRIORITY"
	EnvRabbitmqContentType    = "MF_RABBITMQ_CONTENT_TYPE"
	EnvRabbitmqSubSystem      = "MF_RABBITMQ_SUB_SYSTEM"
	EnvRabbitmqSeverity       = "MF_RABBITMQ_SEVERITY"
	EnvRabbitmqQueues         = "MF_RABBITMQ_QUEUES"
	EnvNatsURL                = "MF_NATS_URL"
	EnvNatsTLSCertificate     = "MF_NATS_TLS_CERTIFICATE"
	EnvNatsTLSKey             = "MF_NATS_TLS_KEY"
	EnvNatsTLSCA              = "MF_NATS_TLS_CA"

	NatsMessagingSystem       = "nats"
	RabbitmqMessagingSystem   = "rabbitmq"
)

// GetSystem returns the type of queue system. 
// If no queue system is specified, it returns the default system.
func GetSystem() string {
	systemType := mainflux.Env(EnvQueueSystem, DefQueueSystem)
	return systemType
}

// GetConfig take the queue system type and returns two maps and an error. 
// The first parameter is a map of the system configuration parameters.
// The second parameter is a map that contains the queue parameters. If the queue system does not have any queue, it returns nil.
// The third parameter is the error in the case that the provided systemType is not valid.
func GetConfig() (map[string]string, map[string]string, error) {
	configs := make(map[string]string)
	systemType := mainflux.Env(EnvQueueSystem, DefQueueSystem)

	if systemType == NatsMessagingSystem {
		configs[EnvNatsURL] = mainflux.Env(EnvNatsURL, DefNatsURL)
		configs[EnvNatsTLSCertificate] = mainflux.Env(EnvNatsTLSCertificate, DefNatsTLSCertificate)
		configs[EnvNatsTLSKey] = mainflux.Env(EnvNatsTLSKey, DefNatsTLSKey)
		configs[EnvNatsTLSCA] = mainflux.Env(EnvNatsTLSCA, DefNatsTLSCA)
		return configs, nil, nil
	} else if systemType == RabbitmqMessagingSystem {
		configs[EnvRabbitmqURL] = mainflux.Env(EnvRabbitmqURL, DefRabbitmqURL)
		configs[EnvRabbitmqTLSCertificate] = mainflux.Env(EnvRabbitmqTLSCertificate, DefRabbitmqTLSCertificate)
		configs[EnvRabbitmqTLSKey] = mainflux.Env(EnvRabbitmqTLSKey, DefRabbitmqTLSKey)
		configs[EnvRabbitmqTLSCA] = mainflux.Env(EnvRabbitmqTLSCA, DefRabbitmqTLSCA)
		configs[EnvRabbitmqDurable] = mainflux.Env(EnvRabbitmqDurable, DefRabbitmqDurable)
		configs[EnvRabbitmqTTL] = mainflux.Env(EnvRabbitmqTTL, DefRabbitmqTTL)
		configs[EnvRabbitmqPriority] = mainflux.Env(EnvRabbitmqPriority, DefRabbitmqPriority)
		configs[EnvRabbitmqContentType] = mainflux.Env(EnvRabbitmqContentType, DefRabbitmqContentType)
		configs[EnvRabbitmqSubSystem] = mainflux.Env(EnvRabbitmqSubSystem, DefRabbitmqSubSystem)
		configs[EnvRabbitmqSeverity] = mainflux.Env(EnvRabbitmqSeverity, DefRabbitmqSeverity)

		queue := make(map[string]string)
		qString := mainflux.Env(EnvRabbitmqQueues, DefRabbitmqQueues)
		if qString != "" {
			items := strings.Split(qString, ",")
			for _, item := range items {
				parts := strings.Split(item, ":")
				queue[parts[0]] = parts[1]
			}
		}

		return configs, queue, nil
	} else {
		return nil, nil, errors.New("Queue type is not valid")
	}
}
