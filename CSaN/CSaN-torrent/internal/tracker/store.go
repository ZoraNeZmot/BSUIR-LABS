package tracker

import (
	"math/rand"
	"sync"
	"time"
)

type Store struct {
	mu     sync.RWMutex
	swarms map[string]*Swarm
	rnd    *rand.Rand
}

func NewStore(seed int64) *Store {
	return &Store{
		swarms: make(map[string]*Swarm),
		rnd:    rand.New(rand.NewSource(seed)),
	}
}

func (s *Store) EnsureSwarm(infoHash string) *Swarm {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ensureSwarmLocked(infoHash)
}

func (s *Store) ensureSwarmLocked(infoHash string) *Swarm {
	sw, ok := s.swarms[infoHash]
	if !ok {
		sw = &Swarm{Peers: make(map[string]*Peer)}
		s.swarms[infoHash] = sw
	}
	return sw
}

func (s *Store) UpsertPeer(infoHash string, peer *Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sw := s.ensureSwarmLocked(infoHash)
	if current, ok := sw.Peers[peer.PeerID]; ok && current.Completed {
		peer.Completed = true
	}
	sw.Peers[peer.PeerID] = peer
}

func (s *Store) RemovePeer(infoHash, peerID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sw, ok := s.swarms[infoHash]
	if !ok {
		return
	}
	delete(sw.Peers, peerID)
	if len(sw.Peers) == 0 {
		delete(s.swarms, infoHash)
	}
}

func (s *Store) MarkCompleted(infoHash, peerID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sw, ok := s.swarms[infoHash]
	if !ok {
		return
	}
	peer, ok := sw.Peers[peerID]
	if !ok {
		return
	}
	if !peer.Completed {
		sw.Downloaded++
		peer.Completed = true
	}
}

func (s *Store) Snapshot(infoHash string) Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sw, ok := s.swarms[infoHash]
	if !ok {
		return Snapshot{}
	}

	var complete int64
	var incomplete int64
	for _, p := range sw.Peers {
		if p.IsSeeder() {
			complete++
		} else {
			incomplete++
		}
	}
	return Snapshot{
		Complete:   complete,
		Incomplete: incomplete,
		Downloaded: sw.Downloaded,
	}
}

func (s *Store) SelectPeers(infoHash, excludePeerID string, numWant int) []CompactPeer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sw, ok := s.swarms[infoHash]
	if !ok {
		return nil
	}

	candidates := make([]CompactPeer, 0, len(sw.Peers))
	for _, p := range sw.Peers {
		if p.PeerID == excludePeerID {
			continue
		}
		candidates = append(candidates, CompactPeer{IP: p.IP, Port: p.Port})
	}

	if len(candidates) <= numWant {
		return candidates
	}

	s.rnd.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	return candidates[:numWant]
}

func (s *Store) RemoveStale(cutoff time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	removed := 0
	for infoHash, sw := range s.swarms {
		for peerID, p := range sw.Peers {
			if p.LastSeen.Before(cutoff) {
				delete(sw.Peers, peerID)
				removed++
			}
		}
		if len(sw.Peers) == 0 {
			delete(s.swarms, infoHash)
		}
	}
	return removed
}

func (s *Store) ScrapeAll() map[string]Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]Snapshot, len(s.swarms))
	for infoHash, sw := range s.swarms {
		var complete int64
		var incomplete int64
		for _, p := range sw.Peers {
			if p.IsSeeder() {
				complete++
			} else {
				incomplete++
			}
		}
		out[infoHash] = Snapshot{
			Complete:   complete,
			Incomplete: incomplete,
			Downloaded: sw.Downloaded,
		}
	}
	return out
}
