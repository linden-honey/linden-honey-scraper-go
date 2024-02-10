package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/linden-honey/linden-honey-api-go/pkg/application/domain/song"
)

var (
	// HINT: `\xA0` match `&nbsp;`(non-breaking space character) that is not included in \s group
	ReOneOrMoreSpace     = regexp.MustCompile(`[\s\xA0]+`)
	ReGrobVerseSeparator = regexp.MustCompile(`(?:<br/>[\s\xA0]*){2,}`)
)

// GrobParser is the implementation of a song parser for the site gr-oborona.ru.
type GrobParser struct {
}

// NewGrobParser returns a pointer to the new instance of [GrobParser].
func NewGrobParser() *GrobParser {
	return &GrobParser{}
}

// ParseSong parses the in html and returns a pointer to the new instance of [song.Entity] or an error.
func (p *GrobParser) ParseSong(in string) (*song.Entity, error) {
	document, err := p.parseHTML(in)
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

// ParsePreviews parses the in html and returns a slice of [song.Metadata] instances or an error.
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

func (p *GrobParser) parseTags(document *goquery.Document) song.Tags {
	tags := make(song.Tags, 0)

	author := strings.TrimSpace(substringAfterLast(document.Find("p:has(strong:contains(Автор))").Text(), ":"))
	if author != "" {
		tags = append(tags, song.Tag{
			Name:  "author",
			Value: author,
		})
	}

	album := strings.TrimSpace(substringAfterLast(document.Find("p:has(strong:contains(Альбом))").Text(), ":"))
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

func (p *GrobParser) parseLyrics(in string) (song.Lyrics, error) {
	out := make(song.Lyrics, 0)
	for _, verseHTML := range ReGrobVerseSeparator.Split(in, -1) {
		verse, err := p.parseVerse(verseHTML)
		if err != nil {
			return nil, fmt.Errorf("failed to parse a verse: %w", err)
		}

		out = append(out, verse)
	}

	return out, nil
}

func (p *GrobParser) parseVerse(in string) (song.Verse, error) {
	out := make(song.Verse, 0)
	for _, quoteHTML := range strings.Split(in, "<br/>") {
		quote, err := p.parseQuote(quoteHTML)
		if err != nil {
			return nil, fmt.Errorf("failed to parse a quote: %w", err)
		}

		out = append(out, quote)
	}

	return out, nil
}

func (p *GrobParser) parseQuote(in string) (song.Quote, error) {
	document, err := p.parseHTML(in)
	if err != nil {
		return "", fmt.Errorf("failed to parse html: %w", err)
	}

	quoteText := document.Text()
	quote := ReOneOrMoreSpace.ReplaceAllLiteralString(strings.TrimSpace(quoteText), " ")

	return song.Quote(quote), nil
}

func (p *GrobParser) parseHTML(in string) (*goquery.Document, error) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(in))
	if err != nil {
		return nil, fmt.Errorf("failed to create a document: %w", err)
	}

	return document, err
}
