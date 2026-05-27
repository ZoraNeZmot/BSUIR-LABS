package metainfo

import (
	"crypto/sha1"
	"fmt"
	"os"
	"strings"

	"torrent/internal/bencode"
)

// Torrent holds parsed metadata from a .torrent file.
type Torrent struct {
	Announce    string
	InfoHash    [20]byte
	PieceLen    int64
	NumPieces   int
	PieceHashes [][]byte // each 20 bytes
	TotalLen    int64
	Name        string
	SingleFile  bool
	Files       []FileEntry // empty when SingleFile
}

type FileEntry struct {
	Length int64
	Path   []string
}

// Load reads a .torrent file and parses metainfo.
func Load(path string) (*Torrent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rawInfo, err := RawInfoDict(data)
	if err != nil {
		return nil, err
	}
	sum := sha1.Sum(rawInfo)

	rootVal, _, err := bencode.Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("decode torrent: %w", err)
	}
	root, ok := rootVal.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("torrent root is not a dict")
	}

	announce, err := stringField(root, "announce")
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(announce, "http://") && !strings.HasPrefix(announce, "https://") {
		return nil, fmt.Errorf("unsupported announce scheme (need http/https): %q", announce)
	}

	infoVal, _, err := bencode.Unmarshal(rawInfo)
	if err != nil {
		return nil, fmt.Errorf("decode info: %w", err)
	}
	info, ok := infoVal.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("info is not a dict")
	}

	pieceLen, err := int64Field(info, "piece length")
	if err != nil {
		return nil, err
	}
	piecesBytes, ok := info["pieces"].([]byte)
	if !ok || len(piecesBytes)%20 != 0 {
		return nil, fmt.Errorf("invalid pieces field")
	}
	numPieces := len(piecesBytes) / 20
	pieceHashes := make([][]byte, numPieces)
	for i := 0; i < numPieces; i++ {
		h := make([]byte, 20)
		copy(h, piecesBytes[i*20:(i+1)*20])
		pieceHashes[i] = h
	}

	t := &Torrent{
		Announce:    announce,
		InfoHash:    sum,
		PieceLen:    pieceLen,
		NumPieces:   numPieces,
		PieceHashes: pieceHashes,
	}

	if _, ok := info["files"]; ok {
		t.SingleFile = false
		name, err := stringField(info, "name")
		if err != nil {
			return nil, err
		}
		t.Name = name
		filesVal, ok := info["files"].([]any)
		if !ok {
			return nil, fmt.Errorf("invalid files list")
		}
		var total int64
		for _, item := range filesVal {
			fd, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid file entry")
			}
			fl, err := int64Field(fd, "length")
			if err != nil {
				return nil, err
			}
			pathParts, err := pathListField(fd, "path")
			if err != nil {
				return nil, err
			}
			t.Files = append(t.Files, FileEntry{Length: fl, Path: pathParts})
			total += fl
		}
		t.TotalLen = total
	} else {
		t.SingleFile = true
		length, err := int64Field(info, "length")
		if err != nil {
			return nil, err
		}
		name, err := stringField(info, "name")
		if err != nil {
			return nil, err
		}
		t.Name = name
		t.TotalLen = length
	}

	if err := validatePieceLayout(t.TotalLen, pieceLen, numPieces); err != nil {
		return nil, err
	}

	return t, nil
}

func validatePieceLayout(total, pieceLen int64, numPieces int) error {
	if numPieces < 1 || pieceLen < 1 || total < 1 {
		return fmt.Errorf("invalid piece layout")
	}
	last := total - pieceLen*int64(numPieces-1)
	if last < 1 || last > pieceLen {
		return fmt.Errorf("piece geometry mismatch: total=%d pieces=%d pieceLen=%d last=%d", total, numPieces, pieceLen, last)
	}
	return nil
}

func stringField(m map[string]any, key string) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", fmt.Errorf("missing %q", key)
	}
	switch s := v.(type) {
	case string:
		return s, nil
	case []byte:
		return string(s), nil
	default:
		return "", fmt.Errorf("field %q has wrong type", key)
	}
}

func int64Field(m map[string]any, key string) (int64, error) {
	v, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("missing %q", key)
	}
	switch n := v.(type) {
	case int64:
		return n, nil
	case int:
		return int64(n), nil
	default:
		return 0, fmt.Errorf("field %q has wrong type", key)
	}
}

func pathListField(m map[string]any, key string) ([]string, error) {
	v, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("missing %q", key)
	}
	list, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("path is not a list")
	}
	out := make([]string, 0, len(list))
	for _, part := range list {
		switch p := part.(type) {
		case []byte:
			out = append(out, string(p))
		case string:
			out = append(out, p)
		default:
			return nil, fmt.Errorf("invalid path part type")
		}
	}
	return out, nil
}
