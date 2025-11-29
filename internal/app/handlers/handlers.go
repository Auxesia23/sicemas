package handlers

type Handlers struct {
	User UserHandler
	Auth AuthHandler
}

func NewHandlers(user UserHandler, auth AuthHandler) *Handlers {
	return &Handlers{
		User: user,
		Auth: auth,
	}
}
