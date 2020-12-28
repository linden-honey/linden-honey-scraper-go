package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
)

// Endpoints represents song endpoints
type Endpoints struct {
	GetSong     endpoint.Endpoint
	GetSongs    endpoint.Endpoint
	GetPreviews endpoint.Endpoint
}

// NewEndpoints returns a pointer to the new instance of Endpoints
func NewEndpoints(svc song.Service) (*Endpoints, error) {
	return &Endpoints{
		GetSong:     makeGetSongEndpoint(svc),
		GetSongs:    makeGetSongsEndpoint(svc),
		GetPreviews: makeGetPreviewsEndpoint(svc),
	}, nil
}

func makeGetSongEndpoint(svc song.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
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

func makeGetSongsEndpoint(svc song.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
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

func makeGetPreviewsEndpoint(svc song.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
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
