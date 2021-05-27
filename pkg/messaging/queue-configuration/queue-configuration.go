package queueConfiguration

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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

type Config struct {
	RabbitmqURL            string
	RabbitmqTLSCertificate string
	RabbitmqTLSKey         string
	RabbitmqTLSCA          string
	RabbitmqDurable        bool
	RabbitmqTTL            time.Duration
	RabbitmqPriority       uint8
	RabbitmqContentType    string
	RabbitmqSubSystem      string
	RabbitmqSeverity       uint8
	NatsURL                string
	NatsTLSCertificate     string
	NatsTLSKey             string
	NatsTLSCA              string
}

// GetSystem returns the type of queue system. 
// If no queue system is specified, it returns the default system.
func GetSystem() string {
	systemType := mainflux.Env(EnvQueueSystem, DefQueueSystem)
	return systemType
}

// GetConfig retrieves and parses the system configuration from environment variables and returns them in a struct and a map.
// If an error occurs, then the struct and map are nils and an error is returned.
// The first parameter is a struct of the system configuration parameters.
// The second parameter is a map that contains the queue parameters. If the queue system does not have any queue, it returns nil.
// The third parameter is the error in the case that the provided systemType is not valid.
func GetConfig() (*Config, map[string]string, error) {
	systemType := mainflux.Env(EnvQueueSystem, DefQueueSystem)

	if systemType == NatsMessagingSystem {
		config := &Config {
			NatsURL: mainflux.Env(EnvNatsURL, DefNatsURL),
			NatsTLSCertificate: mainflux.Env(EnvNatsTLSCertificate, DefNatsTLSCertificate),
			NatsTLSKey: mainflux.Env(EnvNatsTLSKey, DefNatsTLSKey),
			NatsTLSCA: mainflux.Env(EnvNatsTLSCA, DefNatsTLSCA),
		}

		return config, nil, nil
	} else if systemType == RabbitmqMessagingSystem {
		durableValue, err := strconv.ParseBool(mainflux.Env(EnvRabbitmqDurable, DefRabbitmqDurable))

		if err != nil {
			fmt.Println("Unable to parse Durable configuration, defaulting to false")
			durableValue, _ = strconv.ParseBool(DefRabbitmqDurable)
		}
	
		priorityValue, err := strconv.ParseUint(mainflux.Env(EnvRabbitmqPriority, DefRabbitmqPriority), 10, 64)
	
		if err != nil {
			fmt.Println("Unable to parse Priority configuration, defaulting to 1")
			priorityValue, _ = strconv.ParseUint(DefRabbitmqPriority, 10, 8)
		}
	
		ttlValue, err := strconv.ParseUint(mainflux.Env(EnvRabbitmqTTL, DefRabbitmqTTL), 10, 8)
	
		if err != nil {
			fmt.Println("Unable to parse TTL configuration, defaulting to 3600000 milliseconds")
			ttlValue, _ = strconv.ParseUint(DefRabbitmqTTL, 10, 64)
		}
	
		severityValue, err := strconv.ParseUint(mainflux.Env(EnvRabbitmqSeverity, DefRabbitmqSeverity), 10, 64)
	
		if err != nil {
			fmt.Println("Unable to parse Severity configuration, defaulting to 6")
			severityValue, _ = strconv.ParseUint(DefRabbitmqSeverity, 10, 8)
		}

		queue := make(map[string]string)
		qString := mainflux.Env(EnvRabbitmqQueues, DefRabbitmqQueues)
		if qString != "" {
			items := strings.Split(qString, ",")
			for _, item := range items {
				parts := strings.Split(item, ":")
				queue[parts[0]] = parts[1]
			}
		}

		config := &Config{
			RabbitmqURL: mainflux.Env(EnvRabbitmqURL, DefRabbitmqURL),
			RabbitmqTLSCertificate: mainflux.Env(EnvRabbitmqTLSCertificate, DefRabbitmqTLSCertificate),
			RabbitmqTLSKey: mainflux.Env(EnvRabbitmqTLSKey, DefRabbitmqTLSKey),
			RabbitmqTLSCA: mainflux.Env(EnvRabbitmqTLSCA, DefRabbitmqTLSCA),
			RabbitmqDurable: durableValue,
			RabbitmqTTL: time.Duration(ttlValue) * time.Millisecond,
			RabbitmqPriority: uint8(priorityValue),
			RabbitmqContentType: mainflux.Env(EnvRabbitmqContentType, DefRabbitmqContentType),
			RabbitmqSubSystem: mainflux.Env(EnvRabbitmqSubSystem, DefRabbitmqSubSystem),
			RabbitmqSeverity: uint8(severityValue),
		}

		return config, queue, nil
	} else {
		return nil, nil, errors.New("Queue type is not valid")
	}
}
