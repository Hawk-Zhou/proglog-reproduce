// my reflection:
// the checks of un-marshalling in both handleProduce and handleConsume
// are kindle of redundant and should be handled by a middleware or something
// -> after tried, maybe not a good idea to multiplexing these

package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type httpServer struct {
	Log *Log
}

// below are structs to be handled/delivered by http Server
type ProduceRequest struct {
	Record Record `json:"record"`
}
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}
type ConsumeRequest struct {
	// pointer type is a trick to distinguish null value
	Offset *uint64 `json:"offset"`
}
type ConsumeResponse struct {
	Record Record `json:"record"`
}

// handleProduce implements httpHandlerFunc
func (hs *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Record.Value) == 0 {
		http.Error(w, "record should be non-empty", http.StatusBadRequest)
		return
	}

	offset, err := hs.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp = ProduceResponse{Offset: offset}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleConsume implements httpHandlerFunc
func (hs *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Offset == nil {
		http.Error(w, "offset should be non-empty", http.StatusBadRequest)
		return
	}

	rec, err := hs.Log.Read(*req.Offset)
	if err != nil {
		if err == ErrOffsetNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp = ConsumeResponse{Record: rec}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// newHTTPServer inits internal struct httpServer
func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

var _ http.Handler = (*mux.Router)(nil)

func NewHTTPServer(addr string) *http.Server {
	srv := newHTTPServer()
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	})
	r.HandleFunc("/", srv.handleConsume).Methods("GET")
	r.HandleFunc("/", srv.handleProduce).Methods("POST")
	return &http.Server{
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 120,

		Addr:    addr,
		Handler: r,
	}
}
