package parser

import (
	"regexp"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/domain"

	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func substringAfterLast(s, substr string) string {
	const EMPTY = ""
	if len(s) == 0 {
		return s
	}
	startIndex := strings.LastIndex(s, substr)
	if len(substr) == 0 || startIndex == -1 {
		return EMPTY
	}
	return s[startIndex+len(substr):]
}

func parseHTML(html string) *goquery.Document {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal("Error happend during html parsing", err)
	}
	return document
}

// ParseQuote parses html and returns a quote
func ParseQuote(html string) domain.Quote {
	document := parseHTML(html)
	re := regexp.MustCompile(`\s+`)
	phrase := re.ReplaceAllLiteralString(strings.TrimSpace(document.Text()), " ")
	return domain.Quote{
		Phrase: phrase,
	}
}

// ParseVerse parses html and returns a verse
func ParseVerse(html string) domain.Verse {
	var quotes []domain.Quote
	for _, text := range strings.Split(html, "<br/>") {
		quote := ParseQuote(text)
		if len(quote.Phrase) != 0 {
			quotes = append(quotes, quote)
		}

	}
	return domain.Verse{
		Quotes: quotes,
	}
}

func parseLyrics(html string) []domain.Verse {
	verses := make([]domain.Verse, 0)
	re := regexp.MustCompile(`(?:<br/>\s*){2,}`)
	for _, verseHTML := range re.Split(html, -1) {
		verse := ParseVerse(verseHTML)
		verses = append(verses, verse)
	}
	return verses
}

// ParseSong parses html and returns a song
func ParseSong(html string) domain.Song {
	document := parseHTML(html)
	title := document.Find("h2").Text()
	author := substringAfterLast(document.Find("p:has(strong:contains(Автор))").Text(), ": ")
	album := substringAfterLast(document.Find("p:has(strong:contains(Альбом))").Text(), ": ")
	lyricsHTML, _ := document.Find("p:last-of-type").Html()
	verses := parseLyrics(lyricsHTML)
	return domain.Song{
		Title:  title,
		Author: author,
		Album:  album,
		Verses: verses,
	}
}

// ParsePreviews parses html and returns a song
func ParsePreviews(html string) []domain.Preview {
	document := parseHTML(html)
	previews := make([]domain.Preview, 0)
	document.Find("#abc_list a").Each(func(_ int, link *goquery.Selection) {
		path, pathExists := link.Attr("href")
		if pathExists {
			startIndex := strings.LastIndex(path, "/")
			endIndex := strings.Index(path, ".")
			if startIndex != -1 && endIndex != -1 {
				id := path[startIndex+1 : endIndex]
				title := link.Text()
				previews = append(previews, domain.Preview{
					ID:    id,
					Title: title,
				})
			}
		}
	})
	return previews
}
