package bencoding

import (
	"errors"
	"fmt"
	"strconv"
)

func (d *Parser) decode(raw bool) (any, error) {
	switch d.data[d.pos] {
	case 'i':
		digits, err := d.decodeInt()
		if err != nil {
			return nil, err
		}
		return digits, nil
	case 'l':
		return d.decodeList(raw)
	case 'd':
		return d.decodeDict(raw)
	default:
		if d.data[d.pos] >= '0' && d.data[d.pos] <= '9' {
			return d.decodeString()
		}
		return nil, fmt.Errorf("unexpected character at pos %d: %c", d.pos, d.data[d.pos])
	}
}

func (d *Parser) decodeString() (string, error) {
	start := d.pos
	for d.pos < len(d.data) && d.data[d.pos] != ':' {
		d.pos++
	}

	l, err := strconv.Atoi(string(d.data[start:d.pos]))
	if err != nil {
		return "", err
	}
	d.pos++ //skip ':'
	if d.pos+l > len(d.data) {
		return "", errors.New("error")
	}

	str := string(d.data[d.pos : d.pos+l])
	d.pos += l
	return str, nil
}

func (d *Parser) decodeInt() (int, error) {
	d.pos++
	start := d.pos
	for d.pos < len(d.data) && d.data[d.pos] != 'e' {
		if d.data[d.pos] == '-' && start == d.pos {
			// Allow negative sign at the start
			d.pos++
		}
		if d.data[d.pos] < '0' || d.data[d.pos] > '9' {
			return 0, fmt.Errorf("invalid character in integer at pos %d: %c", d.pos, d.data[d.pos])
		}
		d.pos++
	}

	numStr := string(d.data[start:d.pos])
	if d.pos >= len(d.data) || d.data[d.pos] != 'e' || len(numStr) == 0 {
		return 0, errors.New("error parsing integer")
	}
	d.pos++ //skip 'e'
	return strconv.Atoi(numStr)
}

func (d *Parser) decodeList(raw bool) ([]interface{}, error) {
	d.pos++
	var list []interface{}
	for d.pos < len(d.data) && d.data[d.pos] != 'e' {
		item, err := d.decode(raw)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	if d.pos >= len(d.data) {
		return nil, errors.New("unterminated list")
	}

	d.pos++ //skip 'e'
	return list, nil
}

func (d *Parser) decodeDict(raw bool) (map[string]interface{}, error) {
	start := d.pos
	d.pos++
	dict := make(map[string]interface{})
	var afterPieces int
	for d.pos < len(d.data) && d.data[d.pos] != 'e' {
		key, err := d.decodeString()
		if err != nil {
			return nil, err
		}
		val, err := d.decode(raw)
		if err != nil {
			return nil, err
		}
		dict[key] = val

		// Record position immediately after 'pieces' value
		if raw {
			if key == "pieces" {
				afterPieces = d.pos
				dict["after_pieces"] = d.data[afterPieces:]
			}
		}
	}

	if d.pos >= len(d.data) {
		return nil, errors.New("unterminated dictionary")
	}

	d.pos++ // skip 'e'
	if raw {
		dict["raw"] = d.data[start:d.pos]
	}

	return dict, nil
}

func Decode(data []byte, opts ...bool) (any, error) {
	raw := false
	if len(opts) > 0 {
		raw = opts[0]
	}
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	d := Parser{data, 0}
	return d.decode(raw)
}
