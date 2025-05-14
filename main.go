package main

import (
	"fmt"
	"log"
	"os"

	"github.com/morpho/bencoding"
	"github.com/morpho/torrent"
)

func main() {
    data, _ := os.ReadFile("ubuntu-25.04-desktop-amd64.iso.torrent")
    result, err := bencoding.Decode(data)
    if err != nil {
	    log.Fatal(err)
    }
    TorrentFile, err := torrent.Initialize(result)
    if err != nil {
	    log.Fatal(err)
    }

    fmt.Println(TorrentFile.BuildTrackerURL())

    // fmt.Println(result)

    // JSON -
    // jsonData, err := json.MarshalIndent(result, "", "  ")
    // if err != nil {
	//     fmt.Println("Error marshaling to JSON:", err)
	//     return
    // }
    // fmt.Println(string(jsonData))

    // fmt.Println(TorrentFile)

}
