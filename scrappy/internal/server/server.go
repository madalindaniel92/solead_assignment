package server

import (
	"net/http"
	"time"

	"examples/scrappy/internal/es"
)

type State struct {
	client *es.Client
}

func NewServer(addr string, timeout time.Duration, client *es.Client) *http.Server {
	state := &State{client: client}

	return &http.Server{
		Addr:         addr,
		Handler:      router(state),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
}

func router(state *State) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/companies", companiesHandler(state))
	return mux
}
