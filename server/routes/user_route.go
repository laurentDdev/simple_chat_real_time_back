package routes

import (
	"instantmsg/models"
	"instantmsg/server"
	"instantmsg/server/controllers"
)

type UserRoute struct {
	Routes         []server.Route
	UserController *controllers.UserController
	Context        *models.AppContext
}

func NewUserRoute(ctx *models.AppContext) *UserRoute {

	return &UserRoute{
		Routes:         []server.Route{},
		UserController: controllers.NewUserController(ctx),
	}
}

func (ur *UserRoute) GetRoutes() []server.Route {
	routes := []server.Route{
		{
			Path:   "/login",
			Handle: ur.UserController.Login,
		},
		{
			Path:   "/register",
			Handle: ur.UserController.Register,
		},
	}

	return routes
}
