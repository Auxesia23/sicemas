package handlers

type Handlers struct {
	User             UserHandler
	Auth             AuthHandler
	Role             RoleHandler
	Policy           PolicyHandler
	JenisSitus       JenisSitusHandler
	SitusKeagamaan   SitusKeagamaanHandler
	DashboardHandler DashboardHandler
	ActivityHandler  ActivityHandler
}

func NewHandlers(user UserHandler,
	auth AuthHandler,
	role RoleHandler,
	policy PolicyHandler,
	jenisSitus JenisSitusHandler,
	situsKeagamaan SitusKeagamaanHandler,
	dashboardHandler DashboardHandler,
	activityHandler ActivityHandler) *Handlers {
	return &Handlers{
		User:             user,
		Auth:             auth,
		Role:             role,
		Policy:           policy,
		JenisSitus:       jenisSitus,
		SitusKeagamaan:   situsKeagamaan,
		DashboardHandler: dashboardHandler,
		ActivityHandler:  activityHandler,
	}
}
