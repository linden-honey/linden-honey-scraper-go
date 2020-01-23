package controller

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
)

type SongController struct {
	Scraper *scraper.Scraper
}

func (s *SongController) GetSongs(w http.ResponseWriter, r *http.Request) {
	w.Write(([]byte)("SONGS"))
}

func (s *SongController) GetSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Write(([]byte)(vars["songId"]))
}

func (s *SongController) GetPreviews(w http.ResponseWriter, r *http.Request) {
	w.Write(([]byte)("PREVIEWS"))
}
