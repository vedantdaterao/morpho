package main

import (
	"fmt"
	"log"
	"os"

	"github.com/morpho/bencoding"
	"github.com/morpho/torrent"
)

func main() {
	// data, _ := os.ReadFile("ubuntu-25.04-desktop-amd64.iso.torrent")
	// data, _ := os.ReadFile("torrents/assassins_creed_mirage.torrent")
	data, _ := os.ReadFile("torrents/debian.torrent")
	result, err := bencoding.Decode(data)
	if err != nil {
		log.Fatal(err)
	}

	TorrentFile, err := torrent.Initialize(result)
	if err != nil {
		log.Fatal(err)
	}

	torrent.Announce(&TorrentFile)
	peerList := torrent.AllPeerList.GetPeers()
	fmt.Println("all peers -------------------", peerList)

	// peer connection
}

func removeRawKey(data any) any {
	switch val := data.(type) {
	case map[string]any:
		newMap := make(map[string]any)
		for k, v := range val {
			if k == "raw" || k == "pieces" {
				continue
			}
			newMap[k] = removeRawKey(v)
		}
		return newMap
	case []any:
		for i, v := range val {
			val[i] = removeRawKey(v)
		}
		return val
	default:
		return val
	}
}
