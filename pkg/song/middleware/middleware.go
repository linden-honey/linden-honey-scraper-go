package middleware

import (
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song"
)

// Middleware represents the service layer middleware
type Middleware func(song.Service) song.Service

// Compose composes middlewares into a single one
func Compose(mws ...Middleware) Middleware {
	return func(svc song.Service) song.Service {
		for _, mw := range mws {
			svc = mw(svc)
		}

		return svc
	}
}
