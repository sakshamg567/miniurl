package shortener

import (
	"context"
	"fmt"

	"github.com/mattheath/base62"
	"github.com/sakshamg567/miniurl/db"
	shortenerpb "github.com/sakshamg567/miniurl/shared/proto/shortenerpb"
	tokenpb "github.com/sakshamg567/miniurl/shared/proto/tokenpb"
	"google.golang.org/grpc"
)

var MAX_CUSTOM_SHORTCODE_LEN = 16

type shortenerService struct {
	tokenClient tokenpb.TokenServiceClient
	shortenerpb.UnimplementedShortenerServiceServer
}

func NewShortenerService(tokenServiceAddr string) (shortenerpb.ShortenerServiceServer, error) {
	client, err := newTokenClient(tokenServiceAddr)
	if err != nil {
		return nil, err
	}

	return &shortenerService{
		tokenClient: client,
	}, nil
}

func newTokenClient(addr string) (tokenpb.TokenServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return tokenpb.NewTokenServiceClient(conn), nil
}

func (s *shortenerService) CreateShortURL(ctx context.Context, req *shortenerpb.CreateShortURLRequest) (*shortenerpb.CreateShortURLResponse, error) {
	long_url := req.LongUrl
	// custom_short_code := req.CustomShortCode

	idRes, err := s.tokenClient.NextID(ctx, &tokenpb.NextIDReq{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate id : %v", err)
	}

	shortCode := base62.EncodeInt64(int64(idRes.Id))

	if err := db.InsertLongURL(long_url, shortCode); err != nil {
		return nil, fmt.Errorf("failed to insert url into postgres :%v", err)
	}

	return &shortenerpb.CreateShortURLResponse{ShortCode: shortCode}, nil
}
