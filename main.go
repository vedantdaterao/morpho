package main

import (
	"fmt"
	"log"
	"os"

	"github.com/morpho/bencoding"
	"github.com/morpho/torrent"
)

func main() {
    // data, _ := os.ReadFile("assassins_creed_mirage.torrent")
    // data, _ := os.ReadFile("ubuntu-25.04-desktop-amd64.iso.torrent")
    data, _ := os.ReadFile("debian.torrent")
    result, err := bencoding.Decode(data)
    if err != nil {
	    log.Fatal(err)
    }
    
    TorrentFile, err := torrent.Initialize(result)
    if err != nil {
	    log.Fatal(err)
    }

    torrent.Announce(&TorrentFile)
    fmt.Println("all peers -------------------", torrent.AllPeerList.GetPeers())

    // JSON -
    // cleanResult := removeRawKey(result)
    // jsonData, err := json.MarshalIndent(cleanResult, "", "  ")
    // if err != nil {
	//     fmt.Println("Error marshaling to JSON:", err)
	//     return
    // }

}

func removeRawKey(data any) any {
	switch val := data.(type) {
	case map[string]any:
		newMap := make(map[string]any)
		for k, v := range val {
			if k == "raw" || k == "pieces"{
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
