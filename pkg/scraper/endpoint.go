package scraper

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints represents endpoints definition
type Endpoints struct {
	GetSong     endpoint.Endpoint
	GetSongs    endpoint.Endpoint
	GetPreviews endpoint.Endpoint
}

func MakeGetSongEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetSongRequest)

		s, err := svc.GetSong(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return GetSongResponse{
			Result: s,
		}, nil
	}
}

func MakeGetSongsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(GetSongsRequest)

		ss, err := svc.GetSongs(ctx)
		if err != nil {
			return nil, err
		}

		return GetSongsResponse{
			Results: ss,
		}, nil
	}
}

func MakeGetPreviewsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(GetPreviewsRequest)

		pp, err := svc.GetPreviews(ctx)
		if err != nil {
			return nil, err
		}

		return GetPreviewsResponse{
			Results: pp,
		}, nil
	}
}
