package queueConfiguration

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mainflux/mainflux"
)

const (
	defQueueSystem         = "nats"
	defRabbitmqURL         = "amqp://guest:guest@rabbitmq/"
	defRabbitmqTLSConfig   = ""
	defRabbitmqDurable     = "true"
	defRabbitmqTTL         = "345600"
	defRabbitmqPriority    = "5"
	defRabbitmqContentType = "applicationnn/json"
	defRabbitmqSubSystem   = "smc"
	defRabbitmqSeverity    = "6"
	defRabbitmqQueues      = ""
	defNatsURL             = "nats://nats:4222"
	defNatsTLSConfig       = ""

	EnvQueueSystem         = "MF_QUEUE_SYSTEM"
	EnvRabbitmqURL         = "MF_RABBITMQ_URL"
	EnvRabbitmqTLSConfig   = "MF_RABBITMQ_TLS_CONFIG"
	EnvRabbitmqDurable     = "MF_RABBITMQ_DURABLE"
	EnvRabbitmqTTL         = "MF_RABBITMQ_TTL"
	EnvRabbitmqPriority    = "MF_RABBITMQ_PRIORITY"
	EnvRabbitmqContentType = "MF_RABBITMQ_CONTENT_TYPE"
	EnvRabbitmqSubSystem   = "MF_RABBITMQ_SUB_SYSTEM"
	EnvRabbitmqSeverity    = "MF_RABBITMQ_SEVERITY"
	EnvRabbitmqQueues      = "MF_RABBITMQ_QUEUES"
	EnvNatsURL             = "MF_NATS_URL"
	EnvNatsTLSConfig       = "MF_NATS_TLS_CONFIG"
)

// GetSystem returns the type of queue system. 
// If no queue system is specified, it returns the default system.
func GetSystem() string {
	systemType := mainflux.Env(EnvQueueSystem, defQueueSystem)
	fmt.Println("Queue system is ", systemType)

	return systemType
}

// GetConfig take the queue system type and returns two maps and an error. 
// The first parameter is a map of the system configuration parameters.
// The second parameter is a map that contains the queue parameters. If the queue system does not have any queue, it returns nil.
// The third parameter is the error in the case that the provided systemType is not valid.
func GetConfig() (map[string]string, map[string]string, error) {
	configs := make(map[string]string)
	systemType := mainflux.Env(EnvQueueSystem, defQueueSystem)

	if systemType == "nats" {
		configs[EnvNatsURL] = mainflux.Env(EnvNatsURL, defNatsURL)
		configs[EnvNatsTLSConfig] = mainflux.Env(EnvNatsTLSConfig, defNatsTLSConfig)
		return configs, nil, nil
	} else if systemType == "rabbitmq" {
		configs[EnvRabbitmqURL] = mainflux.Env(EnvRabbitmqURL, defRabbitmqURL)
		configs[EnvRabbitmqTLSConfig] = mainflux.Env(EnvRabbitmqTLSConfig, defRabbitmqTLSConfig)
		configs[EnvRabbitmqDurable] = mainflux.Env(EnvRabbitmqDurable, defRabbitmqDurable)
		configs[EnvRabbitmqTTL] = mainflux.Env(EnvRabbitmqTTL, defRabbitmqTTL)
		configs[EnvRabbitmqPriority] = mainflux.Env(EnvRabbitmqPriority, defRabbitmqPriority)
		configs[EnvRabbitmqContentType] = mainflux.Env(EnvRabbitmqContentType, defRabbitmqContentType)
		configs[EnvRabbitmqSubSystem] = mainflux.Env(EnvRabbitmqSubSystem, defRabbitmqSubSystem)
		configs[EnvRabbitmqSeverity] = mainflux.Env(EnvRabbitmqSeverity, defRabbitmqSeverity)

		queue := make(map[string]string)
		qString := mainflux.Env(EnvRabbitmqQueues, defRabbitmqQueues)
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
