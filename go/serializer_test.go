package wampprotoflatbuffers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	wampprotoflatbuffers "github.com/xconnio/wampproto-flatbuffers/go"
	"github.com/xconnio/wampproto-go/messages"
)

func TestAbort(t *testing.T) {
	abort := messages.NewAbort(map[string]any{}, "wamp.error.close_realm", []any{"foo"}, map[string]any{"bar": "baz"})

	serializer := &wampprotoflatbuffers.FlatbuffersSerializer{}

	payload, err := serializer.Serialize(abort)
	require.NoError(t, err)
	require.NotNil(t, payload)

	require.Equal(t, uint8(messages.MessageTypeAbort), payload[0])

	msg, err := serializer.Deserialize(payload)
	require.NoError(t, err)
	require.NotNil(t, msg)
	recreated := msg.(*messages.Abort)

	require.Equal(t, abort.Reason(), recreated.Reason())
	require.Equal(t, abort.Args(), recreated.Args())
	require.Equal(t, abort.KwArgs(), recreated.KwArgs())
}
