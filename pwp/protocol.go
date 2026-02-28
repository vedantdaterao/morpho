package pwp

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func NewBitfield(numPieces int) *Bitfield {
	return &Bitfield{Bits: make([]byte, (numPieces+7)/8)}
}

func BitfieldFromBytes(b []byte) *Bitfield {
	return &Bitfield{Bits: b}
}

func (bf *Bitfield) HasPiece(i int) bool {
	byteIndex := i / 8
	offset := i % 8
	return bf.Bits[byteIndex]>>(7-offset)&1 != 0
}

func (bf *Bitfield) SetPiece(i int) {
	byteIndex := i / 8
	offset := i % 8
	bf.Bits[byteIndex] |= 1 << (7 - offset)
}

func (bf *Bitfield) Bytes() []byte {
	return bf.Bits
}

// ----------------------------------------------------------------------------
type peerConn struct {
	conn net.Conn
	r    *bufio.Reader
	w    *bufio.Writer
}

func (pc *peerConn) WriteMessage(msg Message) error {
	// keep-alive: length=0, no id byte
	if _, ok := msg.(MsgKeepAlive); ok {
		_, err := pc.w.Write([]byte{0, 0, 0, 0})
		if err != nil {
			return err
		}
		return pc.w.Flush()
	}
	if err := WriteMessage(pc.w, msg); err != nil {
		return err
	}
	return pc.w.Flush()
}

func (pc *peerConn) ReadMessage() (Message, error) {
	return ReadMessage(pc.r)
}

func (pc *peerConn) Close() error {
	return pc.conn.Close()
}

func (pc *peerConn) RemoteAddr() string {
	return pc.conn.RemoteAddr().String()
}

// ----------------------------------------------------------------------------
// Handshake

// BitTorrent handshake
// - PeerConn ready for message exchange
func Dial(addr string, infoHash, peerID []byte) (PeerConn, *Msg_handshake, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("dial %s: %w", addr, err)
	}

	pc := &peerConn{
		conn: conn,
		r:    bufio.NewReader(conn),
		w:    bufio.NewWriter(conn),
	}

	// send handshake
	if _, err := conn.Write(SerializeHandshake(infoHash, peerID)); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("send handshake: %w", err)
	}

	// peer handshake
	hs, err := readHandshake(pc.r)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("read handshake: %w", err)
	}

	return pc, hs, nil
}

func readHandshake(r *bufio.Reader) (*Msg_handshake, error) {
	buf := make([]byte, HandshakeLen)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	pstrlen := int(buf[0])
	if pstrlen != ProtocolStrLen {
		return nil, fmt.Errorf("unexpected pstrlen %d", pstrlen)
	}
	if string(buf[1:1+pstrlen]) != ProtocolStr {
		return nil, fmt.Errorf("unexpected protocol string")
	}

	hs := &Msg_handshake{}
	copy(hs.Info_hash[:], buf[1+pstrlen+8:1+pstrlen+8+20])
	copy(hs.Peer_id[:], buf[1+pstrlen+8+20:])
	return hs, nil
}

// ----------------------------------------------------------------------------
// Peer - goroutine message loop

func NewPeer(id [20]byte, conn PeerConn) *Peer {
	return &Peer{
		ID:          id,
		Conn:        conn,
		AmChoking:   true,
		PeerChoking: true,
		incoming:    make(chan Message, 32),
		outgoing:    make(chan Message, 32),
		close:       make(chan struct{}),
	}
}

// start the send/receive goroutines 
// Blocks until the peer disconnects
// or Close() is called
func (p *Peer) Run() {
	go p.readLoop()
	p.writeLoop()
}

func (p *Peer) readLoop() {
	for {
		msg, err := p.Conn.ReadMessage()
		if err != nil {
			close(p.incoming)
			return
		}
		select {
		case p.incoming <- msg:
		case <-p.close:
			return
		}
	}
}

func (p *Peer) writeLoop() {
	for {
		select {
		case msg, ok := <-p.outgoing:
			if !ok {
				return
			}
			p.Conn.WriteMessage(msg) //nolint: errcheck - closed conn on next read
		case <-p.close:
			return
		}
	}
}

// message for sending to THIS peer
func (p *Peer) Send(msg Message) {
	select {
	case p.outgoing <- msg:
	case <-p.close:
	}
}

// channel of messages received from THIS peer
func (p *Peer) Incoming() <-chan Message {
	return p.incoming
}

// shuts down the peer's goroutines and connection
func (p *Peer) Close() {
	select {
	case <-p.close:
	default:
		close(p.close)
	}
	p.Conn.Close()
}

func (p *Peer) UpdateState(msg Message) {
	switch m := msg.(type) {
	case MsgChoke:
		p.PeerChoking = true
	case MsgUnchoke:
		p.PeerChoking = false
	case MsgInterested:
		p.PeerInterested = true
	case MsgNotInterested:
		p.PeerInterested = false
	case MsgBitfield:
		p.Bitfield = BitfieldFromBytes(m.Bitfield)
	case MsgHave:
		if p.Bitfield != nil {
			p.Bitfield.SetPiece(int(m.Index))
		}
	}
}

// ----------------------------------------------------------------------------
// Listen - accept incoming peer connections

type ConnHandler func(conn PeerConn, hs *Msg_handshake)

// Use "127.0.0.1:0" to let the OS pick a free port
func Listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

// accepts connections on ln and calls handler for each one
func Serve(ln net.Listener, infoHash, peerID []byte, handler ConnHandler) error {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn) {
			pc := &peerConn{
				conn: c,
				r:    bufio.NewReader(c),
				w:    bufio.NewWriter(c),
			}
			hs, err := readHandshake(pc.r)
			if err != nil {
				c.Close()
				return
			}
			if _, err := c.Write(SerializeHandshake(infoHash, peerID)); err != nil {
				c.Close()
				return
			}
			handler(pc, hs)
		}(conn)
	}
}

// listens on addr and calls handler for each incoming connection
func ListenAndServe(addr string, infoHash, peerID []byte, handler ConnHandler) error {
	ln, err := Listen(addr)
	if err != nil {
		return err
	}
	return Serve(ln, infoHash, peerID, handler)
}
