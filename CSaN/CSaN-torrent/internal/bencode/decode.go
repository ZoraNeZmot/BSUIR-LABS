package bencode

import (
	"fmt"
	"strconv"
)

// Unmarshal parses one bencoded value from data and returns the value and the
// number of bytes consumed. String values are returned as []byte.
func Unmarshal(data []byte) (any, int, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("empty input")
	}
	switch data[0] {
	case 'i':
		return decodeInt(data)
	case 'l':
		return decodeList(data)
	case 'd':
		return decodeDict(data)
	default:
		if data[0] >= '0' && data[0] <= '9' {
			return decodeString(data)
		}
		return nil, 0, fmt.Errorf("unexpected byte %q", data[0])
	}
}

func decodeInt(data []byte) (int64, int, error) {
	if len(data) < 3 || data[0] != 'i' {
		return 0, 0, fmt.Errorf("invalid int")
	}
	end := 1
	for end < len(data) && data[end] != 'e' {
		end++
	}
	if end >= len(data) {
		return 0, 0, fmt.Errorf("unterminated int")
	}
	raw := string(data[1:end])
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return n, end + 1, nil
}

func decodeString(data []byte) ([]byte, int, error) {
	colon := -1
	for i := 0; i < len(data); i++ {
		if data[i] == ':' {
			colon = i
			break
		}
		if data[i] < '0' || data[i] > '9' {
			return nil, 0, fmt.Errorf("invalid string length")
		}
	}
	if colon <= 0 {
		return nil, 0, fmt.Errorf("missing colon in string")
	}
	length, err := strconv.Atoi(string(data[:colon]))
	if err != nil || length < 0 {
		return nil, 0, fmt.Errorf("invalid string length")
	}
	start := colon + 1
	if start+length > len(data) {
		return nil, 0, fmt.Errorf("string truncated")
	}
	return append([]byte(nil), data[start:start+length]...), start + length, nil
}

func decodeList(data []byte) ([]any, int, error) {
	if len(data) < 2 || data[0] != 'l' {
		return nil, 0, fmt.Errorf("invalid list")
	}
	pos := 1
	var out []any
	for pos < len(data) && data[pos] != 'e' {
		v, n, err := Unmarshal(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, v)
		pos += n
	}
	if pos >= len(data) || data[pos] != 'e' {
		return nil, 0, fmt.Errorf("unterminated list")
	}
	return out, pos + 1, nil
}

func decodeDict(data []byte) (map[string]any, int, error) {
	if len(data) < 2 || data[0] != 'd' {
		return nil, 0, fmt.Errorf("invalid dict")
	}
	out := make(map[string]any)
	pos := 1
	for pos < len(data) && data[pos] != 'e' {
		keyBytes, kn, err := decodeString(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		pos += kn
		key := string(keyBytes)
		val, vn, err := Unmarshal(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		out[key] = val
		pos += vn
	}
	if pos >= len(data) || data[pos] != 'e' {
		return nil, 0, fmt.Errorf("unterminated dict")
	}
	return out, pos + 1, nil
}
