package tracker

import (
	"testing"
	"time"
)

func TestRemoveStale(t *testing.T) {
	store := NewStore(1)
	store.UpsertPeer("hash", &Peer{
		PeerID:   "oldpeer",
		IP:       "127.0.0.1",
		Port:     6881,
		LastSeen: time.Now().Add(-10 * time.Minute),
	})
	store.UpsertPeer("hash", &Peer{
		PeerID:   "newpeer",
		IP:       "127.0.0.2",
		Port:     6882,
		Left:     10,
		LastSeen: time.Now().Add(-10 * time.Second),
	})

	removed := store.RemoveStale(time.Now().Add(-2 * time.Minute))
	if removed != 1 {
		t.Fatalf("RemoveStale() = %d, want 1", removed)
	}

	snap := store.Snapshot("hash")
	if snap.Complete != 0 || snap.Incomplete != 1 {
		t.Fatalf("unexpected snapshot after cleanup: %+v", snap)
	}
}
