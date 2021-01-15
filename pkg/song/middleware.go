package song

// Middleware represents the service layer middleware
type Middleware func(Service) Service

// Compose composes middlewares into a single one
func Compose(mws ...Middleware) Middleware {
	return func(svc Service) Service {
		for _, mw := range mws {
			svc = mw(svc)
		}

		return svc
	}
}
