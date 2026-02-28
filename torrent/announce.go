package torrent

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/torgo/bencoding"
)

func (a *PeerList) AddPeers(newPeers []PeerInfo) {
	a.mu.Lock()
	defer a.mu.Unlock()
	seen := make(map[string]bool)
	for _, p := range a.Peers {
		seen[p.IP.String()+":"+strconv.Itoa(int(p.Port))] = true
	}

	for _, p := range newPeers {
		key := p.IP.String() + ":" + strconv.Itoa(int(p.Port))
		if !seen[key] {
			a.Peers = append(a.Peers, p)
			seen[key] = true
		}
	}
}

func (a *PeerList) GetPeers() []PeerInfo {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]PeerInfo{}, a.Peers...)
}

func (p *PeerList) GetNumberOfPeers() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.Peers)
}

func Announce(tf *TorrentFile) error {
	url, err := tf.BuildTrackerURL()
	if err != nil {
		return errors.New("tracker url error")
	}
	tries := 3
	// fmt.Println("URL - ", url)
	GetTrackerResponse(url)

	if AllPeerList.GetNumberOfPeers() < 30 {
		for _, x := range Torrent.AnnounceList {
			url, err := tf.BuildTrackerURL(x)
			if err != nil {
				return errors.New("tracker url error")
			}
			// fmt.Println("URL - ", url)
			if _, err := GetTrackerResponse(url); err != nil {
				continue
			}
			// trackerResp, _ := GetTrackerResponse(url)
			// jsonData, err := json.MarshalIndent(trackerResp, "", "  ")
			// if err != nil {
			// 	fmt.Println("Error marshaling to JSON:", err)
			// }
			// fmt.Println(string(jsonData))
		}
	}
	// JSON response
	// jsonData, err := json.MarshalIndent(trackerResp, "", "  ")
	// if err != nil {
	//     fmt.Println("Error marshaling to JSON:", err)
	// }
	// fmt.Println(string(jsonData))
	if AllPeerList.GetNumberOfPeers() == 0 {
		if tries == 0{
			return errors.New("can't get peers from trackers")
		}
		GetTrackerResponse(url)
		tries--
	}
	return nil
}

func GetTrackerResponse(fullURL string) (TrackerResponse, error) {
	if strings.HasPrefix(fullURL, "udp://") {
		return TrackerResponse{}, nil
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(fullURL)
	if err != nil {
		return TrackerResponse{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return TrackerResponse{}, fmt.Errorf("failed to read tracker response: %w", err)
	}

	decoded, err := bencoding.Decode(bodyBytes)
	if err != nil {
		return TrackerResponse{}, fmt.Errorf("bencode decode failed: %w", err)
	}

	trackerResp, err := parseTrackerResponse(decoded)
	if err != nil {
		return TrackerResponse{}, fmt.Errorf("response parser failed: %w", err)
	}

	return *trackerResp, nil
}

func parseTrackerResponse(data any) (*TrackerResponse, error) {
	dict, ok := data.(map[string]any)
	if !ok {
		return nil, errors.New("not bencoded dict")
	}

	resp := &TrackerResponse{}

	for key, val := range dict {
		switch key {
		case "falure reason":
			resp.FailureReason = val.(string)
			return resp, nil
		case "warning message":
			resp.WarningMessage = val.(string)
		case "tracker id":
			resp.TrackerID = val.(string)
		case "interval":
			resp.Interval = val.(int)
		case "min interval":
			resp.MinInterval = val.(int)
		case "complete":
			resp.Complete = val.(int)
		case "incomplete":
			resp.Incomplete = val.(int)
		case "peers":
			switch v := val.(type) {
			case string:
				bin := []byte(v)
				chunks, err := chunkPeers(bin)
				if err != nil {
					return nil, err
				}
				AllPeerList.AddPeers(parseCompactChunks(chunks))
			case []any:
				peersList := make([]PeerInfo, 0, len(v))
				for _, p := range v {
					pdict := p.(map[string]any)
					peer := PeerInfo{}
					peer.PeerID = pdict["peer id"].(string)
					peer.IP = pdict["ip"].(net.IP)
					peer.Port = pdict["port"].(uint16)
					peersList = append(peersList, peer)
				}
				AllPeerList.AddPeers(peersList)
			}
		}
	}
	return resp, nil
}

func chunkPeers(data []byte) ([][6]byte, error) {
	if len(data)%6 != 0 {
		return nil, fmt.Errorf("invalid compact peer length: %d", len(data))
	}
	var chunks [][6]byte
	for i := 0; i < len(data); i += 6 {
		var block [6]byte
		copy(block[:], data[i:i+6])
		chunks = append(chunks, block)
	}
	return chunks, nil
}

func parseCompactChunks(chunks [][6]byte) []PeerInfo {
	var peers []PeerInfo
	for _, b := range chunks {
		ip := net.IP(b[0:4])
		port := binary.BigEndian.Uint16(b[4:6])
		peers = append(peers, PeerInfo{IP: ip, Port: port})
	}
	return peers
}
