package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const apiPrefix = "/api/"

type ServerData struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewServerData() *ServerData {
	return &ServerData{
		store: make(map[string]string),
	}
}

func (s *ServerData) getterHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, apiPrefix)
	if path == "" || strings.Contains(path, "/") {
		http.Error(w, "malformed get", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	value, found := s.store[path]
	s.mu.RUnlock()

	if !found {
		http.Error(w, strconv.Itoa(http.StatusNotFound)+" request resource not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(value))
}

func (s *ServerData) putterHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, apiPrefix)
	if path == "" || strings.Contains(path, "/") {
		http.Error(w, "malformed put", http.StatusBadRequest)
		return
	}

	value := r.URL.Query().Get("value")

	s.mu.Lock()
	s.store[path] = value
	s.mu.Unlock()

	w.Write([]byte("ok"))
}

func main() {
	s := NewServerData()

	http.HandleFunc("GET "+apiPrefix, s.getterHandler)
	http.HandleFunc("PUT "+apiPrefix, s.putterHandler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
