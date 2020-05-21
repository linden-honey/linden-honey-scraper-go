package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
)

// songController represents the song controller implementation
type songController struct {
	logger  *log.Logger
	scraper scraper.Scraper
}

// NewSongController returns a pointer to the new instance of songController
func NewSongController(logger *log.Logger, scraper scraper.Scraper) *songController {
	return &songController{
		logger:  logger,
		scraper: scraper,
	}
}

func (c *songController) safeWrite(supplier func() (interface{}, error), w http.ResponseWriter) {
	data, err := supplier()
	if err != nil {
		WriteError(err, http.StatusInternalServerError, w)
	} else {
		err = WriteJSON(data, http.StatusOK, w)
		if err != nil {
			c.logger.Error(err)
			WriteError(err, http.StatusInternalServerError, w)
		}
	}
}

// GetSong handles getting song via http
func (c *songController) GetSong(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	songID := pathVariables["songId"]
	supplier := func() (interface{}, error) {
		return c.scraper.GetSong(songID)
	}
	c.safeWrite(supplier, w)
}

// GetSongs handles getting songs via http
func (c *songController) GetSongs(w http.ResponseWriter, r *http.Request) {
	supplier := func() (interface{}, error) {
		return c.scraper.GetSongs()
	}
	c.safeWrite(supplier, w)
}

// GetPreviews handles getting previews via http
func (c *songController) GetPreviews(w http.ResponseWriter, r *http.Request) {
	supplier := func() (interface{}, error) {
		return c.scraper.GetPreviews()
	}
	c.safeWrite(supplier, w)
}
