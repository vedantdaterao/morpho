package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/morpho/bencoding"
	"github.com/morpho/pwp"
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

	// JSON -
	// cleanResult := removeRawKey(result)
	// jsonData, err := json.MarshalIndent(cleanResult, "", "  ")
	// if err != nil {
	//     fmt.Println("Error marshaling to JSON:", err)
	//     return
	// }

	// peer connection
	peer, err := pwp.Connect(peerList[1].IP, peerList[1].Port, TorrentFile.InfoHash[:], []byte("-XXXXXXXXXXXXXXXXXXX"))

	if err != nil {
		fmt.Println(err)
	}

	// jsondata, _ := json.MarshalIndent(peer, "", " ")
	fmt.Println("peer -", peer)

	// msgChan := make(chan interface{})

	// go peer.Listner(msgChan)

	// for msg := range msgChan {
	// 	fmt.Println("Received: -", msg)
	// }

	// fmt.Print("\n\n peer - ", peerAddr(peerList[8].IP, peerList[8].Port), "\n\n")

	// conn, err := net.Dial("tcp", peerAddr(peerList[2].IP, peerList[2].Port))
	// if err != nil{
	// 	log.Fatal("\n\n Failed to Connect to the Peer\n\n")
	// }
	// conn.Write(pwp.SerializeHandshake(TorrentFile.InfoHash[:], []byte("-XXXXXXXXXXXXXXXXXXX")))

	// handleClient(conn)

}

func peerAddr(ip net.IP, port uint16) string {
	return fmt.Sprintf("%s:%s", ip, strconv.FormatUint(uint64(port), 10))
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Create a buffer to read data into
	buffer := make([]byte, 2048)

	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Process and use the data (here, we'll just print it)
		fmt.Print("Received: \n", hex.Dump(buffer[:n]))
	}
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
