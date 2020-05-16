package controller

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
)

type SongController struct {
	Scraper scraper.Scraper
}

func NewSongController(scraper scraper.Scraper) *SongController {
	return &SongController{
		Scraper: scraper,
	}
}

func (c *SongController) GetSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := c.Scraper.FetchSongs()
	if err != nil {
		WriteError(err, http.StatusInternalServerError, w)
	} else {
		WriteJSON(songs, http.StatusOK, w)
	}
}

func (c *SongController) GetSong(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	songID := pathVariables["songId"]
	song, err := c.Scraper.FetchSong(songID)
	if err != nil {
		WriteError(err, http.StatusInternalServerError, w)
	} else {
		WriteJSON(song, http.StatusOK, w)
	}
}

func (c *SongController) GetPreviews(w http.ResponseWriter, r *http.Request) {
	previews, err := c.Scraper.FetchPreviews()
	if err != nil {
		WriteError(err, http.StatusInternalServerError, w)
	} else {
		WriteJSON(previews, http.StatusOK, w)
	}
}
