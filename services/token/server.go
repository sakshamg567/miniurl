package token

import (
	"context"

	"github.com/redis/go-redis/v9"
	tokenpb "github.com/sakshamg567/miniurl/shared/proto/tokenpb"
)

type TokenService struct {
	db *redis.Client
	tokenpb.UnimplementedTokenServiceServer
}

func NewTokenService(cacheAddr string) *TokenService {
	rdb := redis.NewClient(&redis.Options{Addr: cacheAddr})

	ctx := context.Background()
	const startVal = int64(100000000000 - 1)

	ok, err := rdb.Exists(ctx, "url_counter").Result()
	if err != nil {
		panic(err)
	}
	if ok == 0 {
		if err := rdb.Set(ctx, "url_counter", startVal, 0).Err(); err != nil {
			panic(err)
		}
	}

	return &TokenService{
		db: rdb,
	}
}

func (s *TokenService) NextID(ctx context.Context, req *tokenpb.NextIDReq) (*tokenpb.NextIDRes, error) {
	defer s.db.Close()
	id, err := s.db.Incr(ctx, "url_counter").Result()
	if err != nil {
		return nil, err
	}
	return &tokenpb.NextIDRes{Id: uint64(id)}, nil
}
