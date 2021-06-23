package sqMessaging

import (
	// "fmt"
	"encoding/binary"
	"testing"
	"time"

	"github.com/mainflux/mainflux/pkg/messaging"
)

func TestToSQMessageEmpty(t *testing.T) {
	// test for empty MF message
	mfMsg := messaging.Message{}
	sqMsg := ToSQMessage(mfMsg)

	if sqMsg.Source != "" {
		t.Error("ToSGMessage with empty message failed, expected empty source, got", sqMsg.Source)
	}

	if sqMsg.Timestamp != 0 {
		t.Error("ToSGMessage with empty message failed, expected 0 timestamp, got", sqMsg.Timestamp)
	}

	if sqMsg.Type != "" {
		t.Error("ToSGMessage with empty message failed, expected empty type, got", sqMsg.Type)
	}

	if sqMsg.Version != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Version, got", sqMsg.Version)
	}

	if sqMsg.ContentType != "" {
		t.Error("ToSGMessage with empty message failed, expected empty ContentType, got", sqMsg.ContentType)
	}

	if len(sqMsg.Body) != 0 {
		t.Error("ToSGMessage with empty message failed, expected 0 length Body, got", sqMsg.Body)
	}

	if sqMsg.Token != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Token, got", sqMsg.Token)
	}

	if sqMsg.Priority != 0 {
		t.Error("ToSGMessage with empty message failed, expected 0 Priority, got", sqMsg.Priority)
	}

	if sqMsg.Initiator != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Initiator, got", sqMsg.Initiator)
	}

	if len(sqMsg.Metadata) != 0 {
		t.Error("ToSGMessage with empty message failed, expected 0 length Metadata, got", sqMsg.Metadata)
	}
}

func TestToSQMessage(t *testing.T) {
	// test for nil as string metadata map
	mfMsg := messaging.Message{
		Metadata: map[string][]byte{
			"Token": []byte(nil),
		},
	}
	sqMsg := ToSQMessage(mfMsg)

	if sqMsg.Token != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Token, got", sqMsg.Token)
	}

	// test for empty string in metadata map
	mfMsg = messaging.Message{
		Metadata: map[string][]byte{
			"Version": []byte(""),
		},
	}
	sqMsg = ToSQMessage(mfMsg)

	if sqMsg.Version != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Version, got", sqMsg.Version)
	}

	// test for negative value as priority in metadata map
	mfPriority := make([]byte, 4)
	binary.PutVarint(mfPriority, -4)

	mfMsg = messaging.Message{
		Metadata: map[string][]byte{
			"Priority": mfPriority,
		},
	}
	sqMsg = ToSQMessage(mfMsg)

	if sqMsg.Priority != -4 {
		t.Error("ToSGMessage with empty message failed, expected -4 as Priority, got", sqMsg.Priority)
	}

	// test for zero value as priority in metadata map
	mfPriority = make([]byte, 4)
	binary.PutVarint(mfPriority, 0)

	mfMsg = messaging.Message{
		Metadata: map[string][]byte{
			"Priority": mfPriority,
		},
	}
	sqMsg = ToSQMessage(mfMsg)

	if sqMsg.Priority != 0 {
		t.Error("ToSGMessage with empty message failed, expected 0 as Priority, got", sqMsg.Priority)
	}

	// test for positive value as priority in metadata map
	mfPriority = make([]byte, 4)
	binary.PutVarint(mfPriority, 5)

	mfMsg = messaging.Message{
		Metadata: map[string][]byte{
			"Priority": mfPriority,
		},
	}
	sqMsg = ToSQMessage(mfMsg)

	if sqMsg.Priority != 5 {
		t.Error("ToSGMessage with empty message failed, expected 5 as Priority, got", sqMsg.Priority)
	}

	// test for additional fields in metadata
	mfMsg = messaging.Message{
		Metadata: map[string][]byte{
			"TestField": []byte("SomeTestData"),
		},
	}
	sqMsg = ToSQMessage(mfMsg)
	testStr := string(sqMsg.Metadata["TestField"])

	if testStr != "SomeTestData" {
		t.Error("ToSGMessage with empty message failed, expected SomeTestData as TestField, got", testStr)
	}

	// test for timestamp
	creationTime := time.Now().UnixNano()
	mfMsg = messaging.Message{
		Created: creationTime,
	}
	sqMsg = ToSQMessage(mfMsg)

	if sqMsg.Timestamp != creationTime {
		t.Error("ToSGMessage with empty message failed, expected ", creationTime, " as Timestamp, got", sqMsg.Timestamp)
	}
}

