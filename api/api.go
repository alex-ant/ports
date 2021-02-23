package api

import (
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/alex-ant/ports/ports"
	"github.com/go-zoo/bone"
)

// PortInfoFetcher defines port info fetcher function.
type PortInfoFetcher func() ([]*ports.PortInfo, error)

// API contains server's settings.
type API struct {
	port     int
	listener net.Listener
	mux      *bone.Mux

	portInfoFetcher PortInfoFetcher
}

// New returns new API.
func New(port int, portInfoFetcher PortInfoFetcher) (*API, error) {
	if portInfoFetcher == nil {
		return nil, errors.New("nil portInfoFetcher provided")
	}

	return &API{
		port:            port,
		portInfoFetcher: portInfoFetcher,
	}, nil
}

func (a *API) defineMux() {
	a.mux = bone.New()

	a.mux.Get("/fetch-port-info", http.HandlerFunc(a.fetchPortInfoHandler))
}

// Start starts the HTTP server.
func (a *API) Start() (err error) {
	a.defineMux()

	a.listener, err = net.Listen("tcp", ":"+strconv.Itoa(a.port))
	if err != nil {
		return
	}

	go http.Serve(a.listener, a.mux)

	return
}

// Stop stops the server.
func (a *API) Stop() {
	a.listener.Close()
}
