package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"torrent/internal/tracker"
)

func discardLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	l.SetLevel(logrus.PanicLevel)
	return l
}

func TestHealthz(t *testing.T) {
	store := tracker.NewStore(1)
	svc := tracker.NewService(store, 120, 30, 100)
	h := NewHandlers(discardLogger(), svc, store)

	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if rr.Body.String() != "ok" {
		t.Fatalf("body = %q, want ok", rr.Body.String())
	}
}

func TestAnnounceAndScrape(t *testing.T) {
	store := tracker.NewStore(1)
	svc := tracker.NewService(store, 120, 30, 100)
	h := NewHandlers(discardLogger(), svc, store)

	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/announce?info_hash=12345678901234567890&peer_id=ABCDEFGHIJKLMNOPQRST&port=6881&uploaded=0&downloaded=0&left=10", nil)
	req.RemoteAddr = "127.0.0.1:1111"
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("announce status = %d, want 200", rr.Code)
	}
	if rr.Body.Len() == 0 {
		t.Fatalf("announce body empty")
	}

	scrapeReq := httptest.NewRequest(http.MethodGet, "/scrape?info_hash=12345678901234567890", nil)
	scrapeRR := httptest.NewRecorder()
	mux.ServeHTTP(scrapeRR, scrapeReq)
	if scrapeRR.Code != http.StatusOK {
		t.Fatalf("scrape status = %d, want 200", scrapeRR.Code)
	}
	if scrapeRR.Body.Len() == 0 {
		t.Fatalf("scrape body empty")
	}
}

func TestAnnounceInvalid(t *testing.T) {
	store := tracker.NewStore(1)
	svc := tracker.NewService(store, 120, 30, 100)
	h := NewHandlers(discardLogger(), svc, store)

	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/announce?peer_id=ABCDEFGHIJKLMNOPQRST", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}
