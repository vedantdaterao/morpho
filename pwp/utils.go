package pwp

import "fmt"

// pads or truncates s to exactly 20 bytes
func PeerIDFromString(s string) [20]byte {
	var id [20]byte
	copy(id[:], s)
	return id
}

// 40-char hex string into a 20-byte info hash
func InfoHashFromHex(h string) ([20]byte, error) {
	var hash [20]byte
	b, err := hexDecode(h)
	if err != nil || len(b) != 20 {
		return hash, fmt.Errorf("invalid info hash %q", h)
	}
	copy(hash[:], b)
	return hash, nil
}

func hexDecode(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("odd hex length")
	}
	b := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		hi := hexVal(s[i])
		lo := hexVal(s[i+1])
		if hi < 0 || lo < 0 {
			return nil, fmt.Errorf("invalid hex char")
		}
		b[i/2] = byte(hi<<4 | lo)
	}
	return b, nil
}

func hexVal(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return -1
}

// standard request block size, clamped to pieceRemaining
func BlockLength(pieceRemaining int) uint32 {
	const standard = 16 * 1024 // 16 KiB
	if pieceRemaining < standard {
		return uint32(pieceRemaining)
	}
	return standard
}

// index of the piece
func PieceIndex(offset, pieceLen int) int {
	return offset / pieceLen
}

// byte offset within a piece
func PieceOffset(offset, pieceLen int) uint32 {
	return uint32(offset % pieceLen)
}