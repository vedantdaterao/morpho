package torrent

type TorrentMeta struct {
	Announce		string			`torrent:"announce"`
	AnnounceList	[][]string		`torrent:"announce-list,omitempty"`
	Info			*InfoDict		`torrent:"info"`
}

type InfoDict struct {
	PieceLength		uint			`torrent:"piece length"`
	Pieces			[][20]byte		`torrent:"pieces"` // concatenated 20-byte SHA1 hashes
	Private     	*int			`torrent:"private,omitempty"`

	// Shared between single and multi-file modes
	Name 			string 			`torrent:"name"`

	// Single-file mode
	Length 			int  			`torrent:"length,omitempty"`

	// Multi-file mode
	Files 			[]FileEntry 	`torrent:"files,omitempty"`
}

type FileEntry struct {
	Length 			int			`torrent:"length"`
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
