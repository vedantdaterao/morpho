package pwp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func Connect(ip net.IP, port uint16, infoHash []byte, peerID []byte) (*Peer, error) {
	peer := Peer{}

	address := net.JoinHostPort(ip.String(), strconv.Itoa(int(port)))

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, &ConnectionErr{msg: "can't connect to peer"}
	}

	peer.addr = conn.RemoteAddr()
	fmt.Println(peer.addr)
	reader := bufio.NewReader(conn)

	// Handshake Request
	handshake := SerializeHandshake(infoHash, peerID)
	if _, err := conn.Write(handshake); err != nil {
		conn.Close()
		return nil, err
	}

	// Handshake Response
	if err := read_handshake(reader, infoHash, &peer); err != nil {
		conn.Close()
		return nil, err
	}

	peer.connStatus = true
	peer.conn = conn
	peer.reader = reader

	peer.amChoking = true     // true
	peer.amInterested = false // true

	return &peer, nil
}

func (p *Peer) ReadMesage() (*Message, error) {
	m, err := read_message(p.reader)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *Peer) WriteMessage() {

}

func RequestPiece(index uint, pieceLen uint, pieceHash [20]byte) {

}

func (p *Peer) Listner(msgChan chan<- interface{}) error {
	for {
		msg, err := read_message(p.reader)
		if err != nil {
			if err == io.EOF {
				p.conn.Close()
				p.connStatus = false
				return &ConnectionErr{msg: "peer disconnrcted"}
				// fmt.Println("Peer disconnected:", p.addr)
			} else {
				p.conn.Close()
				p.connStatus = false
				return &ConnectionErr{msg: "error reading from peer"}
				// fmt.Println("Error reading from peer:", err)
			}
			// break
		}
		// block, _ := p.handleMsg(msg)
		fmt.Println("msg - ", msg)
		msgChan <- msg
	}
}
