package endpoint

import (
	"github.com/linden-honey/linden-honey-scraper-go/pkg/docs"
)

//GetSpecRequest represents a request object
type GetSpecRequest struct {
}

//GetSpecResponse represents a response object
type GetSpecResponse struct {
	Spec *docs.Spec
}
