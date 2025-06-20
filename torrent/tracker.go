package torrent

import (
	"net/url"
	"strconv"
)

func (t *TorrentFile) BuildTrackerURL(customURL ...string) (string, error) {
    announceURL := t.Announce
	if len(customURL) > 0 && customURL[0] != "" {
		announceURL = customURL[0]
	}
    base, err := url.Parse(announceURL)
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
        "numwant":    []string{"50"},
    }
    base.RawQuery = params.Encode()
    return base.String(), nil
}

