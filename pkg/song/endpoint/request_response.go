package endpoint

import (
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
)

type GetSongRequest struct {
	ID string
}

type GetSongResponse struct {
	Result *song.Song
}

type GetSongsRequest struct {
}

type GetSongsResponse struct {
	Results []song.Song
}

type GetPreviewsRequest struct {
}

type GetPreviewsResponse struct {
	Results []song.Preview
}
