from wampproto import messages, serializers

from wampprotofbs.parsers.abort import abort_to_flatbuffer, flatbuffer_to_abort


class FlatBuffersSerializer(serializers.Serializer):
    def serialize(self, message: messages.Message) -> bytes:
        if isinstance(message, messages.Abort):
            return abort_to_flatbuffer(message)
        else:
            raise TypeError(f"unknown message type {type(message)}")

    def deserialize(self, data: bytes) -> messages.Message:
        message_type = data[0]

        match message_type:
            case messages.Abort.TYPE:
                return flatbuffer_to_abort(data)
            case _:
                raise ValueError(f"Unsupported message type {message_type}")
