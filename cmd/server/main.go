package main

import (
	"fmt"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/service/scraper"
)

func main() {
	s := scraper.Create(&scraper.Properties{
		BaseURL: "http://www.gr-oborona.ru",
	})
	fmt.Println(s.FetchPreviews())
}
