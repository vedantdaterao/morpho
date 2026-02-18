package pwp

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/torgo/torrent"
)

func pieceLengthFor(tf torrent.TorrentFile, index int) int {
	standard := int(tf.PieceLength)
	totalPieces := len(tf.PieceHashes)

	if index < totalPieces-1 {
		return standard
	}
	remainder := tf.Length - (standard * (totalPieces - 1))
	if remainder <= 0 {
		return standard
	}
	return remainder
}

func downloadPiece(peer torrent.PeerInfo, tf torrent.TorrentFile, work pieceWork) ([]byte, error) {
	addr := fmt.Sprintf("%s:%d", peer.IP, peer.Port)

	peerID := []byte("-GO0001-123456789012")
	conn, hs, err := Dial(addr, tf.InfoHash[:], peerID)
	if err != nil {
		return nil, err
	}

	if hs.Info_hash != tf.InfoHash {
		conn.Close()
		return nil, fmt.Errorf("info hash mismatch")
	}

	p := NewPeer(hs.Peer_id, conn)
	go p.Run()
	defer p.Close()

	p.Send(MsgInterested{})

	for msg := range p.Incoming() {
		p.UpdateState(msg)
		if _, ok := msg.(MsgUnchoke); ok {
			break
		}
	}

	if p.PeerChoking {
		return nil, fmt.Errorf("peer still choking")
	}

	if p.Bitfield == nil {
		return nil, fmt.Errorf("peer bitfield nil")
	}
	if !p.Bitfield.HasPiece(work.index) {
		return nil, fmt.Errorf("peer missing piece")
	}

	pieceLen := pieceLengthFor(tf, work.index)
	buf := make([]byte, pieceLen)

	blockSize := 16 * 1024
	numBlocks := (pieceLen + blockSize - 1) / blockSize

	type blockState struct {
		received bool
	}

	blocks := make([]blockState, numBlocks)
	receivedBlocks := 0
	requested := 0
	backlog := 0
	maxBacklog := 5

	for receivedBlocks < numBlocks {
		for !p.PeerChoking &&
			backlog < maxBacklog &&
			requested < pieceLen {

			blockLen := blockSize
			if pieceLen-requested < blockLen {
				blockLen = pieceLen - requested
			}

			p.Send(MsgRequest{
				Index:  uint32(work.index),
				Begin:  uint32(requested),
				Length: uint32(blockLen),
			})

			requested += blockLen
			backlog++
		}

		msg, ok := <-p.Incoming()
		if !ok {
			return nil, fmt.Errorf("peer disconnected")
		}

		p.UpdateState(msg)

		switch m := msg.(type) {
		case MsgPiece:
			if int(m.Index) != work.index {
				continue
			}
			blockIndex := int(m.Begin) / blockSize
			if blockIndex < 0 || blockIndex >= numBlocks {
				continue
			}
			if blocks[blockIndex].received {
				continue
			}
			copy(buf[m.Begin:], m.Block)
			blocks[blockIndex].received = true
			receivedBlocks++
			backlog--

		case MsgChoke:
			return nil, fmt.Errorf("peer re-choked")
		}
	}

	// fmt.Println("piece data -")
	// fmt.Println(hex.Dump(buf))

	sum := sha1.Sum(buf)
	if sum != work.hash {
		return nil, fmt.Errorf("piece hash mismatch")
	}

	return buf, nil
}

func worker(
	peer torrent.PeerInfo, 
	tf torrent.TorrentFile, 
	queue chan pieceWork, 
	results chan<- struct {
		index int
		data  []byte
	},
	) {
	const maxConsecutiveFails = 3
	consecutiveFails := 0

	for work := range queue {
		data, err := downloadPiece(peer, tf, work)
		if err != nil {
			log.Printf("peer %s failed piece %d: %v", peer.IP, work.index, err)
			queue <- work
			consecutiveFails++
			if consecutiveFails >= maxConsecutiveFails {
				log.Printf("peer %s dropped after %d consecutive failures", peer.IP, consecutiveFails)
				return
			}
			continue
		}
		consecutiveFails = 0
		results <- struct {
			index int
			data  []byte
		}{work.index, data}
	}
}

func Download(tf torrent.TorrentFile, peers []torrent.PeerInfo) error {
	numPieces := len(tf.PieceHashes)

	queue := make(chan pieceWork, numPieces)
	results := make(chan struct {
		index int
		data  []byte
	}, numPieces)

	for i, hash := range tf.PieceHashes {
		queue <- pieceWork{i, hash}
	}

	workerCount := 50
	if len(peers) < workerCount {
		workerCount = len(peers)
	}
	for i := 0; i < workerCount; i++ {
		go worker(peers[i], tf, queue, results)
	}

	outFile, err := os.Create(tf.Name)
	if err != nil {
		return err
	}
	if err := outFile.Truncate(int64(tf.Length)); err != nil {
		outFile.Close()
		return err
	}

	var downloadedPieces int32
	for int(downloadedPieces) < numPieces {
		r := <-results
		offset := int64(r.index) * int64(tf.PieceLength)
		if _, err := outFile.WriteAt(r.data, offset); err != nil {
			outFile.Close()
			return fmt.Errorf("write piece %d: %w", r.index, err)
		}
		atomic.AddInt32(&downloadedPieces, 1)
		log.Printf("piece %d/%d done", int(downloadedPieces), numPieces)
	}

	close(queue)
	fmt.Println("\nDownload complete.")
	return outFile.Close()
}
