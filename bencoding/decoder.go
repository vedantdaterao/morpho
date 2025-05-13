package bencoding

import (
	"errors"
	"fmt"
	"strconv"
)

func (d *Parser) decode() (any, error) {
	switch d.data[d.pos] {
		case 'i':
			return d.decodeInt()
		case 'l':
			return d.decodeList()
		case 'd':
			return d.decodeDict()
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
		d.pos++
	}
	numStr := string(d.data[start:d.pos])
	d.pos++ //skip 'e'
	return strconv.Atoi(numStr)
}

func (d *Parser) decodeList() ([]interface{}, error) {
	d.pos++
	var list []interface{}
	for d.pos < len(d.data) && d.data[d.pos] != 'e' {
		item, err := d.decode()
		if err != nil{
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

func (d *Parser) decodeDict() (map[string]interface{}, error) {
	d.pos++
	dict := make(map[string]interface{})
	for d.pos < len(d.data) && d.data[d.pos] != 'e' {
		key, err := d.decodeString()
		if err != nil{
			return nil, err
		}

		val, err := d.decode()
		if err != nil{
			return nil, err
		}
		dict[key] = val
	}
	if d.pos >= len(d.data) {
		return nil, errors.New("unterminated dictionary")
	}

	d.pos++ //skip 'e'
	return dict, nil
}

func Decode(data []byte) (any, error){
	d := Parser{data, 0}
	return d.decode()
}