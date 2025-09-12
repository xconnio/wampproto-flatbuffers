package wampprotoflatbuffers

import (
	"fmt"

	"github.com/xconnio/wampproto-flatbuffers/go/parsers"
	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
)

type FlatbuffersSerializer struct{}

var _ serializers.Serializer = &FlatbuffersSerializer{}

func (c *FlatbuffersSerializer) Serialize(message messages.Message) ([]byte, error) {
	switch message.Type() {
	case messages.MessageTypeAbort:
		msg := message.(*messages.Abort)
		return parsers.AbortToFlatbuffers(msg)
	default:
		return nil, fmt.Errorf("unknown message type: %v", message.Type())
	}
}

func (c *FlatbuffersSerializer) Deserialize(data []byte) (messages.Message, error) {
	messageData, payloadData, err := parsers.ExtractMessage(data)
	if err != nil {
		return nil, err
	}

	switch uint64(data[0]) {
	case messages.MessageTypeAbort:
		return parsers.FlatbuffersToAbort(messageData, payloadData)
	default:
		return nil, fmt.Errorf("unknown message type: %v", data[0])
	}
}

func (c *FlatbuffersSerializer) Static() bool {
	return true
}
