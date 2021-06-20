package scraper

// Middleware represents the service layer middleware
type Middleware func(Service) Service
