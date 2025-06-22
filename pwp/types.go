package pwp

import (
	"bufio"
	"net"
	"sync"
)

const (
	ProtocolStr    = "BitTorrent protocol"
	ProtocolStrLen = len(ProtocolStr)      // 19
	HandshakeLen   = 49 + len(ProtocolStr) // total: 68 bytes
)

const (
	Msg_Choke         byte = 0
	Msg_Unchoke       byte = 1
	Msg_Interested    byte = 2
	Msg_NotInterested byte = 3
	Msg_Have          byte = 4
	Msg_Bitfield      byte = 5
	Msg_Request       byte = 6
	Msg_Piece         byte = 7
	Msg_Cancel        byte = 8
)

// handshake: <pstrlen><pstr><reserved><info_hash><peer_id>
type Msg_handshake struct {
	Info_hash [20]byte
	Peer_id   [20]byte
}

type PeerInfo struct {
	addr     net.Addr
	peer_id  []byte
	bitfield []byte

	connStatus bool
	reader     *bufio.Reader
	conn       net.Conn

	amChoking      bool
	amInterested   bool
	peerChoking    bool
	peerInterested bool

	// Mutex sync.Mutex
}

type Message struct {
	id      byte
	payload []byte
}

type ClientInfo struct {
	expectedBitfieldLen int
	Bitfield            []byte
	infoHash            []byte
	pieceHashes         [][20]byte
	piecesLength        uint

	pieceChan chan Block
}

const (
	BlockSize = 16 * 1024
)

type Block struct {
	index uint
	begin uint
	block []byte
}

type Piece struct {
	index    uint
	data     []byte
	received []bool
	done     bool
}

type ActivePieces struct {
	pieces map[int]*Piece
}

type Data struct {
	mu     sync.Mutex
	pieces map[uint]*Piece
}

var data Data
