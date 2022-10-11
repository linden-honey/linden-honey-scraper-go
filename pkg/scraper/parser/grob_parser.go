package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/linden-honey/linden-honey-api-go/pkg/song"
)

// GrobParser represents the parser implementation for gr-oborona.ru
type GrobParser struct {
}

// NewGrobParser returns a pointer to the new instance of GrobParser or an error
func NewGrobParser() *GrobParser {
	return &GrobParser{}
}

func (p *GrobParser) parseHTML(in string) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(in))
	if err != nil {
		return nil, fmt.Errorf("failed to create a document: %w", err)
	}

	return document, err
}

func (p *GrobParser) parseQuote(in string) (*song.Quote, error) {
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

func (p *GrobParser) parseVerse(in string) (*song.Verse, error) {
	quotes := make([]song.Quote, 0)
	for _, text := range strings.Split(in, "<br/>") {
		quote, err := p.parseQuote(text)
		if err != nil {
			return nil, fmt.Errorf("failed to parse a quote: %w", err)
		}
		quotes = append(quotes, *quote)
	}

	return &song.Verse{
		Quotes: quotes,
	}, nil
}

func (p *GrobParser) parseLyrics(in string) (song.Lyrics, error) {
	verses := make(song.Lyrics, 0)
	// hint: match nbsp; character (\xA0) that not included in \s group
	re := regexp.MustCompile(`(?:<br/>[\s\xA0]*){2,}`)
	for _, verseHTML := range re.Split(in, -1) {
		verse, err := p.parseVerse(verseHTML)
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

	title := document.Find("h2").Text()
	author := substringAfterLast(document.Find("p:has(strong:contains(Автор))").Text(), ": ")
	album := substringAfterLast(document.Find("p:has(strong:contains(Альбом))").Text(), ": ")
	lyricsHTML, err := document.Find("p:last-of-type").Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get lyrics html: %w", err)
	}

	lyrics, err := p.parseLyrics(lyricsHTML)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lyrics: %w", err)
	}

	return &song.Song{
		Metadata: song.Metadata{
			Title: title,
			Tags: song.Tags{
				{
					Name:  "author",
					Value: author,
				},
				{
					Name:  "album",
					Value: album,
				},
			},
		},
		Lyrics: lyrics,
	}, nil
}

// ParsePreviews parses html and returns a slice of pointers of the Preview instances
func (p *GrobParser) ParsePreviews(in string) ([]song.Metadata, error) {
	document, err := p.parseHTML(in)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	previews := make([]song.Metadata, 0)
	document.Find("#abc_list a").Each(func(_ int, link *goquery.Selection) {
		path, pathExists := link.Attr("href")
		if pathExists {
			startIndex := strings.LastIndex(path, "/")
			endIndex := strings.Index(path, ".")
			if startIndex != -1 && endIndex != -1 {
				id := path[startIndex+1 : endIndex]
				title := link.Text()
				previews = append(previews, song.Metadata{
					ID:    id,
					Title: title,
				})
			}
		}
	})

	return previews, nil
}
