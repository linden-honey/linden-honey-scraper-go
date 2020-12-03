package endpoint

import (
	"github.com/linden-honey/linden-honey-scraper-go/pkg/scraper/domain"
)

type GetSongRequest struct {
	ID string
}

type GetSongResponse struct {
	Result *domain.Song
}

type GetSongsRequest struct {
}

type GetSongsResponse struct {
	Results []domain.Song
}

type GetPreviewsRequest struct {
}

type GetPreviewsResponse struct {
	Results []domain.Preview
}
