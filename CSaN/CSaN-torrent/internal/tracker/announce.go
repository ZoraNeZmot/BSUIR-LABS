package tracker

import (
	"encoding/binary"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var errMissingParam = errors.New("missing required parameter")

type Service struct {
	store          *Store
	interval       int64
	defaultNumWant int
	maxNumWant     int
}

func NewService(store *Store, intervalSec int64, defaultNumWant, maxNumWant int) *Service {
	return &Service{
		store:          store,
		interval:       intervalSec,
		defaultNumWant: defaultNumWant,
		maxNumWant:     maxNumWant,
	}
}

type AnnounceRequest struct {
	InfoHash   string
	PeerID     string
	IP         string
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	NumWant    int
	Event      Event
}

type AnnounceResponse struct {
	Interval   int64
	Complete   int64
	Incomplete int64
	Peers      []CompactPeer
}

func (s *Service) ParseAnnounce(r *http.Request) (AnnounceRequest, error) {
	query := r.URL.Query()

	infoHash, err := requiredRaw(query, "info_hash", 20)
	if err != nil {
		return AnnounceRequest{}, err
	}
	peerID, err := requiredRaw(query, "peer_id", 20)
	if err != nil {
		return AnnounceRequest{}, err
	}

	port, err := requiredInt(query, "port")
	if err != nil {
		return AnnounceRequest{}, err
	}
	uploaded, err := requiredInt64(query, "uploaded")
	if err != nil {
		return AnnounceRequest{}, err
	}
	downloaded, err := requiredInt64(query, "downloaded")
	if err != nil {
		return AnnounceRequest{}, err
	}
	left, err := requiredInt64(query, "left")
	if err != nil {
		return AnnounceRequest{}, err
	}

	numWant := s.defaultNumWant
	if raw := query.Get("numwant"); raw != "" {
		if parsed, parseErr := strconv.Atoi(raw); parseErr == nil && parsed > 0 {
			numWant = parsed
		}
	}
	if numWant > s.maxNumWant {
		numWant = s.maxNumWant
	}

	ip := clientIP(r)
	if queryIP := query.Get("ip"); queryIP != "" && net.ParseIP(queryIP) != nil {
		ip = queryIP
	}

	return AnnounceRequest{
		InfoHash:   infoHash,
		PeerID:     peerID,
		IP:         ip,
		Port:       port,
		Uploaded:   uploaded,
		Downloaded: downloaded,
		Left:       left,
		NumWant:    numWant,
		Event:      parseEvent(query.Get("event")),
	}, nil
}

func (s *Service) Announce(req AnnounceRequest) AnnounceResponse {
	if req.Event == EventStopped {
		s.store.RemovePeer(req.InfoHash, req.PeerID)
	} else {
		peer := &Peer{
			PeerID:     req.PeerID,
			IP:         req.IP,
			Port:       req.Port,
			Uploaded:   req.Uploaded,
			Downloaded: req.Downloaded,
			Left:       req.Left,
			LastSeen:   time.Now().UTC(),
		}
		s.store.UpsertPeer(req.InfoHash, peer)
		if req.Event == EventCompleted {
			s.store.MarkCompleted(req.InfoHash, req.PeerID)
		}
	}

	snap := s.store.Snapshot(req.InfoHash)
	peers := s.store.SelectPeers(req.InfoHash, req.PeerID, req.NumWant)

	return AnnounceResponse{
		Interval:   s.interval,
		Complete:   snap.Complete,
		Incomplete: snap.Incomplete,
		Peers:      peers,
	}
}

func CompactPeers(peers []CompactPeer) []byte {
	out := make([]byte, 0, len(peers)*6)
	for _, p := range peers {
		ip := net.ParseIP(p.IP).To4()
		if ip == nil {
			continue
		}
		out = append(out, ip...)
		port := make([]byte, 2)
		binary.BigEndian.PutUint16(port, uint16(p.Port))
		out = append(out, port...)
	}
	return out
}

func requiredRaw(values url.Values, key string, expectedLen int) (string, error) {
	raw, ok := values[key]
	if !ok || len(raw) == 0 || raw[0] == "" {
		return "", errMissingParam
	}
	value := raw[0]
	if expectedLen > 0 && len(value) != expectedLen {
		return "", errors.New("invalid parameter length")
	}
	return value, nil
}

func requiredInt(values url.Values, key string) (int, error) {
	raw := values.Get(key)
	if raw == "" {
		return 0, errMissingParam
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 0 {
		return 0, errors.New("invalid integer parameter")
	}
	return v, nil
}

func requiredInt64(values url.Values, key string) (int64, error) {
	raw := values.Get(key)
	if raw == "" {
		return 0, errMissingParam
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || v < 0 {
		return 0, errors.New("invalid integer parameter")
	}
	return v, nil
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "127.0.0.1"
	}
	if net.ParseIP(host) == nil {
		return "127.0.0.1"
	}
	return host
}
