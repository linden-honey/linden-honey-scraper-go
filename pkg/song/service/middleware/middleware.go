package middleware

import (
	"github.com/linden-honey/linden-honey-scraper-go/pkg/song/service"
)

type Middleware func(service.Service) service.Service
