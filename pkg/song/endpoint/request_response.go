package endpoint

import (
	"github.com/linden-honey/linden-honey-go/pkg/song"
)

//GetSongRequest represents a request object
type GetSongRequest struct {
	ID string
}

//GetSongResponse represents a response object
type GetSongResponse struct {
	Result *song.Song
}

//GetSongsRequest represents a request object
type GetSongsRequest struct {
}

//GetSongsResponse represents a response object
type GetSongsResponse struct {
	Results []song.Song
}

//GetPreviewsRequest represents a request object
type GetPreviewsRequest struct {
}

//GetPreviewsResponse represents a response object
type GetPreviewsResponse struct {
	Results []song.Preview
}
