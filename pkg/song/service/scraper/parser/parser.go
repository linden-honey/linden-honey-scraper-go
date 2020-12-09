package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/domain"
)

// Parser represents the parser implementation
type Parser struct {
}

// NewParser returns a pointer to the new instance of defaultParser
func NewParser() *Parser {
	return &Parser{}
}

// parseHTML parse html and returns a pointer to the Document instance
func (p *Parser) parseHTML(html string) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return document, err
}

// ParseQuote parses html and returns a pointer to the Quote instance
func (p *Parser) ParseQuote(html string) (*domain.Quote, error) {
	document, err := p.parseHTML(html)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	re := regexp.MustCompile(`\s+`)
	phrase := re.ReplaceAllLiteralString(strings.TrimSpace(document.Text()), " ")
	quote := &domain.Quote{
		Phrase: phrase,
	}

	return quote, nil
}

// ParseVerse parses html and returns a pointer to the Verse instance
func (p *Parser) ParseVerse(html string) (*domain.Verse, error) {
	quotes := make([]domain.Quote, 0)
	for _, text := range strings.Split(html, "<br/>") {
		quote, err := p.ParseQuote(text)
		if err != nil {
			return nil, fmt.Errorf("failed to parse quote: %w", err)
		}
		quotes = append(quotes, *quote)
	}

	return &domain.Verse{
		Quotes: quotes,
	}, nil
}

// ParseVerse parses html and returns a slice of pointers of the Verse instances
func (p *Parser) parseLyrics(html string) ([]domain.Verse, error) {
	verses := make([]domain.Verse, 0)
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
func (p *Parser) ParseSong(html string) (*domain.Song, error) {
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

	return &domain.Song{
		Title:  title,
		Author: author,
		Album:  album,
		Verses: verses,
	}, nil
}

// ParsePreviews parses html and returns a slice of pointers of the Preview instances
func (p *Parser) ParsePreviews(html string) ([]domain.Preview, error) {
	document, err := p.parseHTML(html)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

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

	return previews, nil
}
