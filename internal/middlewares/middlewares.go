package middlewares

type Middlewares struct {
	Auth AuthMiddleware
}

func NewMiddlewares(auth AuthMiddleware) *Middlewares {
	return &Middlewares{
		Auth: auth,
	}
}
