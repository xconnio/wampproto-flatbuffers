from wampproto import messages

from wampprotofbs.serializer import FlatBuffersSerializer


def test_serialize_deserialize():
    serializer = FlatBuffersSerializer()
    hello = messages.Abort(messages.AbortFields({}, "reason", ["args"], {"k": "va"}))
    data = serializer.serialize(hello)
    message: messages.Abort = serializer.deserialize(data)
    assert message.TYPE == hello.TYPE
    assert message.details == hello.details
    assert message.reason == hello.reason
    assert message.args == hello.args
    assert message.kwargs == hello.kwargs
