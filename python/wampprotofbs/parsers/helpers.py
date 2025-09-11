from typing import Any

from wampproto.serializers.cbor import CBOR_SERIALIZER_ID
from wampproto.serializers.payload import deserialize_payload

HEADER_LENGTH = 3


def select_payload_serializer(options: dict[str, Any]) -> int:
    if isinstance(options, dict) and "x_payload_serializer" in options:
        return int(options["x_payload_serializer"])

    return CBOR_SERIALIZER_ID


def prepend_header(message_type: int, message_data: bytes, payload_data: bytes) -> bytes:
    result = [message_type]

    result.extend(len(message_data).to_bytes(2, byteorder="big"))

    result.extend(message_data)

    if payload_data:
        result.extend(payload_data)

    return bytes(result)


def extract_message(data: bytes) -> tuple[bytes, bytes]:
    if len(data) < HEADER_LENGTH:
        raise ValueError(f"invalid message length, must be at least {HEADER_LENGTH} bytes")

    message_length = int.from_bytes(data[1:HEADER_LENGTH], byteorder="big")
    if len(data) < HEADER_LENGTH + message_length:
        raise ValueError("invalid message length")

    message_data = data[HEADER_LENGTH : HEADER_LENGTH + message_length]
    payload_data = data[HEADER_LENGTH + message_length :]

    return message_data, payload_data


class PayloadExpander:
    def __init__(self, payload: bytes, serializer: int):
        self._expanded = False
        self._payload = payload
        self._serializer = serializer
        self._args: list[Any] | None = None
        self._kwargs: dict[str, Any] | None = None

    def _expand(self) -> None:
        args, kwargs = deserialize_payload(self._serializer, self._payload)
        self._args = args
        self._kwargs = kwargs
        self._expanded = True

    @property
    def args(self) -> list[Any]:
        if not self._expanded:
            self._expand()

        return self._args

    @property
    def kwargs(self) -> dict[str, Any]:
        if not self._expanded:
            self._expand()

        return self._kwargs

    @property
    def payload(self) -> bytes:
        return self._payload
