package main

import (
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const baseUrl = "http://localhost:8080/"

var storage = make(map[string]string)

func handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/":
		handlePost(w, r)
		return
	case r.Method == http.MethodGet && r.URL.Path != "/":
		handleGet(w, r)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
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

	shortUrl := baseUrl + id

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortUrl))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
