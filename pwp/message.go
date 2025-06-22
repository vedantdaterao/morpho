package pwp

import (
	"encoding/hex"
	"fmt"
)

func SerializeHandshake(InfoHash []byte, peerID []byte) []byte {
	buf := make([]byte, HandshakeLen)
	buf[0] = byte(ProtocolStrLen)
	copy(buf[1:], []byte(ProtocolStr))
	copy(buf[1+len(ProtocolStr):], make([]byte, 8)) 					// reserved (zero bytes)
	copy(buf[1+len(ProtocolStr)+8:], InfoHash[:])
	copy(buf[1+len(ProtocolStr)+8+20:], peerID[:])
	fmt.Print("Handshake Data -\n",hex.Dump(buf))
	return buf
}

