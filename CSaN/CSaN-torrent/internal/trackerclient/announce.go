package trackerclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"torrent/internal/bencode"
)

// Peer is IPv4 + port from compact peer list.
type Peer struct {
	IP   string
	Port int
}

// Announce calls the HTTP tracker and returns peers (compact) and interval seconds.
func Announce(announceBase string, infoHash [20]byte, peerID [20]byte, port int, uploaded, downloaded, left int64, event string) ([]Peer, int64, error) {
	u, err := url.Parse(announceBase)
	if err != nil {
		return nil, 0, err
	}
	q := u.Query()
	q.Set("info_hash", string(infoHash[:]))
	q.Set("peer_id", string(peerID[:]))
	q.Set("port", strconv.Itoa(port))
	q.Set("uploaded", strconv.FormatInt(uploaded, 10))
	q.Set("downloaded", strconv.FormatInt(downloaded, 10))
	q.Set("left", strconv.FormatInt(left, 10))
	q.Set("compact", "1")
	if event != "" {
		q.Set("event", event)
	}
	u.RawQuery = q.Encode()

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("tracker HTTP %d: %s", resp.StatusCode, string(body))
	}

	val, _, err := bencode.Unmarshal(body)
	if err != nil {
		return nil, 0, fmt.Errorf("decode tracker response: %w", err)
	}
	dict, ok := val.(map[string]any)
	if !ok {
		return nil, 0, fmt.Errorf("tracker response is not a dict")
	}
	if fr, ok := dict["failure reason"]; ok {
		msg := formatBencodeString(fr)
		return nil, 0, fmt.Errorf("tracker failure: %s", msg)
	}

	interval := int64(120)
	if iv, ok := dict["interval"]; ok {
		switch n := iv.(type) {
		case int64:
			interval = n
		case int:
			interval = int64(n)
		}
	}

	var peersBytes []byte
	switch p := dict["peers"].(type) {
	case []byte:
		peersBytes = p
	case string:
		peersBytes = []byte(p)
	default:
		return nil, 0, fmt.Errorf("missing or invalid compact peers field")
	}
	plist, err := parseCompactPeers(peersBytes)
	if err != nil {
		return nil, 0, err
	}
	return plist, interval, nil
}

func formatBencodeString(v any) string {
	switch s := v.(type) {
	case []byte:
		return string(s)
	case string:
		return s
	default:
		return fmt.Sprint(v)
	}
}

func parseCompactPeers(b []byte) ([]Peer, error) {
	if len(b)%6 != 0 {
		return nil, fmt.Errorf("invalid compact peers length %d", len(b))
	}
	var out []Peer
	for i := 0; i < len(b); i += 6 {
		ip := fmt.Sprintf("%d.%d.%d.%d", b[i], b[i+1], b[i+2], b[i+3])
		port := int(b[i+4])<<8 | int(b[i+5])
		out = append(out, Peer{IP: ip, Port: port})
	}
	return out, nil
}
