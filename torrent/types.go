package torrent

import (
	"net"
	"sync"
)

type TorrentMeta struct {
	Announce		string			`torrent:"announce"`
	AnnounceList	[]string		`torrent:"announce-list,omitempty"`
	Info			*InfoDict		`torrent:"info"`
}

type InfoDict struct {
	PieceLength		uint			`torrent:"piece length"`
	Pieces			[][20]byte		`torrent:"pieces"` 				// concatenated 20-byte SHA1 hashes
	Private     	*int			`torrent:"private,omitempty"`
	Name 			string 			`torrent:"name"`				// Shared between single and multi-file modes
	Length 			int  			`torrent:"length,omitempty"`	// Single-file mode
	Files 			[]FileEntry 	`torrent:"files,omitempty"`		// Multi-file mode
}

type FileEntry struct {
	Length 			int				`torrent:"length"`
	Path   			[]string		`torrent:"path"`
}

var Info InfoDict

var Torrent TorrentMeta

type TorrentFile struct {
    Announce    string
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength uint
    Length      int
	Name        string
	Files		[]FileEntry
}


type TrackerResponse struct {
	FailureReason string        `torrent:"failure reason,omitempty"` // Only present if failure
	WarningMessage string       `torrent:"warning message,omitempty"`
	Interval       int          `torrent:"interval"`                 // Required
	MinInterval    int          `torrent:"min interval,omitempty"`
	TrackerID      string       `torrent:"tracker id,omitempty"`
	Complete       int          `torrent:"complete,omitempty"`       // Seeders
	Incomplete     int          `torrent:"incomplete,omitempty"`     // Leechers
}

type PeerInfo struct {
	PeerID string 				`torrent:"peer id"` // 20-byte peer ID
	IP     net.IP 				`torrent:"ip"`      // IPv4, IPv6 or DNS
	Port   uint16 				`torrent:"port"`    // TCP port number
}

type PeerList struct {
	mu 		sync.Mutex
	Peers 	[]PeerInfo
}

var AllPeerList PeerList