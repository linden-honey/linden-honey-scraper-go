package main

import (
	"encoding/json"
	"net/http"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
)

func main() {
	s := scraper.Create(&scraper.Properties{
		BaseURL: "http://www.gr-oborona.ru",
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		songs := s.FetchSongs()
		bytes, _ := json.Marshal(songs)
		w.Write(bytes)
		w.WriteHeader(200)
	})
	http.ListenAndServe(":8080", nil)
}
