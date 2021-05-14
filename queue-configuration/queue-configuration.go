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

	envQueueSystem         = "MF_QUEUE_SYSTEM"
	envRabbitmqURL         = "MF_RABBITMQ_URL"
	envRabbitmqTLSConfig   = "MF_RABBITMQ_TLS_CONFIG"
	envRabbitmqDurable     = "MF_RABBITMQ_DURABLE"
	envRabbitmqTTL         = "MF_RABBITMQ_TTL"
	envRabbitmqPriority    = "MF_RABBITMQ_PRIORITY"
	envRabbitmqContentType = "MF_RABBITMQ_CONTENT_TYPE"
	envRabbitmqSubSystem   = "MF_RABBITMQ_SUB_SYSTEM"
	envRabbitmqSeverity    = "MF_RABBITMQ_SEVERITY"
	envRabbitmqQueues      = "MF_RABBITMQ_QUEUES"
	envNatsURL             = "MF_NATS_URL"
	envNatsTLSConfig       = "MF_NATS_TLS_CONFIG"
)

// GetSystem returns the type of queue system. 
// If no queue system is specified, it returns the default system.
func GetSystem() string {
	systemType := mainflux.Env(envQueueSystem, defQueueSystem)
	fmt.Println("Queue system is ", systemType)

	return systemType
}

// GetConfig take the queue system type and returns two maps and an error. 
// The first parameter is a map of the system configuration parameters.
// The second parameter is a map that contains the queue parameters. If the queue system does not have any queue, it returns nil.
// The third parameter is the error in the case that the provided systemType is not valid.
func GetConfig(systemType string) (map[string]string, map[string]string, error) {
	configs := make(map[string]string)

	if systemType == "nats" {
		configs["natsURL"] = mainflux.Env(envNatsURL, defNatsURL)
		configs["natsTLSConfig"] = mainflux.Env(envNatsTLSConfig, defNatsTLSConfig)
		return configs, nil, nil
	} else if systemType == "rabbitmq" {
		configs["rabbitmqURL"] = mainflux.Env(envRabbitmqURL, defRabbitmqURL)
		configs["rabbitmqTLSConfig"] = mainflux.Env(envRabbitmqTLSConfig, defRabbitmqTLSConfig)
		configs["rabbitmqDurable "] = mainflux.Env(envRabbitmqDurable, defRabbitmqDurable)
		configs["rabbitmqTTL"] = mainflux.Env(envRabbitmqTTL, defRabbitmqTTL)
		configs["rabbitmqPriority"] = mainflux.Env(envRabbitmqPriority, defRabbitmqPriority)
		configs["rabbitmqContentType"] = mainflux.Env(envRabbitmqContentType, defRabbitmqContentType)
		configs["rabbitmqSubSystem"] = mainflux.Env(envRabbitmqSubSystem, defRabbitmqSubSystem)
		configs["rabbitmqSeverity"] = mainflux.Env(envRabbitmqSeverity, defRabbitmqSeverity)
		
		queue := make(map[string]string)
		qString := mainflux.Env(envRabbitmqQueues, defRabbitmqQueues)
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
