package parser

import (
	"errors"

	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service/scraper"
)

// NewParser factory function that returns scraper.Parser instance by id
func NewParser(id string) (scraper.Parser, error) {
	switch id {
	case GrobParserID:
		return NewGrobParser()
	default:
		return nil, errors.New("unknown parser")
	}
}
