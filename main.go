package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/torgo/bencoding"
	"github.com/torgo/pwp"
	"github.com/torgo/torrent"
)


func formatBytes(bytes int) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	torrentFile := flag.String("f", "", "path to .torrent file")
	verbose := flag.Bool("v", false, "enable verbose logging")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -f <file.torrent> [-v]\n", os.Args[0])
	}
	flag.Parse()

	if *torrentFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Suppress logs unless -v is set
	if !*verbose {
		log.SetOutput(io.Discard)
	}

	data, err := os.ReadFile(*torrentFile)
	if err != nil {
		log.Fatal(err)
	}

	meta, err := bencoding.Decode(data)
	if err != nil {
		log.Fatal(err)
	}

	tf, err := torrent.Initialize(meta)
	if err != nil {
		log.Fatal(err)
	}

	// Display torrent info
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Name:         %s\n", tf.Name)
	fmt.Printf("Size:         %s (%d bytes)\n", formatBytes(tf.Length), tf.Length)
	fmt.Printf("Pieces:       %d × %s\n", len(tf.PieceHashes), formatBytes(int(tf.PieceLength)))
	fmt.Printf("Tracker:      %s\n", tf.Announce)
	if len(tf.Files) > 0 {
		fmt.Printf("Files:        %d\n", len(tf.Files))
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	torrent.Announce(&tf)

	peers := torrent.AllPeerList.GetPeers()
	if len(peers) == 0 {
		log.Fatal("no peers received from tracker")
	}

	fmt.Printf("Peers:        %d\n", len(peers))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	if err := pwp.Download(tf, peers); err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("✓ Download complete.")
}



// func main() {
// 	// data, _ := os.ReadFile("torrent_files/CoD-BlackOps6.torrent")
// 	// data, _ := os.ReadFile("torrent_files/debian.iso.torrent")
// 	data, _ := os.ReadFile("torrent_files/debian.iso.torrent")
// 	result, err := bencoding.Decode(data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	torrent.torrent.torrent.TorrentFile, err := torrent.Initialize(result)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	torrent.Announce(&torrent.torrent.torrent.TorrentFile)
// 	peerList := torrent.AllPeerList.GetPeers()
// 	fmt.Println("all peers -------------------", peerList)

// 	// peer connection
// }
