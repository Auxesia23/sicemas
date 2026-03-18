package middlewares

type Middlewares struct {
	Auth    AuthMiddleware
	Limiter RateLimiter
}

func NewMiddlewares(auth AuthMiddleware, limiter *RateLimiter) *Middlewares {
	return &Middlewares{
		Auth:    auth,
		Limiter: *limiter,
	}
}
