package http

import (
	"encoding/hex"
	"net/http"

	"github.com/sirupsen/logrus"
	"torrent/internal/bencode"
	"torrent/internal/tracker"
)

type Handlers struct {
	log     *logrus.Logger
	tracker *tracker.Service
	store   *tracker.Store
}

func NewHandlers(log *logrus.Logger, service *tracker.Service, store *tracker.Store) *Handlers {
	if log == nil {
		log = logrus.StandardLogger()
	}
	return &Handlers{log: log, tracker: service, store: store}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/announce", h.announce)
	mux.HandleFunc("/scrape", h.scrape)
}

func (h *Handlers) healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handlers) announce(w http.ResponseWriter, r *http.Request) {
	req, err := h.tracker.ParseAnnounce(r)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"remote": r.RemoteAddr,
			"path":   r.URL.Path,
			"reason": err.Error(),
		}).Warn("announce rejected")
		h.failure(w, http.StatusBadRequest, "invalid announce request")
		return
	}
	resp := h.tracker.Announce(req)

	h.log.WithFields(logrus.Fields{
		"info_hash":      hex.EncodeToString([]byte(req.InfoHash)),
		"event":          string(req.Event),
		"peer_port":      req.Port,
		"left":           req.Left,
		"remote":         r.RemoteAddr,
		"peers_returned": len(resp.Peers),
		"complete":       resp.Complete,
		"incomplete":     resp.Incomplete,
	}).Info("announce")

	payload := map[string]any{
		"interval":   resp.Interval,
		"complete":   resp.Complete,
		"incomplete": resp.Incomplete,
		"peers":      tracker.CompactPeers(resp.Peers),
	}
	h.writeBencoded(w, http.StatusOK, payload)
}

func (h *Handlers) scrape(w http.ResponseWriter, r *http.Request) {
	infoHash := r.URL.Query().Get("info_hash")

	files := map[string]any{}
	if infoHash != "" {
		snap := h.store.Snapshot(infoHash)
		files[infoHash] = map[string]any{
			"complete":   snap.Complete,
			"incomplete": snap.Incomplete,
			"downloaded": snap.Downloaded,
		}
		h.log.WithFields(logrus.Fields{
			"remote":     r.RemoteAddr,
			"scope":      "single",
			"info_hash":  hex.EncodeToString([]byte(infoHash)),
			"complete":   snap.Complete,
			"incomplete": snap.Incomplete,
			"downloaded": snap.Downloaded,
		}).Info("scrape")
	} else {
		all := h.store.ScrapeAll()
		for ih, snap := range all {
			files[ih] = map[string]any{
				"complete":   snap.Complete,
				"incomplete": snap.Incomplete,
				"downloaded": snap.Downloaded,
			}
		}
		h.log.WithFields(logrus.Fields{
			"remote":   r.RemoteAddr,
			"scope":    "all",
			"torrents": len(all),
		}).Info("scrape")
	}

	h.writeBencoded(w, http.StatusOK, map[string]any{"files": files})
}

func (h *Handlers) failure(w http.ResponseWriter, code int, reason string) {
	h.writeBencoded(w, code, map[string]any{"failure reason": reason})
}

func (h *Handlers) writeBencoded(w http.ResponseWriter, code int, payload map[string]any) {
	data, err := bencode.Marshal(payload)
	if err != nil {
		h.log.WithError(err).Error("bencode marshal failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	_, _ = w.Write(data)
}
