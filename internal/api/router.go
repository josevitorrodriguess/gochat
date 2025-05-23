package api

func (api *Api) RegisterRoutes() {
	api.Echo.POST("/signup", api.SignUpHandler)

}
