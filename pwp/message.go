package pwp

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

// ref: https://wiki.theory.org/BitTorrentSpecification#Messages

// keep-alive:
// 		<len=0000>
type MsgKeepAlive struct{}

func (MsgKeepAlive) ID() byte        { return 255 }
func (MsgKeepAlive) Payload() []byte { return nil }

// choke: 
// 		<len=0001><id=0>
type MsgChoke struct{}

func (MsgChoke) ID() byte        { return Msg_Choke }
func (MsgChoke) Payload() []byte { return nil }

// unchoke: 
// 		<len=0001><id=1>
type MsgUnchoke struct{}

func (MsgUnchoke) ID() byte        { return Msg_Unchoke }
func (MsgUnchoke) Payload() []byte { return nil }

// interested: 
// 		<len=0001><id=2>
type MsgInterested struct {}

func (MsgInterested) ID() byte 			{ return Msg_Interested }
func (MsgInterested) Payload() []byte 	{ return nil }

// not interested: 
// 		<len=0001><id=3>
type MsgNotInterested struct{}

func (MsgNotInterested) ID() byte			{ return Msg_NotInterested }
func (MsgNotInterested) Payload() []byte 	{ return nil }

// have: 
// 		<len=0005><id=4><piece index>
type MsgHave struct {
	Index uint32
}

func (MsgHave) ID() byte { return Msg_Have }
func (m MsgHave) Payload() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, m.Index)
	return b
}

// bitfield: 
// 		<len=0001+X><id=5><bitfield>
type MsgBitfield struct {
	Bitfield []byte
}

func (MsgBitfield) ID() byte { return Msg_Bitfield}
func (m MsgBitfield) Payload() []byte {
	return m.Bitfield
}

// request: 
// 		<len=0013><id=6><index><begin><length>
type MsgRequest struct {
	Index  uint32
	Begin  uint32
	Length uint32
}

func (MsgRequest) ID() byte { return Msg_Request }
func (m MsgRequest) Payload() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf[0:], m.Index)
	binary.BigEndian.PutUint32(buf[4:], m.Begin)
	binary.BigEndian.PutUint32(buf[8:], m.Length)
	return buf
}

// piece: 
// 		<len=0009+X><id=7><index><begin><block>
// The piece message is variable length, where X is the length of the block.
// The payload contains the following information:
// index: integer specifying the zero-based piece index
// begin: integer specifying the zero-based byte offset within the piece
// block: block of data, which is a subset of the piece specified by index.
type MsgPiece struct {
	Index uint32
	Begin uint32
	Block []byte
}

func (MsgPiece) ID() byte { return Msg_Piece}
func (m MsgPiece) Payload() []byte {
	buf := make([]byte, 8 + len(m.Block))
	binary.BigEndian.PutUint32(buf[0:], m.Index)
	binary.BigEndian.PutUint32(buf[4:], m.Begin)
	copy(buf[8:], m.Block)
	return buf
}

// cancel: 
// 		<len=0013><id=8><index><begin><length>
// The cancel message is fixed length, and is used to cancel block requests. 
// The payload is identical to that of the "request" message. It is typically used during "End Game".
type MsgCancel struct {
	Index  uint32
	Begin  uint32
	Length uint32
}

func (MsgCancel) ID() byte { return Msg_Cancel }
func (m MsgCancel) Payload() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf[0:], m.Index)
	binary.BigEndian.PutUint32(buf[4:], m.Begin)
	binary.BigEndian.PutUint32(buf[8:], m.Length)
	return buf
}

// --------------
// Read message from stream
func ReadMessage(r *bufio.Reader) (Message, error) {
	// read length prefix (payload + id byte)
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lenBuf)
	if length == 0 {
		// return nil, fmt.Errorf("keep-alive not handled here")
		return MsgKeepAlive{}, nil
	}

	// length includes the message id
	msgBuf := make([]byte, length)
	if _, err := io.ReadFull(r, msgBuf); err != nil {
		return nil, err
	}

	id := msgBuf[0]
	payload := msgBuf[1:]

	// decode message
	switch id {
	case Msg_Choke:
		return MsgChoke{}, nil
	case Msg_Unchoke:
		return MsgUnchoke{}, nil
	case Msg_Interested:
		return MsgInterested{}, nil
	case Msg_NotInterested:
		return MsgNotInterested{}, nil
	case Msg_Have:
		return MsgHave{Index: binary.BigEndian.Uint32(payload)}, nil
	case Msg_Bitfield:
		return MsgBitfield{Bitfield: payload}, nil
	case Msg_Request:
		return MsgRequest{
			Index:  binary.BigEndian.Uint32(payload[0:4]),
			Begin:  binary.BigEndian.Uint32(payload[4:8]),
			Length: binary.BigEndian.Uint32(payload[8:12]),
		}, nil
	case Msg_Piece:
		return MsgPiece{
			Index: binary.BigEndian.Uint32(payload[0:4]),
			Begin: binary.BigEndian.Uint32(payload[4:8]),
			Block: payload[8:],
		}, nil
	case Msg_Cancel:
		return MsgCancel{
			Index:  binary.BigEndian.Uint32(payload[0:4]),
			Begin:  binary.BigEndian.Uint32(payload[4:8]),
			Length: binary.BigEndian.Uint32(payload[8:12]),
		}, nil
	default:
		return nil, fmt.Errorf("unknown message id %d", id)
	}
}

func WriteMessage(w *bufio.Writer, m Message) error {
	// <length prefix><message ID><payload>
	payload := m.Payload()
	length := uint32(len(payload) + 1)

	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = m.ID()
	copy(buf[5:], payload)
	
	_, err := w.Write(buf)
	return err
}

func SerializeMessage(m Message) []byte {
	payload := m.Payload()
	length := uint32(len(payload) + 1)

	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = m.ID()
	copy(buf[5:], payload)

	return buf
}

func SerializeHandshake(infoHash []byte, peerID []byte) []byte {
	buf := make([]byte, HandshakeLen)
	
	buf[0] = byte(ProtocolStrLen)
	copy(buf[1:], ProtocolStr)
	copy(buf[1+ProtocolStrLen:], make([]byte, 8)) // reserved
	copy(buf[1+ProtocolStrLen+8:], infoHash[:20])
	copy(buf[1+ProtocolStrLen+8+20:], peerID[:20])
	
	// fmt.Print("Handshake Data -\n", hex.Dump(buf))
	
	return buf
}

// example main 
// func main() {
// 	// simulate input stream — writing a MsgHave message into a buffer
// 	buf := make([]byte, 0)
// 	writer := appendLengthPrefixedMessage(&buf, MsgHave{Index: 123})

// 	reader := bufio.NewReader(writer)
// 	msg, err := readMessage(reader)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("Received message ID=%d", msg.ID())
// 	switch m := msg.(type) {
// 	case MsgHave:
// 		log.Printf("Have: %d", m.Index)
// 	default:
// 		log.Printf("Unhandled message type")
// 	}
// }

// // Helper: encode and append a message
// func appendLengthPrefixedMessage(buf *[]byte, m Message) *os.File {
// 	payload := m.Payload()
// 	msg := append([]byte{m.ID()}, payload...)
// 	length := uint32(len(msg))

// 	lb := make([]byte, 4)
// 	binary.BigEndian.PutUint32(lb, length)

// 	*buf = append(*buf, lb...)
// 	*buf = append(*buf, msg...)

// 	tmpfile, _ := os.CreateTemp("", "msgbuf")
// 	tmpfile.Write(*buf)
// 	tmpfile.Seek(0, 0)
// 	return tmpfile
// }
