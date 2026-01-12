package handlers

type Handlers struct {
	User   UserHandler
	Auth   AuthHandler
	Role   RoleHandler
	Policy PolicyHandler
}

func NewHandlers(user UserHandler, auth AuthHandler, role RoleHandler, policy PolicyHandler) *Handlers {
	return &Handlers{
		User:   user,
		Auth:   auth,
		Role:   role,
		Policy: policy,
	}
}
