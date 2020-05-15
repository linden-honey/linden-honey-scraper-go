package parser

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/domain"
)

func substringAfterLast(s, substr string) string {
	if len(s) == 0 {
		return s
	}
	startIndex := strings.LastIndex(s, substr)
	if len(substr) == 0 || startIndex == -1 {
		return ""
	}
	return s[startIndex+len(substr):]
}

func parseHTML(html string) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during document creation")
	}
	return document, err
}

// ParseQuote parses html and returns a quote
func ParseQuote(html string) (*domain.Quote, error) {
	document, err := parseHTML(html)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during quote parsing")
	}
	re := regexp.MustCompile(`\s+`)
	phrase := re.ReplaceAllLiteralString(strings.TrimSpace(document.Text()), " ")
	quote := &domain.Quote{
		Phrase: phrase,
	}
	return quote, nil
}

// ParseVerse parses html and returns a verse
func ParseVerse(html string) (*domain.Verse, error) {
	quotes := make([]*domain.Quote, 0)
	for _, text := range strings.Split(html, "<br/>") {
		quote, err := ParseQuote(text)
		if err != nil {
			log.Println("Error happened during verse parsing", err)
		} else {
			quotes = append(quotes, quote)
		}
	}
	verse := &domain.Verse{
		Quotes: quotes,
	}
	return verse, nil
}

func parseLyrics(html string) []*domain.Verse {
	verses := make([]*domain.Verse, 0)
	re := regexp.MustCompile(`(?:<br/>\s*){2,}`)
	for _, verseHTML := range re.Split(html, -1) {
		verse, err := ParseVerse(verseHTML)
		if err != nil {
			log.Println("Error happened during lyrics parsing", err)
		} else {
			verses = append(verses, verse)
		}
	}
	return verses
}

// ParseSong parses html and returns a song
func ParseSong(html string) (*domain.Song, error) {
	document, err := parseHTML(html)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during html parsing")
	}
	title := document.Find("h2").Text()
	author := substringAfterLast(document.Find("p:has(strong:contains(Автор))").Text(), ": ")
	album := substringAfterLast(document.Find("p:has(strong:contains(Альбом))").Text(), ": ")
	lyricsHTML, _ := document.Find("p:last-of-type").Html()
	verses := parseLyrics(lyricsHTML)
	song := &domain.Song{
		Title:  title,
		Author: author,
		Album:  album,
		Verses: verses,
	}
	return song, nil
}

// ParsePreviews parses html and returns a song
func ParsePreviews(html string) ([]*domain.Preview, error) {
	document, err := parseHTML(html)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during html parsing")
	}
	previews := make([]*domain.Preview, 0)
	document.Find("#abc_list a").Each(func(_ int, link *goquery.Selection) {
		path, pathExists := link.Attr("href")
		if pathExists {
			startIndex := strings.LastIndex(path, "/")
			endIndex := strings.Index(path, ".")
			if startIndex != -1 && endIndex != -1 {
				id := path[startIndex+1 : endIndex]
				title := link.Text()
				previews = append(previews, &domain.Preview{
					ID:    id,
					Title: title,
				})
			}
		}
	})
	return previews, nil
}
