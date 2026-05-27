package tracker

import "time"

type Peer struct {
	PeerID     string
	IP         string
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	LastSeen   time.Time
	Completed  bool
}

func (p Peer) IsSeeder() bool {
	return p.Left == 0
}

type Swarm struct {
	Peers      map[string]*Peer
	Downloaded int64
}

type Snapshot struct {
	Complete   int64
	Incomplete int64
	Downloaded int64
}

type CompactPeer struct {
	IP   string
	Port int
}

type Event string

const (
	EventStarted   Event = "started"
	EventStopped   Event = "stopped"
	EventCompleted Event = "completed"
	EventEmpty     Event = ""
)

func parseEvent(raw string) Event {
	switch raw {
	case string(EventStarted):
		return EventStarted
	case string(EventStopped):
		return EventStopped
	case string(EventCompleted):
		return EventCompleted
	default:
		return EventEmpty
	}
}
