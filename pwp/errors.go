package pwp

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidHandshake = errors.New("invalid handshake")
	ErrInfoHashMismatch = errors.New("infohash mismatch")
	ErrMessageLen 		= errors.New("invalid message length")
)

type ConnectionErr struct {
	msg string
}

func (e *ConnectionErr) Error() string {
	return fmt.Sprintf("connnection error: %s", e.msg)
}