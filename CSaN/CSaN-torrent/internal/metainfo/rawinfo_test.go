package metainfo

import (
	"crypto/sha1"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRawInfoDictRoundTrip(t *testing.T) {
	pieces := make([]byte, 20)
	info := []byte("d6:lengthi4e4:name4:test12:piece lengthi4e6:pieces20:")
	info = append(info, pieces...)
	info = append(info, 'e')

	announce := "http://127.0.0.1:8080/announce"
	root := []byte("d8:announce")
	root = append(root, []byte(strconv.Itoa(len(announce)))...)
	root = append(root, ':')
	root = append(root, []byte(announce)...)
	root = append(root, []byte("4:info")...)
	root = append(root, info...)
	root = append(root, 'e')

	tmp := filepath.Join(t.TempDir(), "test.torrent")
	if err := os.WriteFile(tmp, root, 0o644); err != nil {
		t.Fatal(err)
	}

	raw, err := RawInfoDict(root)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != string(info) {
		t.Fatalf("raw info mismatch")
	}
	sum := sha1.Sum(raw)
	mt, err := Load(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if mt.InfoHash != sum {
		t.Fatalf("info hash mismatch")
	}
	if mt.TotalLen != 4 {
		t.Fatalf("total %d", mt.TotalLen)
	}
}
