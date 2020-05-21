package parser

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/domain"
)

// substringAfterLast util function to get last substring after some inclusion
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

// Fetcher represents the fetcher interface
type Parser interface {
	ParseSong(html string) (*domain.Song, error)
	ParsePreviews(html string) ([]*domain.Preview, error)
}

// defaultParser represents the default parser implementation
type defaultParser struct {
	logger *log.Logger
}

// NewDefaultParser returns a pointer to the new instance of defaultParser
func NewDefaultParser(logger *log.Logger) *defaultParser {
	return &defaultParser{
		logger: logger,
	}
}

// parseHTML parse html and returns a pointer to the Document instance
func (p *defaultParser) parseHTML(html string) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during document creation")
	}
	return document, err
}

// ParseQuote parses html and returns a pointer to the Quote instance
func (p *defaultParser) ParseQuote(html string) (*domain.Quote, error) {
	document, err := p.parseHTML(html)
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

// ParseVerse parses html and returns a pointer to the Verse instance
func (p *defaultParser) ParseVerse(html string) (*domain.Verse, error) {
	quotes := make([]*domain.Quote, 0)
	for _, text := range strings.Split(html, "<br/>") {
		quote, err := p.ParseQuote(text)
		if err != nil {
			p.logger.Println("Error happened during quote parsing", err)
		} else {
			quotes = append(quotes, quote)
		}
	}
	verse := &domain.Verse{
		Quotes: quotes,
	}
	return verse, nil
}

// ParseVerse parses html and returns a slice of pointers of the Verse instances
func (p *defaultParser) parseLyrics(html string) []*domain.Verse {
	verses := make([]*domain.Verse, 0)
	re := regexp.MustCompile(`(?:<br/>\s*){2,}`)
	for _, verseHTML := range re.Split(html, -1) {
		verse, err := p.ParseVerse(verseHTML)
		if err != nil {
			p.logger.Println("Error happened during lyrics parsing", err)
		} else {
			verses = append(verses, verse)
		}
	}
	return verses
}

// ParseSong parses html and returns a pointer to the Song instance
func (p *defaultParser) ParseSong(html string) (*domain.Song, error) {
	document, err := p.parseHTML(html)
	if err != nil {
		return nil, errors.Wrap(err, "Error happened during html parsing")
	}
	title := document.Find("h2").Text()
	author := substringAfterLast(document.Find("p:has(strong:contains(Автор))").Text(), ": ")
	album := substringAfterLast(document.Find("p:has(strong:contains(Альбом))").Text(), ": ")
	lyricsHTML, _ := document.Find("p:last-of-type").Html()
	verses := p.parseLyrics(lyricsHTML)
	song := &domain.Song{
		Title:  title,
		Author: author,
		Album:  album,
		Verses: verses,
	}
	return song, nil
}

// ParsePreviews parses html and returns a slice of pointers of the Preview instances
func (p *defaultParser) ParsePreviews(html string) ([]*domain.Preview, error) {
	document, err := p.parseHTML(html)
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
