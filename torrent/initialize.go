package torrent

import (
	"crypto/sha1"
	"fmt"
)

func Initialize(parsed interface{}) (TorrentFile, error) {

	root, ok := parsed.(map[string]interface{})
	if !ok {
		return TorrentFile{}, fmt.Errorf("top-level bencode is not a dictionary")
	}

	// Extract TorrentInfo
	if announce, ok := root["announce"].(string); ok {
		Torrent.Announce = announce
	}
	if announceList, ok := root["announce-list"].([][]string); ok {
		Torrent.AnnounceList = announceList
	}
	

	// Extract info dictionary
	infoMap, ok := root["info"].(map[string]interface{})
	if !ok {
		return TorrentFile{}, fmt.Errorf("missing or invalid 'info' dictionary")
	}
	if pieceLen, ok := infoMap["piece length"].(uint); ok {
		Info.PieceLength = pieceLen
	}
	if priv, ok := infoMap["private"].(int); ok {
		Info.Private = &priv
	}
	if name, ok := infoMap["name"].(string); ok {
		Info.Name = name
	}

	// single file or multi-file 
	if files, ok := infoMap["files"].([]interface{}); ok {
		var file FileEntry
		for _, f := range files{
			fileDict, ok := f.(map[string]interface{})
			if !ok {
				return TorrentFile{}, fmt.Errorf("invalid file entry in 'files'")
			}
			if length, ok := fileDict["length"].(int); ok {
				file.Length = length
			}
			if pathList, ok := fileDict["path"].([]interface{}); ok {
				for _, p := range pathList {
					if s, ok := p.(string); ok {
						file.Path = append(file.Path, s)
					}
				}
			}
			}
			Info.Files = append(Info.Files, file)
	} else {
		// Single-file mode
		if length, ok := infoMap["length"].(int); ok {
			Info.Length = length
			fmt.Println(Info.Length)
		}
	}
	
	if piecesRaw, ok := infoMap["after_pieces"].([]byte); ok{
		for i := 0; i+20 <= len(piecesRaw); i += 20 {
    		var hash [20]byte
    		copy(hash[:], piecesRaw[i:i+20])
    		Info.Pieces = append(Info.Pieces, hash)
		}
	}
	
	if pieceLen, ok := infoMap["piece length"].(uint); ok {
		Torrent.Info.PieceLength = pieceLen
	}

	Torrent.Info = &Info
	
	var hash [20]byte
	if rawInfo, ok := infoMap["raw"].([]byte); ok{
		hash = sha1.Sum(rawInfo)
	}

	t := TorrentFile{
		Announce: Torrent.Announce,
		InfoHash: hash,
		PieceHashes: Torrent.Info.Pieces,
		PieceLength: Torrent.Info.PieceLength,
		Name: Torrent.Info.Name,
	}
	
	if Torrent.Info.Length > int(0) {
		// Single-file mode
		t.Length = int(Torrent.Info.Length)
	} else if len(Torrent.Info.Files) > 0 {
		// Multi-file mode
		t.Files = Torrent.Info.Files
		// Optional: calculate total length if needed
		total := 0
		for _, file := range Torrent.Info.Files {
			total += int(file.Length)
		}
		t.Length = total
	} else {
		panic("invalid torrent: no length or files")
	}

	return t, nil
}

