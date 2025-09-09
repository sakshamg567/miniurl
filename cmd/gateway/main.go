package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sakshamg567/miniurl/shared/proto/redirectpb"
	"github.com/sakshamg567/miniurl/shared/proto/shortenerpb"
	"google.golang.org/grpc"
)

var (
	shortenerClient shortenerpb.ShortenerServiceClient
	redirectClient  redirectpb.RedirectServiceClient
)

func main() {
	sConn, err := grpc.NewClient("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to Connect shortener: %v", err)
	}
	shortenerClient = shortenerpb.NewShortenerServiceClient(sConn)

	rConn, err := grpc.NewClient("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to Connect redirectSvc: %v", err)
	}

	redirectClient = redirectpb.NewRedirectServiceClient(rConn)

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	log.Println("HTTP Gateway running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		LongURL string `json:"long_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	res, err := shortenerClient.CreateShortURL(
		context.Background(),
		&shortenerpb.CreateShortURLRequest{LongUrl: req.LongURL},
	)
	if err != nil {
		http.Error(w, "shorten failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"short_code": res.ShortCode,
	})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortCode := r.URL.Path[1:]
	if shortCode == "" {
		http.Error(w, "missing short code", http.StatusBadRequest)
		return
	}

	res, err := redirectClient.GetLongURL(
		context.Background(),
		&redirectpb.GetLongURLRequest{ShortCode: shortCode},
	)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, res.LongUrl, http.StatusMovedPermanently)
}
