package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

const baseUrl = "http://localhost:8080/"

var storage = make(map[string]string)

func handlePost(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "text/plain") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originUrl := strings.TrimSpace(string(body))
	if originUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedUrl, err := url.ParseRequestURI(originUrl)
	if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := generateShortUrl()
	storage[id] = originUrl

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(baseUrl + id))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originUrl, ok := storage[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func generateShortUrl() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 5

	shortNewUrl := make([]byte, length)
	for i := range shortNewUrl {
		shortNewUrl[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortNewUrl)
}

func main() {
	r := chi.NewRouter()

	r.Post("/", handlePost)
	r.Get("/{id}", handleGet)

	log.Fatal(http.ListenAndServe(`:8080`, r))
}
