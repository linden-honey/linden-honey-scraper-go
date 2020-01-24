package controller

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
)

type SongController struct {
	Scraper scraper.Scraper
}

func (c *SongController) GetSongs(w http.ResponseWriter, r *http.Request) {
	songs := c.Scraper.FetchSongs()
	writeJSON(songs, w, r)
}

func (c *SongController) GetSong(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	songID := pathVariables["songId"]
	song := c.Scraper.FetchSong(songID)
	if song != nil {
		writeJSON(song, w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write(([]byte)("Song not found"))
	}
}

func (c *SongController) GetPreviews(w http.ResponseWriter, r *http.Request) {
	previews := c.Scraper.FetchPreviews()
	writeJSON(previews, w, r)
}
