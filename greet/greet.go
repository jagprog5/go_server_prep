package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type GreetRequest struct {
	Name string `json:"name"`
	// Version string `json:"version"`
}

type GreetResponse struct {
	Message string `json:"message"`
}

// curl -XPOST http://localhost:8080/greet  -H "Content-Type: application/json" -H "Accept: application/json" -d '{"name": "Alice"}'
func post_greet(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	accept := r.Header.Get("Accept")
	if accept != "" && !strings.Contains(accept, "application/json") {
		http.Error(w, "Accept header must include application/json", http.StatusNotAcceptable)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req GreetRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "empty name", http.StatusBadRequest)
		return
	}

	var response GreetResponse
	response.Message = req.Name + ", hi!"

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("POST /greet", post_greet)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
