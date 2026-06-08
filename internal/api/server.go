package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/baltej223/dukedb/internal/node"
)

type Server struct {
	Addr string
	me   *node.Node
}

func NewServer(addr string, me *node.Node) *Server {
	return &Server{
		Addr: addr,
		me:   me,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/health",
		s.handleHealth,
	)

	mux.HandleFunc(
		"/put",
		s.handlePut,
	)

	mux.HandleFunc(
		"/get",
		s.handleGet,
	)

	log.Printf(
		"HTTP API listening on %s",
		s.Addr,
	)

	return http.ListenAndServe(
		s.Addr,
		mux,
	)
}

func (s *Server) handleHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(
		[]byte("OK"),
	)
}

type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PutResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (s *Server) handlePut(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var req PutRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&req)
	if err != nil {
		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)
		return
	}

	var resp PutResponse
	err = node.PUT(req.Key, req.Value, s.me)
	if err != nil {
		resp = PutResponse{
			Success: false,
			Error:   string(err.Error()),
		}
	} else {
		resp = PutResponse{
			Success: true,
		}
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(resp)
}

type GetResponse struct {
	Found bool   `json:"found"`
	Value string `json:"value,omitempty"`
	Error string `json:"error,omitempty"`
}

func (s *Server) handleGet(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(
			w,
			"missing key",
			http.StatusBadRequest,
		)
		return
	}

	// TODO:
	// call Duke here
	value, err := node.GET(key, s.me)
	var resp GetResponse

	if err != nil {
		resp = GetResponse{
			Found: false,
			Value: "",
		}
	} else {
		resp = GetResponse{
			Found: true,
			Value: value,
		}
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(resp)
}
