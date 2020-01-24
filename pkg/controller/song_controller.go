package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
)

type SongController struct {
	Scraper scraper.Scraper
}

func writeJSON(data interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		err := errors.Wrap(err, "Error happend during data marshalling")
		log.Println(err)
		w.WriteHeader(500)
		w.Write(([]byte)(err.Error()))
	} else {
		w.WriteHeader(200)
		w.Write(jsonBytes)
	}
}

func (s *SongController) GetSongs(w http.ResponseWriter, r *http.Request) {
	songs := s.Scraper.FetchSongs()
	writeJSON(songs, w, r)
}

func (s *SongController) GetSong(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	songID := pathVariables["songId"]
	song := s.Scraper.FetchSong(songID)
	if song != nil {
		writeJSON(song, w, r)
	} else {
		w.WriteHeader(404)
		w.Write(([]byte)("Song not found"))
	}
}

func (s *SongController) GetPreviews(w http.ResponseWriter, r *http.Request) {
	previews := s.Scraper.FetchPreviews()
	writeJSON(previews, w, r)
}
