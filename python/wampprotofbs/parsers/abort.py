from typing import Any

import flatbuffers
from wampproto.messages.abort import Abort, IAbortFields
from wampproto.serializers.payload import serialize_payload

from wampprotofbs.parsers import helpers
from wampprotofbs.gen import Abort as Abort_FBS


class AbortFields(IAbortFields):
    def __init__(self, gen: Abort_FBS.Abort, payload: bytes):
        self._gen = gen
        self._ex = helpers.PayloadExpander(payload, gen.PayloadSerializerId())

    @property
    def reason(self) -> str:
        return self._gen.Reason().decode("utf-8")

    @property
    def details(self) -> dict[str, Any]:
        return {}

    @property
    def args(self) -> list[Any]:
        return self._ex.args

    @property
    def kwargs(self) -> dict[str, Any]:
        return self._ex.kwargs

    @property
    def payload_serializer(self) -> int:
        return self._gen.PayloadSerializerId()

    @property
    def payload(self) -> bytes:
        return self._ex.payload


def abort_to_flatbuffer(a: Abort) -> bytes:
    builder = flatbuffers.Builder(1024)

    reason_offset = builder.CreateString(a.reason)
    payload_serializer = helpers.select_payload_serializer(a.details)

    Abort_FBS.AbortStart(builder)
    Abort_FBS.AbortAddReason(builder, reason_offset)
    Abort_FBS.AbortAddPayloadSerializerId(builder, payload_serializer)
    abort_msg = Abort_FBS.AbortEnd(builder)

    builder.Finish(abort_msg)
    payload = serialize_payload(payload_serializer, a.args, a.kwargs)

    return helpers.prepend_header(Abort.TYPE, bytes(builder.Output()), payload)


def flatbuffer_to_abort(data: bytes) -> Abort:
    message_data, payload_data = helpers.extract_message(data)
    fb_abort = Abort_FBS.Abort.GetRootAsAbort(message_data, 0)

    return Abort(AbortFields(fb_abort, payload_data))
