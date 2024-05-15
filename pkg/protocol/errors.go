package protocol

import (
	"errors"
)

// ErrInvalidHeaderLength is the error when the header length is invalid
var ErrInvalidHeaderLength = errors.New("invalid header length")

// ErrInvalidPayloadLength is the error when the payload length is invalid
var ErrInvalidPayloadLength = errors.New("invalid payload length")
