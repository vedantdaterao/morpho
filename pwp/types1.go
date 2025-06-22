package pwp

import (
	"bufio"
	"net"
)

// type Conn interface {
// 	ReadMsg() (error)
// 	SendMsg() (error)
// 	RequestPiece() (error)
// }

type Peer struct {
	addr 			net.Addr
	peer_id 		[]byte
	bitfield		[]byte

	connStatus 		bool
	reader 			*bufio.Reader
	conn 			net.Conn

	amChoking		bool
	amInterested 	bool
	peerChoking 	bool
	peerInterested 	bool

	ap 				ActivePieces
	// MsgChan 		chan 	interface{}
}