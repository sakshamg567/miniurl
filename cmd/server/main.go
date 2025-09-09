package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/sakshamg567/miniurl/cache"
	"github.com/sakshamg567/miniurl/db"
	"github.com/sakshamg567/miniurl/services/redirect"
	"github.com/sakshamg567/miniurl/services/shortener"
	"github.com/sakshamg567/miniurl/services/token"
	"github.com/sakshamg567/miniurl/shared/proto/redirectpb"
	"github.com/sakshamg567/miniurl/shared/proto/shortenerpb"
	"github.com/sakshamg567/miniurl/shared/proto/tokenpb"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	dsn := os.Getenv("dsn")
	if dsn == "" {
		log.Fatal("dsn string not set")
	}

	cacheAddr := os.Getenv("cache")
	if cacheAddr == "" {
		log.Fatal("cache addr not found")
	}

	db.Connect(dsn)
	cache.Connect(cacheAddr)

	// TokenService
	go func() {
		lis, _ := net.Listen("tcp", ":50051")
		s := grpc.NewServer()
		tokenService := token.NewTokenService(cacheAddr)
		tokenpb.RegisterTokenServiceServer(s, tokenService)
		log.Println("Token service running on :50051")
		log.Fatal(s.Serve(lis))
	}()

	go func() {
		lis, _ := net.Listen("tcp", ":50052")
		s := grpc.NewServer()
		shortenerService, _ := shortener.NewShortenerService("localhost:50051")
		shortenerpb.RegisterShortenerServiceServer(s, shortenerService)
		log.Println("Shortener service running on :50051")
		log.Fatal(s.Serve(lis))
	}()

	go func() {
		lis, _ := net.Listen("tcp", ":50053")
		s := grpc.NewServer()
		redirectService := redirect.NewRedirectService()
		redirectpb.RegisterRedirectServiceServer(s, redirectService)
		log.Println("Shortener service running on :50053")
		log.Fatal(s.Serve(lis))
	}()

	select {}

}
