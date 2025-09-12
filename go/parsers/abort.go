package parsers

import (
	"github.com/google/flatbuffers/go"

	"github.com/xconnio/wampproto-flatbuffers/go/gen"
	"github.com/xconnio/wampproto-go/messages"
	"github.com/xconnio/wampproto-go/serializers"
)

type Abort struct {
	gen *gen.Abort
	ex  *PayloadExpander
}

func NewAbortFields(g *gen.Abort, payload []byte) messages.AbortFields {
	return &Abort{
		gen: g,
		ex:  &PayloadExpander{payload: payload, serializer: g.PayloadSerializerId()},
	}
}

func (a *Abort) Reason() string {
	return string(a.gen.Reason())
}

func (a *Abort) Details() map[string]any {
	return map[string]any{}
}

func (a *Abort) Args() []any {
	return a.ex.Args()
}

func (a *Abort) KwArgs() map[string]any {
	return a.ex.Kwargs()
}

func AbortToFlatbuffers(m *messages.Abort) ([]byte, error) {
	builder := flatbuffers.NewBuilder(256)

	reason := builder.CreateString(m.Reason())

	gen.AbortStart(builder)
	gen.AbortAddReason(builder, reason)
	payloadSerializer := selectPayloadSerializer(m.Details())
	gen.AbortAddPayloadSerializerId(builder, payloadSerializer)
	abort := gen.AbortEnd(builder)

	builder.Finish(abort)

	payload, err := serializers.SerializePayload(payloadSerializer, m.Args(), m.KwArgs())
	if err != nil {
		return nil, err
	}

	return PrependHeader(messages.MessageTypeAbort, builder.FinishedBytes(), payload), nil
}

func FlatbuffersToAbort(data, payload []byte) (*messages.Abort, error) {
	abort := gen.GetRootAsAbort(data, 0)
	return messages.NewAbortWithFields(NewAbortFields(abort, payload)), nil
}
