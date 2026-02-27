package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"
)

const ApplicationName = "strx"

var urlStore Store

func main() {
	port := flag.String("port", "3000", "Port to run the server on")
	store := flag.String("store", "file", "Store type: memory or file")
	flag.Parse()

	if *store == "file" {
		fileStore := NewFilesStore(ApplicationName)
		log.Printf("Using file store: %s", fileStore.directoryPath)
		urlStore = fileStore
	} else {
		log.Printf("Using in-memory store")
		urlStore = NewInMemoryStore()
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("POST /create", handleCreate)
	mux.HandleFunc("GET /{alias}", handleResolve)

	handler := loggingMiddleware(mux)

	addr := ":" + *port
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func acceptsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Only match exact root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	all := urlStore.All()

	if acceptsJSON(r) {
		writeJSON(w, http.StatusOK, all)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := HTML(all, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		URL   string `json:"url"`
		Alias string `json:"alias,omitempty"`
	}

	req := new(Request)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if req.URL == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "URL is required"})
		return
	}

	alias := req.Alias
	if alias == "" {
		word1, word2 := randomWords()
		alias = word1 + "-" + word2
	}

	urlStore.Set(alias, req.URL)
	writeJSON(w, http.StatusOK, map[string]string{"alias": alias, "url": req.URL})
}

func handleResolve(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("alias")

	url, exists := urlStore.Get(alias)
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Alias not found"})
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
