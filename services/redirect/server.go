package redirect

import (
	"context"
	"fmt"
	"time"

	"github.com/sakshamg567/miniurl/cache"
	"github.com/sakshamg567/miniurl/db"
	redirectpb "github.com/sakshamg567/miniurl/shared/proto/redirectpb"
)

type redirectService struct {
	redirectpb.UnimplementedRedirectServiceServer
}

func NewRedirectService() redirectpb.RedirectServiceServer {
	return &redirectService{}
}

func (s *redirectService) GetLongURL(ctx context.Context, req *redirectpb.GetLongURLRequest) (*redirectpb.GetLongURLResponse, error) {

	longUrl, err := cache.Get(req.ShortCode)
	if err == nil {
		return &redirectpb.GetLongURLResponse{LongUrl: longUrl}, nil
	}

	longUrl, err = db.GetLongURL(req.ShortCode)
	if err != nil {
		return nil, fmt.Errorf("URL Not found :%v", err)
	}

	_ = cache.Set(req.ShortCode, longUrl, 24*time.Hour)
	return &redirectpb.GetLongURLResponse{LongUrl: longUrl}, nil

	// 5. log in prometheus *

}
