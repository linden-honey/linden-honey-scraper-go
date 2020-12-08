package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service"
)

type Endpoints struct {
	GetSong     endpoint.Endpoint
	GetSongs    endpoint.Endpoint
	GetPreviews endpoint.Endpoint
}

func NewEndpoints(svc service.Service) *Endpoints {
	return &Endpoints{
		GetSong:     makeGetSongEndpoint(svc),
		GetSongs:    makeGetSongsEndpoint(svc),
		GetPreviews: makeGetPreviewsEndpoint(svc),
	}
}

func makeGetSongEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetSongRequest)

		song, err := svc.GetSong(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return GetSongResponse{
			Result: song,
		}, nil
	}
}

func makeGetSongsEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(GetSongsRequest)

		songs, err := svc.GetSongs(ctx)
		if err != nil {
			return nil, err
		}

		return GetSongsResponse{
			Results: songs,
		}, nil
	}
}

func makeGetPreviewsEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(GetPreviewsRequest)

		previews, err := svc.GetPreviews(ctx)
		if err != nil {
			return nil, err
		}

		return GetPreviewsResponse{
			Results: previews,
		}, nil
	}
}
