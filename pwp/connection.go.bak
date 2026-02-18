package pwp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func (p *PeerInfo) Listner(client ClientInfo) error {
	for {
		msg, err := read_message(p.reader) 
		if err != nil {
			if err == io.EOF {
                fmt.Println("Peer disconnected:", p.addr)
            } else {
                fmt.Println("Error reading from peer:", err)
            }
            p.conn.Close()
            p.connStatus = false
            break
		}

		p.handleMsg(msg, client)
	}
	return nil
}


func SavePieces(pieceChan chan Block, client *ClientInfo) {
	for b := range pieceChan {
		data.mu.Lock()

		piece, ok := data.pieces[b.index]
		switch {
		case !ok:
			piece = newPiece(b.index, int(client.piecesLength))
			data.pieces[b.index] = piece
		case ok:
			piece = piece.saveBlockData(b)
			data.pieces[b.index] = piece
		}

		data.mu.Unlock()
	}
}


func DialPeer(ip net.IP, port uint16, client ClientInfo, peerID []byte) (*PeerInfo, error){
	peerI := PeerInfo{}
	 
	address := net.JoinHostPort(ip.String(), strconv.Itoa(int(port)))
	
	conn, err := net.DialTimeout("tcp", address, 5 * time.Second) 
	if err != nil {
		return nil, errors.New("connection error: cant connect to the peer")
	}
	peerI.addr = conn.RemoteAddr()
	reader := bufio.NewReader(conn)
	
	// Handshake Request
	handshake := SerializeHandshake(client.infoHash, peerID)
	if _, err := conn.Write(handshake); err != nil {
		conn.Close()
		return nil, err
	}

	// Handshake Response
	// if err := read_handshake(reader, client.infoHash, &peerI); err != nil {
	// 	conn.Close()
	// 	return nil, err
	// } 

	peerI = PeerInfo{
		connStatus: true,
		conn: conn,
		reader: reader,

		amChoking: true, 		// true
		amInterested: false, 	// true
	}
	return &peerI, nil
}

func InitializeClient(infoHash []byte, pieceHashes [][20]byte, pieceLength uint) *ClientInfo{
	return &ClientInfo{
		expectedBitfieldLen: (len(pieceHashes) + 7 )/ 8,
		infoHash: infoHash,
		pieceHashes: pieceHashes,
		piecesLength: pieceLength,
		Bitfield: make([]byte, (len(pieceHashes) + 7) / 8),
		pieceChan: make(chan Block),
	}
}