func TestToMFMessageEmpty(t *testing.T) {
	// test for empty SQ message
	sqMsg := SQMessage{}
	mfMsg := ToMFMessage(sqMsg)

	if mfMsg.Publisher != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Publisher, got", mfMsg.Publisher)
	}

	if mfMsg.Created != 0 {
		t.Error("ToSGMessage with empty message failed, expected 0 Created, got", mfMsg.Created)
	}

	if len(mfMsg.Payload) != 0 {
		t.Error("ToSGMessage with empty message failed, expected 0 length Payload, got", mfMsg.Payload)
	}

	if len(mfMsg.Metadata) != 6 {
		t.Error("ToSGMessage with empty message failed, expected 6 length Metadata, got", mfMsg.Metadata)
	}

	if string(mfMsg.Metadata["Type"]) != "" {
		t.Error("ToSGMessage with empty message failed, expected empty type, got", mfMsg.Metadata["Type"])
	}

	if string(mfMsg.Metadata["Version"]) != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Version, got", string(mfMsg.Metadata["Version"]))
	}

	if string(mfMsg.Metadata["ContentType"]) != "" {
		t.Error("ToSGMessage with empty message failed, expected empty ContentType, got", string(mfMsg.Metadata["ContentType"]))
	}

	if string(mfMsg.Metadata["Token"]) != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Token, got", string(mfMsg.Metadata["Token"]))
	}

	priority, _ := binary.Varint(mfMsg.Metadata["Priority"])
	if priority != 0 {
		t.Error("ToSGMessage with empty message failed, expected 0 Priority, got", priority)
	}

	if string(mfMsg.Metadata["Initiator"]) != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Initiator, got", string(mfMsg.Metadata["Initiator"]))
	}

}

func TestToMFMessage(t *testing.T) {
	// test for empty string
	sqMsg := SQMessage{
		Type: "",
	}
	mfMsg := ToMFMessage(sqMsg)

	if string(mfMsg.Metadata["Type"]) != "" {
		t.Error("ToSGMessage with empty message failed, expected empty Type, got", string(mfMsg.Metadata["Type"]))
	}

	// test for negative value as priority
	sqMsg = SQMessage{
		Priority: -3,
	}
	mfMsg = ToMFMessage(sqMsg)

	priority, _ := binary.Varint(mfMsg.Metadata["Priority"])
	if priority != -3 {
		t.Error("ToSGMessage with empty message failed, expected -3 as Priority, got", priority)
	}

	// test for zero value as priority
	sqMsg = SQMessage{
		Priority: 0,
	}
	mfMsg = ToMFMessage(sqMsg)

	priority, _ = binary.Varint(mfMsg.Metadata["Priority"])
	if priority != 0 {
		t.Error("ToSGMessage with empty message failed, expected 0 as Priority, got", priority)
	}

	// test for positive value as priority
	sqMsg = SQMessage{
		Priority: 6,
	}
	mfMsg = ToMFMessage(sqMsg)

	priority, _ = binary.Varint(mfMsg.Metadata["Priority"])
	if priority != 6 {
		t.Error("ToSGMessage with empty message failed, expected 6 as Priority, got", priority)
	}

	// test for metadata fields
	sqMsg = SQMessage{
		Metadata: map[string][]byte{
			"TestField": []byte("TestData"),
		},
	}
	mfMsg = ToMFMessage(sqMsg)

	if len(mfMsg.Metadata) != 7 {
		t.Error("ToSGMessage with empty message failed, expected length 7 for Metadata, got", len(mfMsg.Metadata))
	}

	if string(mfMsg.Metadata["TestField"]) != "TestData" {
		t.Error("ToSGMessage with empty message failed, expected TestData as TestField, got", string(mfMsg.Metadata["TestField"]))
	}

	// test for created
	creationTime := time.Now().UnixNano()
	sqMsg = SQMessage{
		Timestamp: creationTime,
	}
	mfMsg = ToMFMessage(sqMsg)

	if mfMsg.Created != creationTime {
		t.Error("ToSGMessage with empty message failed, expected ", creationTime, " as Created, got", mfMsg.Created)
	}
}
