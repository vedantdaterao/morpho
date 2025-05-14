package torrent

import (
	"net/url"
	"strconv"
)

func (t *TorrentFile) BuildTrackerURL() (string, error) {
    base, err := url.Parse(t.Announce)
    if err != nil {
        return "", err
    }
    params := url.Values{
        "info_hash":  []string{string(t.InfoHash[:])},
        "peer_id":    []string{"XXXXXXXXXXXXXXXXXXXX"},
        "port":       []string{strconv.Itoa(6881)}, 		//[]string{strconv.Itoa(int(Port))}, //  Ports reserved for BitTorrent are typically 6881-6889
        "uploaded":   []string{"0"},
        "downloaded": []string{"0"},
        "left":       []string{strconv.Itoa(t.Length)},
        "compact":    []string{"1"},
    }
    base.RawQuery = params.Encode()
    return base.String(), nil
}

