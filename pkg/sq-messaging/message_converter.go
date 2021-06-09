package sqMessaging

import (
	"encoding/binary"

	"git.willowglen.ca/sq/third-party/mainflux/pkg/messaging"
)

func ToSQMessage(msg messaging.Message) SQMessage {
	sqPriority, _ := binary.Varint(msg.Metadata["Priority"])

	sqMetadata := make(map[string][]byte)
	for key, element := range msg.Metadata {
		sqMetadata[key] = element
	}

	sq_msg := SQMessage{
		Source:      msg.Publisher,
		Timestamp:   msg.Created,
		Type:        string(msg.Metadata["Type"]),
		Version:     string(msg.Metadata["Version"]),
		ContentType: string(msg.Metadata["ContentType"]),
		Body:        msg.Payload,
		Token:       string(msg.Metadata["Token"]),
		Priority:    sqPriority,
		Initiator:   string(msg.Metadata["Initiator"]),
		Metadata:    sqMetadata,
	}

	return sq_msg
}

func ToMFMessage(msg SQMessage) messaging.Message {
	mfPriority := make([]byte, 4)
	binary.PutVarint(mfPriority, msg.Priority)

	mfMetadata := make(map[string][]byte)
	for key, element := range msg.Metadata {
		mfMetadata[key] = element
	}

	mfMetadata["Type"] = []byte(msg.Type)
	mfMetadata["Version"] = []byte(msg.Version)
	mfMetadata["ContentType"] = []byte(msg.ContentType)
	mfMetadata["Token"] = []byte(msg.Token)
	mfMetadata["Priority"] = mfPriority
	mfMetadata["Initiator"] = []byte(msg.Initiator)

	mf_msg := messaging.Message{
		Publisher: msg.Source,
		Payload:   msg.Body,
		Created:   msg.Timestamp,
		Metadata:  mfMetadata,
	}

	return mf_msg
}
