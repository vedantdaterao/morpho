package pwp

const (
	ProtocolStr    = "BitTorrent protocol"
	ProtocolStrLen = len(ProtocolStr)      // 19
	HandshakeLen   = 49 + len(ProtocolStr) // total: 68 bytes
)

const (
	Msg_Choke         byte = iota
	Msg_Unchoke       
	Msg_Interested    
	Msg_NotInterested 
	Msg_Have          
	Msg_Bitfield      
	Msg_Request       
	Msg_Piece         
	Msg_Cancel        
)

type Message interface {
	ID() 		byte
	Payload() 	[]byte
}

// handshake: <pstrlen><pstr><reserved><info_hash><peer_id>
type Msg_handshake struct {
	Info_hash [20]byte
	Peer_id   [20]byte
}

type Bitfield struct {
	Bits 	[]byte
}

type PeerConn interface {
    WriteMessage(msg Message) 	error
    ReadMessage	() 				(Message, error)
    Close		() 				error
    RemoteAddr	() 				string
}

type Peer struct {
    ID       [20]byte
    Conn     PeerConn
    Bitfield *Bitfield
    
    // State
    AmChoking      bool
    AmInterested   bool
    PeerChoking    bool
    PeerInterested bool
    
    // Channels for communication with downloader
    incoming chan Message
    outgoing chan Message
    close    chan struct{}     // Signal to stop 
}


type pieceWork struct {
	index int
	hash  [20]byte
}