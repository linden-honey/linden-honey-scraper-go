package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/linden-honey/linden-honey-go/pkg/song"
)

// GrobParser represents the parser implementation for gr-oborona.ru
type GrobParser struct {
}

// NewGrobParser returns a pointer to the new instance of GrobParser or an error
func NewGrobParser() *GrobParser {
	return &GrobParser{}
}

// parseHTML parse html and returns a pointer to the Document instance
func (p *GrobParser) parseHTML(in string) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(in))
	if err != nil {
		return nil, fmt.Errorf("failed to create a document: %w", err)
	}

	return document, err
}

// ParseQuote parses html and returns a pointer to the Quote instance
func (p *GrobParser) ParseQuote(in string) (*song.Quote, error) {
	document, err := p.parseHTML(in)
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
func (p *GrobParser) ParseVerse(in string) (*song.Verse, error) {
	quotes := make([]song.Quote, 0)
	for _, text := range strings.Split(in, "<br/>") {
		quote, err := p.ParseQuote(text)
		if err != nil {
			return nil, fmt.Errorf("failed to parse a quote: %w", err)
		}
		quotes = append(quotes, *quote)
	}

	return &song.Verse{
		Quotes: quotes,
	}, nil
}

// parseVerses parses html and returns a slice of pointers of the Verse instances
func (p *GrobParser) parseVerses(in string) ([]song.Verse, error) {
	verses := make([]song.Verse, 0)
	// hint: match nbsp; character (\xA0) that not included in \s group
	re := regexp.MustCompile(`(?:<br/>[\s\xA0]*){2,}`)
	for _, verseHTML := range re.Split(in, -1) {
		verse, err := p.ParseVerse(verseHTML)
		if err != nil {
			return nil, fmt.Errorf("failed to parse a verse: %w", err)
		}

		verses = append(verses, *verse)
	}

	return verses, nil
}

// ParseSong parses html and returns a pointer to the Song instance
func (p *GrobParser) ParseSong(in string) (*song.Song, error) {
	document, err := p.parseHTML(in)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	id := document.Url.Query().Get("id")
	title := document.Find("h2").Text()
	author := substringAfterLast(document.Find("p:has(strong:contains(Автор))").Text(), ": ")
	album := substringAfterLast(document.Find("p:has(strong:contains(Альбом))").Text(), ": ")
	versesHTML, err := document.Find("p:last-of-type").Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get lyrics html: %w", err)
	}

	verses, err := p.parseVerses(versesHTML)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lyrics: %w", err)
	}

	return &song.Song{
		Meta: song.Meta{
			ID:     id,
			Title:  title,
			Author: author,
			Album:  album,
		},
		Verses: verses,
	}, nil
}

// ParsePreviews parses html and returns a slice of pointers of the Preview instances
func (p *GrobParser) ParsePreviews(in string) ([]song.Meta, error) {
	document, err := p.parseHTML(in)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	previews := make([]song.Meta, 0)
	document.Find("#abc_list a").Each(func(_ int, link *goquery.Selection) {
		path, pathExists := link.Attr("href")
		if pathExists {
			startIndex := strings.LastIndex(path, "/")
			endIndex := strings.Index(path, ".")
			if startIndex != -1 && endIndex != -1 {
				id := path[startIndex+1 : endIndex]
				title := link.Text()
				previews = append(previews, song.Meta{
					ID:    id,
					Title: title,
				})
			}
		}
	})

	return previews, nil
}
