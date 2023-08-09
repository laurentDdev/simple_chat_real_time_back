package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type routerInterface interface {
	AddRoutes(route []Route)
	RegisterRoutes()
}

type Router struct {
	Router *mux.Router
	Routes []Route
}

func NewRouter() *Router {
	fmt.Println("Création du router")
	return &Router{
		Router: mux.NewRouter(),
		Routes: []Route{},
	}
}

func (r *Router) AddRoutes(routes []Route) {
	for _, route := range routes {
		r.Routes = append(r.Routes, route)
		fmt.Printf("Ajout des routes au router %v\n", route.Path)
	}
}

func (r *Router) RegisterRoutes() {
	for _, route := range r.Routes {
		handle := route.Handle
		if route.MiddleWare != nil {
			handle = route.MiddleWare(http.HandlerFunc(route.Handle)).ServeHTTP
		}
		fmt.Printf("Route[path=%v]\n", route.Path)
		r.Router.HandleFunc(route.Path, handle).Methods("POST")
	}
}

type Route struct {
	Path       string
	Handle     http.HandlerFunc
	MiddleWare func(next http.Handler) http.Handler
}
