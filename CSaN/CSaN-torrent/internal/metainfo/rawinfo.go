package metainfo

import (
	"fmt"
)

// RawInfoDict returns the exact bencoded bytes of the top-level "info"
// dictionary from a .torrent file. This is required for correct info_hash.
func RawInfoDict(torrentFile []byte) ([]byte, error) {
	idx := findInfoDictStart(torrentFile)
	if idx < 0 {
		return nil, fmt.Errorf("info dict not found")
	}
	end, err := skipValue(torrentFile, idx)
	if err != nil {
		return nil, err
	}
	return torrentFile[idx:end], nil
}

func findInfoDictStart(data []byte) int {
	marker := []byte("4:info")
	i := 0
	for {
		j := indexFrom(data, marker, i)
		if j < 0 {
			return -1
		}
		pos := j + len(marker)
		if pos < len(data) && data[pos] == 'd' {
			return pos
		}
		i = j + 1
	}
}

func indexFrom(haystack, needle []byte, from int) int {
	if from >= len(haystack) {
		return -1
	}
	idx := bytesIndex(haystack[from:], needle)
	if idx < 0 {
		return -1
	}
	return from + idx
}

func bytesIndex(b, sub []byte) int {
outer:
	for i := 0; i+len(sub) <= len(b); i++ {
		for j := 0; j < len(sub); j++ {
			if b[i+j] != sub[j] {
				continue outer
			}
		}
		return i
	}
	return -1
}

// skipValue returns the index just past the end of the bencoded value
// starting at pos (pos must point at d, l, i, or digit).
func skipValue(data []byte, pos int) (int, error) {
	if pos >= len(data) {
		return 0, fmt.Errorf("unexpected EOF")
	}
	switch data[pos] {
	case 'd':
		p := pos + 1
		for p < len(data) && data[p] != 'e' {
			// key
			kn, err := skipString(data, p)
			if err != nil {
				return 0, err
			}
			p += kn
			// value
			vn, err := skipValueSize(data, p)
			if err != nil {
				return 0, err
			}
			p += vn
		}
		if p >= len(data) || data[p] != 'e' {
			return 0, fmt.Errorf("unterminated dict")
		}
		return p + 1, nil
	case 'l':
		p := pos + 1
		for p < len(data) && data[p] != 'e' {
			n, err := skipValueSize(data, p)
			if err != nil {
				return 0, err
			}
			p += n
		}
		if p >= len(data) || data[p] != 'e' {
			return 0, fmt.Errorf("unterminated list")
		}
		return p + 1, nil
	case 'i':
		for p := pos + 1; p < len(data); p++ {
			if data[p] == 'e' {
				return p + 1, nil
			}
		}
		return 0, fmt.Errorf("unterminated int")
	default:
		if data[pos] >= '0' && data[pos] <= '9' {
			n, err := skipString(data, pos)
			if err != nil {
				return 0, err
			}
			return pos + n, nil
		}
		return 0, fmt.Errorf("invalid value at %d", pos)
	}
}

func skipValueSize(data []byte, pos int) (int, error) {
	end, err := skipValue(data, pos)
	if err != nil {
		return 0, err
	}
	return end - pos, nil
}

func skipString(data []byte, pos int) (int, error) {
	colon := -1
	for i := pos; i < len(data); i++ {
		if data[i] == ':' {
			colon = i
			break
		}
		if data[i] < '0' || data[i] > '9' {
			return 0, fmt.Errorf("invalid string")
		}
	}
	if colon < 0 {
		return 0, fmt.Errorf("missing colon")
	}
	lengthStr := string(data[pos:colon])
	var length int
	for _, c := range lengthStr {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid length")
		}
		length = length*10 + int(c-'0')
	}
	end := colon + 1 + length
	if end > len(data) {
		return 0, fmt.Errorf("string truncated")
	}
	return end - pos, nil
}
