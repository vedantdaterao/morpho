package pwp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

func read_handshake(r *bufio.Reader, infoHash []byte, peer *Peer) error {
	buf := make([]byte, HandshakeLen)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return ErrMessageLen
	}

	if !bytes.Equal(buf[1+len(ProtocolStr)+8:48], infoHash) {
		return ErrInfoHashMismatch
	}

	peer.peer_id = buf[1+len(ProtocolStr)+8+20:]
	fmt.Print("\npeer handshake -\n", hex.Dump(buf))
	return nil
}

func read_message(r *bufio.Reader) (*Message, error) {
	// <length prefix><message ID><payload>
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lenBuf)
	if err != nil {
		return nil, err
	}
	len := binary.BigEndian.Uint32(lenBuf)

	msgBuf := make([]byte, len)
	_, err = io.ReadFull(r, msgBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		id:      uint8(msgBuf[0]),
		payload: msgBuf[1:],
	}

	// log.Printf("\npeer message -\nlen - %s\nmsg -", len, m)
	return &m, nil
}

// func (m *Message) handle_msg_bitfield(expectedLen int) ([]byte) {
// 	if len(m.payload) != expectedLen {
// 		return nil
// 	}
// 	return m.payload[:]
// }

func (m *Message) handle_msg_piece() (*Block, error) {
	if len(m.payload) < 8 {
		return nil, ErrMessageLen
	}
	index := binary.BigEndian.Uint32(m.payload[0:4])
	begin := binary.BigEndian.Uint32(m.payload[4:8])
	block := m.payload[8:]

	return &Block{
		index: uint(index),
		begin: uint(begin),
		block: block,
	}, nil
}

func newPiece(index uint, totalLength int) *Piece {
	numBlocks := (totalLength + BlockSize - 1) / BlockSize
	return &Piece{
		index:    index,
		data:     make([]byte, totalLength),
		received: make([]bool, numBlocks),
		done:     false,
	}
}

func (p *Piece) saveBlockData(b Block) *Piece {
	if p.index == b.index {
		copy(p.data[b.begin:], b.block[:])
		p.received[b.index] = true
	}
	for _, v := range p.received {
		if !v {
			return p
		}
	}
	p.done = true
	return p
}

func (p *PeerInfo) handleMsg(m *Message, client ClientInfo) error {
	switch m.id {

	case Msg_Choke:
		p.peerChoking = true
	case Msg_Unchoke:
		p.peerChoking = false
	case Msg_Interested:
		p.peerInterested = true
	case Msg_NotInterested:
		p.peerInterested = false

	case Msg_Have:
		index := binary.BigEndian.Uint32(m.payload[:4])
		SetPiece(p.bitfield, int(index))
	case Msg_Bitfield:
		p.bitfield = m.payload[:]
	case Msg_Piece:
		pMsg, _ := m.handle_msg_piece()
		// send this block (pMsg) to the piece downloader func
		client.pieceChan <- *pMsg
		// download_piece(pMsg.index, pMsg, client)
		// TODO: download piece

		// fmt.Printf("Unknown message ID: %d\n", m.id)
	}
	return nil
}

func (p *Peer) handleMsg(m *Message) (interface{}, error) {
	switch m.id {

	case Msg_Choke:
		p.peerChoking = true
	case Msg_Unchoke:
		p.peerChoking = false
	case Msg_Interested:
		p.peerInterested = true
	case Msg_NotInterested:
		p.peerInterested = false

	case Msg_Have:
		index := binary.BigEndian.Uint32(m.payload[:4])
		SetPiece(p.bitfield, int(index))
	case Msg_Bitfield:
		p.bitfield = m.payload[:]
	case Msg_Piece:
		pMsg, _ := m.handle_msg_piece()
		// send this block (pMsg) to the piece downloader func
		// download_piece(pMsg.index, pMsg, client)
		// TODO: download piece
		// if p.ap.pieces[int(pMsg.index)] == nil {
		// 	p.ap.pieces[int(pMsg.index)] = // Piece
		// } else {

		// }
		return pMsg, nil

		// fmt.Printf("Unknown message ID: %d\n", m.id)
	}
	return nil, nil
}

// func sad(a []byte, b []byte) {
// 	copy(a, b)
// }

// func download_piece(i uint, b *Block, c ClientInfo) {
// 	data.mu.Lock()
// 	defer data.mu.Unlock()
// 	// Initialize the piece if it's not already in memory

// 	totalLength := c.piecesLength
// 	if i == len(c.pieceHashes) - 1 {
// 		totalLength = c.totalLength - (c.pieceLength * (len(c.pieceHashes) - 1))
// 	}
// 	blockIndex := b.begin / BlockSize

// 	if data.pieces[i].data == nil {
// 		data.pieces[i] = &Piece{
// 			index: i,
// 			data: make([]byte, c.totalLength),
// 			received: make([]bool, ),
// 			// data: copy(data.pieces[i].data[b.begin:], b.block),
// 			// received[blockIndex]: true,
// 		}
// 	}

// 	piece := &data.pieces[index]
// 	begin := b.begin
// 	copy(piece.data[begin:], b.block)

// 	blockIndex := begin / BlockSize
// 	piece.received[blockIndex] = true

// 	// Check if all blocks are received
// 	if isPieceComplete(piece.received) {
// 		piece.done = true
// 		go savePieceToDisk(index, piece.data)

// 		// Free memory
// 		data.pieces[index] = Piece{}
// 	}
// }

func (p *PeerInfo) HasPiece(i int) bool {
	byteIndex := i / 8
	offset := i % 8
	return p.bitfield[byteIndex]>>(7-offset)&1 != 0
}

func SetPiece(bf []byte, i int) {
	byteIndex := i / 8
	offset := i % 8
	bf[byteIndex] |= 1 << (7 - offset)
}
