package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
)

// GrobParser represents the parser implementation for gr-oborona.ru
type GrobParser struct {
}

// NewGrobParser returns a pointer to the new instance of GrobParser
func NewGrobParser() (*GrobParser, error) {
	return &GrobParser{}, nil
}

// parseHTML parse html and returns a pointer to the Document instance
func (p *GrobParser) parseHTML(html string) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return document, err
}

// ParseQuote parses html and returns a pointer to the Quote instance
func (p *GrobParser) ParseQuote(html string) (*song.Quote, error) {
	document, err := p.parseHTML(html)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	re := regexp.MustCompile(`\s+`)
	phrase := re.ReplaceAllLiteralString(strings.TrimSpace(document.Text()), " ")
	quote := &song.Quote{
		Phrase: phrase,
	}

	return quote, nil
}

// ParseVerse parses html and returns a pointer to the Verse instance
func (p *GrobParser) ParseVerse(html string) (*song.Verse, error) {
	quotes := make([]song.Quote, 0)
	for _, text := range strings.Split(html, "<br/>") {
		quote, err := p.ParseQuote(text)
		if err != nil {
			return nil, fmt.Errorf("failed to parse quote: %w", err)
		}
		quotes = append(quotes, *quote)
	}

	return &song.Verse{
		Quotes: quotes,
	}, nil
}

// ParseVerse parses html and returns a slice of pointers of the Verse instances
func (p *GrobParser) parseLyrics(html string) ([]song.Verse, error) {
	verses := make([]song.Verse, 0)
	re := regexp.MustCompile(`(?:<br/>\s*){2,}`)
	for _, verseHTML := range re.Split(html, -1) {
		verse, err := p.ParseVerse(verseHTML)
		if err != nil {
			return nil, fmt.Errorf("failed to parse verse: %w", err)
		}

		verses = append(verses, *verse)
	}

	return verses, nil
}

// ParseSong parses html and returns a pointer to the Song instance
func (p *GrobParser) ParseSong(html string) (*song.Song, error) {
	document, err := p.parseHTML(html)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	title := document.Find("h2").Text()
	author := substringAfterLast(document.Find("p:has(strong:contains(Автор))").Text(), ": ")
	album := substringAfterLast(document.Find("p:has(strong:contains(Альбом))").Text(), ": ")
	lyricsHTML, _ := document.Find("p:last-of-type").Html()

	verses, err := p.parseLyrics(lyricsHTML)
	if err != nil {

	}

	return &song.Song{
		Title:  title,
		Author: author,
		Album:  album,
		Verses: verses,
	}, nil
}

// ParsePreviews parses html and returns a slice of pointers of the Preview instances
func (p *GrobParser) ParsePreviews(html string) ([]song.Preview, error) {
	document, err := p.parseHTML(html)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	previews := make([]song.Preview, 0)
	document.Find("#abc_list a").Each(func(_ int, link *goquery.Selection) {
		path, pathExists := link.Attr("href")
		if pathExists {
			startIndex := strings.LastIndex(path, "/")
			endIndex := strings.Index(path, ".")
			if startIndex != -1 && endIndex != -1 {
				id := path[startIndex+1 : endIndex]
				title := link.Text()
				previews = append(previews, song.Preview{
					ID:    id,
					Title: title,
				})
			}
		}
	})

	return previews, nil
}
