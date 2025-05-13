package torrent

type bencodeInfo struct {
	Pieces      string `torrent:"pieces"`
	PieceLength int    `torrent:"piece length"`
	Length      int    `torrent:"length"`
	Name        string `torrent:"name"`
}

type torrentInfo struct {
	Announce string      `torrent:"announce"`
	Info     bencodeInfo `torrent:"info"`
}

type TorrentFile struct {
    Announce    string
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength int
    Length      int
    Name        string
}

type bencodeTrackerResp struct {
	Interval int    `torrent:"interval"`
	Peers    string `torrent:"peers"`
}
