package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	store = make(map[string]string)
	mu    sync.RWMutex
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/key/")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	mu.RLock()
	value, ok := store[key]
	mu.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK) 
	w.Write([]byte(value))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/key/")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	mu.Lock()
	store[key] = string(body)
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/key/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getHandler(w, r)
		case http.MethodPut:
			putHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
