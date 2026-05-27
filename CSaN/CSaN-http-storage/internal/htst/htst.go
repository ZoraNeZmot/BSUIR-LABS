package htst

import (
	"context"
	"htst/pkg/config"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

/*
HTST is a HTTP based file server.
HTST is a server that listens for HTTP based requests and return files from the local filesystem.
*/
type HTST struct {
	server *http.Server
	mux    *http.ServeMux
	config *config.Config
}

func New(config *config.Config) *HTST {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    net.JoinHostPort(config.App.Listen.Host, config.App.Listen.Port),
		Handler: mux,
	}
	// mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/", HTSTHandler)
	return &HTST{
		server: server,
		mux:    mux,
		config: config,
	}
}

func (h *HTST) Start() error {
	return h.server.ListenAndServe()
}
func (h *HTST) Shutdown(context context.Context) error {
	logrus.Info("HTST server shutdown")
	return h.server.Shutdown(context)
}
