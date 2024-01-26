package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/linden-honey/linden-honey-api-go/pkg/application/domain/song"
)

// GrobParser is the implementation of a song parser for the site gr-oborona.ru.
type GrobParser struct {
}

// NewGrobParser returns a pointer to the new instance of [GrobParser].
func NewGrobParser() *GrobParser {
	return &GrobParser{}
}

// ParseSong parses the input html and returns a pointer to the new instance of [song.Entity] or an error.
func (p *GrobParser) ParseSong(input string) (*song.Entity, error) {
	document, err := p.parseHTML(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	title := document.Find("h2").Text()
	tags := p.parseTags(document)
	lyricsHTML, err := document.Find("p:last-of-type").Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get lyrics html: %w", err)
	}

	lyrics, err := p.parseLyrics(lyricsHTML)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lyrics: %w", err)
	}

	return &song.Entity{
		Metadata: song.Metadata{
			Title: title,
			Tags:  tags,
		},
		Lyrics: lyrics,
	}, nil
}

// ParsePreviews parses the input html and returns a slice of [song.Metadata] instances or an error.
func (p *GrobParser) ParsePreviews(input string) ([]song.Metadata, error) {
	document, err := p.parseHTML(input)
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

func (p *GrobParser) parseTags(document *goquery.Document) song.Tags {
	tags := make(song.Tags, 0)

	author := strings.TrimSpace(substringAfterLast(document.Find("p:has(strong:contains(Автор))").Text(), ": "))
	if author != "" {
		tags = append(tags, song.Tag{
			Name:  "author",
			Value: author,
		})
	}

	album := strings.TrimSpace(substringAfterLast(document.Find("p:has(strong:contains(Альбом))").Text(), ": "))
	if album != "" {
		if artist, ok := findKeyByValueInMultiValueMap(grobArtistAlbums, album); ok {
			tags = append(tags, song.Tag{
				Name:  "artist",
				Value: artist,
			})
		}

		if a, ok := findKeyByValueInMultiValueMap(grobAlbumInvalidVariants, album); ok {
			album = a
		}
		tags = append(tags, song.Tag{
			Name:  "album",
			Value: album,
		})
	}

	return tags
}

func (p *GrobParser) parseLyrics(input string) (song.Lyrics, error) {
	verses := make(song.Lyrics, 0)
	// hint: match nbsp; character (\xA0) that not included input \s group
	re := regexp.MustCompile(`(?:<br/>[\s\xA0]*){2,}`)
	for _, verseHTML := range re.Split(input, -1) {
		verse, err := p.parseVerse(verseHTML)
		if err != nil {
			return nil, fmt.Errorf("failed to parse a verse: %w", err)
		}

		verses = append(verses, *verse)
	}

	return verses, nil
}

func (p *GrobParser) parseVerse(input string) (*song.Verse, error) {
	quotes := make([]song.Quote, 0)
	for _, text := range strings.Split(input, "<br/>") {
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

func (p *GrobParser) parseQuote(input string) (*song.Quote, error) {
	document, err := p.parseHTML(input)
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

func (p *GrobParser) parseHTML(input string) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("failed to create a document: %w", err)
	}

	return document, err
}
