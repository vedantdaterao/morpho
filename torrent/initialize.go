package torrent

import (
	"fmt"
)

func Initialize(parsed interface{}) (*torrentInfo, error) {
	root, ok := parsed.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("top-level bencode is not a dictionary")
	}

	torrent := &torrentInfo{}

	// Extract announce
	if announce, ok := root["announce"].(string); ok {
		torrent.Announce = announce
	}

	// Extract info dictionary
	infoMap, ok := root["info"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'info' dictionary")
	}

	info := bencodeInfo{}
	if pieces, ok := infoMap["pieces"].(string); ok {
		info.Pieces = pieces
	}
	if pieceLen, ok := infoMap["piece length"].(int); ok {
		info.PieceLength = pieceLen
	}
	if length, ok := infoMap["length"].(int); ok {
		info.Length = length
	}
	if name, ok := infoMap["name"].(string); ok {
		info.Name = name
	}

	torrent.Info = info
	return torrent, nil
}
