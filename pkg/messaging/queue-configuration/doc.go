// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package queueConfiguration holds the configuration parser for 
// messaging queue settings. The settings are set through environment
// variables and are defined as follows:
// MF_QUEUE_SYSTEM - The messaging queue system to use. (nats or rabbitmq)
// MF_RABBITMQ_URL - The URL used to connect to rabbitmq.
// MF_RABBITMQ_TLS_CERTFICATE - The file path to the TLS client certificate.
// MF_RABBITMQ_TLS_KEY - The file path to the TLS client certificate key.
// MF_RABBITMQ_TLS_CA - The file path to the TLS Certificate Authority used to validate the rabbitmq server.
// MF_RABBITMQ_DURABLE - A boolean for the durability of a rabbitmq message.
// MF_RABBITMQ_TTL - An integer for a rabbitmq message's time to live (in ms).
// MF_RABBITMQ_PRIORITY - An integer for the priority of a message.
// MF_RABBITMQ_CONTENT_TYPE - The media type of the message
// MF_RABBITMQ_SUB_SYSTEM - The name of the subsystem that is sending the message.
// MF_RABBITMQ_SEVERITY - An integer specifying the severity of a message.
// MF_RABBITMQ_QUEUES - Unneeded
// MF_NATS_URL - The URL used to connect to nats.
// MF_NATS_TLS_CERTFICATE - The file path to the TLS client certificate.
// MF_NATS_TLS_KEY - The file path to the TLS client certificate key.
// MF_NATS_TLS_CA - The file path to the TLS Certificate Authority used to validate the nats server.
package queueConfiguration
