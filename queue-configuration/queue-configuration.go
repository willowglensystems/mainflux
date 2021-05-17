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
		configs["MF_NATS_URL"] = mainflux.Env(envNatsURL, defNatsURL)
		configs["MF_NATS_TLS_CONFIG"] = mainflux.Env(envNatsTLSConfig, defNatsTLSConfig)
		return configs, nil, nil
	} else if systemType == "rabbitmq" {
		configs["MF_RABBITMQ_URL"] = mainflux.Env(envRabbitmqURL, defRabbitmqURL)
		configs["MF_RABBITMQ_TLS_CONFIG"] = mainflux.Env(envRabbitmqTLSConfig, defRabbitmqTLSConfig)
		configs["MF_RABBITMQ_DURABLE"] = mainflux.Env(envRabbitmqDurable, defRabbitmqDurable)
		configs["MF_RABBITMQ_TTL"] = mainflux.Env(envRabbitmqTTL, defRabbitmqTTL)
		configs["MF_RABBITMQ_PRIORITY"] = mainflux.Env(envRabbitmqPriority, defRabbitmqPriority)
		configs["MF_RABBITMQ_CONTENT_TYPE"] = mainflux.Env(envRabbitmqContentType, defRabbitmqContentType)
		configs["MF_RABBITMQ_SUB_SYSTEM"] = mainflux.Env(envRabbitmqSubSystem, defRabbitmqSubSystem)
		configs["MF_RABBITMQ_SEVERITY"] = mainflux.Env(envRabbitmqSeverity, defRabbitmqSeverity)
		
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
