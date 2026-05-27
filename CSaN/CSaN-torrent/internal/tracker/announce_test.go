package tracker

import (
	"net/http/httptest"
	"testing"
)

func TestParseAnnounce(t *testing.T) {
	store := NewStore(1)
	svc := NewService(store, 120, 30, 100)

	req := httptest.NewRequest("GET", "/announce?info_hash=12345678901234567890&peer_id=ABCDEFGHIJKLMNOPQRST&port=6881&uploaded=1&downloaded=2&left=3&event=started", nil)
	req.RemoteAddr = "10.0.0.2:1234"

	parsed, err := svc.ParseAnnounce(req)
	if err != nil {
		t.Fatalf("ParseAnnounce() error = %v", err)
	}
	if parsed.IP != "10.0.0.2" {
		t.Fatalf("parsed.IP = %q, want 10.0.0.2", parsed.IP)
	}
	if parsed.Event != EventStarted {
		t.Fatalf("parsed.Event = %q, want started", parsed.Event)
	}
}

func TestAnnounceLifecycle(t *testing.T) {
	store := NewStore(1)
	svc := NewService(store, 120, 30, 100)
	ih := "12345678901234567890"

	respA := svc.Announce(AnnounceRequest{
		InfoHash: ih, PeerID: "AAAAAAAAAAAAAAAAAAAA", IP: "127.0.0.1", Port: 6881, Left: 10, NumWant: 50, Event: EventStarted,
	})
	if respA.Incomplete != 1 || respA.Complete != 0 {
		t.Fatalf("unexpected counters after first peer: %+v", respA)
	}

	respB := svc.Announce(AnnounceRequest{
		InfoHash: ih, PeerID: "BBBBBBBBBBBBBBBBBBBB", IP: "127.0.0.2", Port: 6882, Left: 0, NumWant: 50, Event: EventStarted,
	})
	if len(respB.Peers) != 1 {
		t.Fatalf("len(peers) = %d, want 1", len(respB.Peers))
	}
	if respB.Complete != 1 || respB.Incomplete != 1 {
		t.Fatalf("unexpected counters with two peers: %+v", respB)
	}

	_ = svc.Announce(AnnounceRequest{
		InfoHash: ih, PeerID: "AAAAAAAAAAAAAAAAAAAA", IP: "127.0.0.1", Port: 6881, Left: 0, NumWant: 50, Event: EventCompleted,
	})
	snap := store.Snapshot(ih)
	if snap.Downloaded != 1 {
		t.Fatalf("Downloaded = %d, want 1", snap.Downloaded)
	}

	_ = svc.Announce(AnnounceRequest{
		InfoHash: ih, PeerID: "BBBBBBBBBBBBBBBBBBBB", Event: EventStopped,
	})
	snap = store.Snapshot(ih)
	if snap.Complete != 1 || snap.Incomplete != 0 {
		t.Fatalf("unexpected counters after stop: %+v", snap)
	}
}

func TestCompactPeers(t *testing.T) {
	b := CompactPeers([]CompactPeer{
		{IP: "127.0.0.1", Port: 6881},
	})
	if len(b) != 6 {
		t.Fatalf("compact length = %d, want 6", len(b))
	}
}
