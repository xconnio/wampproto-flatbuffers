package parsers

import (
	"encoding/binary"
	"fmt"

	"github.com/xconnio/wampproto-go/serializers"
)

const HeaderLength = 3

func PrependHeader(messageType uint64, messageData, payloadData []byte) []byte {
	payloadLen := len(payloadData)
	totalLen := HeaderLength + len(messageData) + payloadLen

	result := make([]byte, totalLen)

	result[0] = uint8(messageType)                                               //nolint:gosec
	binary.BigEndian.PutUint16(result[1:HeaderLength], uint16(len(messageData))) //nolint:gosec

	copy(result[HeaderLength:], messageData)

	if payloadLen > 0 {
		copy(result[HeaderLength+len(messageData):], payloadData)
	}

	return result
}

func ExtractMessage(data []byte) ([]byte, []byte, error) {
	if len(data) < HeaderLength {
		return nil, nil, fmt.Errorf("invalid message length, must be at least %d bytes", HeaderLength)
	}

	messageLength := binary.BigEndian.Uint16(data[1:HeaderLength])
	if len(data) < HeaderLength+int(messageLength) {
		return nil, nil, fmt.Errorf("invalid message length")
	}

	messageData := data[HeaderLength : HeaderLength+int(messageLength)]
	payloadData := data[HeaderLength+int(messageLength):]

	return messageData, payloadData, nil
}

type PayloadExpander struct {
	expanded   bool
	payload    []byte
	serializer uint64

	args   []any
	kwargs map[string]any
}

func (p *PayloadExpander) NewPayloadExpander(serializer uint64, payload []byte) *PayloadExpander {
	return &PayloadExpander{
		serializer: serializer,
		payload:    payload,
	}
}

func (p *PayloadExpander) Expand() error {
	args, kwargs, err := serializers.DeserializePayload(p.serializer, p.payload)
	if err != nil {
		return err
	}

	p.args = args
	p.kwargs = kwargs
	p.expanded = true
	return nil
}

func (p *PayloadExpander) Args() []any {
	if !p.expanded {
		if err := p.Expand(); err != nil {
			return nil
		}
	}

	return p.args
}

func (p *PayloadExpander) Kwargs() map[string]any {
	if !p.expanded {
		if err := p.Expand(); err != nil {
			return nil
		}
	}

	return p.kwargs
}

func (p *PayloadExpander) Payload() []byte {
	return p.payload
}

func selectPayloadSerializer(options map[string]any) uint64 {
	payloadSerializer, ok := options["x_payload_serializer"].(uint64)
	if !ok {
		payloadSerializer = serializers.CBORSerializerID
	}

	return payloadSerializer
}
