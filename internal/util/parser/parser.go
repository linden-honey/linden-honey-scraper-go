package parser

import (
	"github.com/linden-honey/linden-honey-scraper-go/pkg/domain"

	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseHTML(html string) *goquery.Document {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal("Error happend during html parsing", err)
	}
	return document
}

// ParseQuote parses html and returns a quote
func ParseQuote(html string) domain.Quote {
	return domain.Quote{}
}

// ParseVerse parses html and returns a verse
func ParseVerse(html string) domain.Verse {
	return domain.Verse{}
}

// ParseSong parses html and returns a song
func ParseSong(html string) domain.Song {
	return domain.Song{}
}

// ParsePreviews parses html and returns a song
func ParsePreviews(html string) (previews []domain.Preview) {
	document := parseHTML(html)
	document.Find("#abc_list a").Each(func(_ int, link *goquery.Selection) {
		path, pathExists := link.Attr("href")
		if pathExists {
			startIndex := strings.LastIndex(path, "/")
			endIndex := strings.Index(path, ".")
			if startIndex != -1 && endIndex != -1 {
				title := link.Text()
				id := path[startIndex+1 : endIndex]
				previews = append(previews, domain.Preview{
					ID:    id,
					Title: title,
				})
			}
		}
	})
	return previews
}
