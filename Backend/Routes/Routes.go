package routes

import (
	"net/http"

	handler "social-network/Handler"
	middleware "social-network/Middleware"
)

/*
This function takes 1 argument:
  - an *http.ServeMux who is used to create the server

The purpose of this function is to create all the server endpoints and define the function call for each endpoint.

The function have no return
*/
func Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", handler.Redirect)

	mux.HandleFunc("POST /login", handler.Login)
	mux.HandleFunc("POST /register", middleware.RegisterMiddleware(handler.Register))

	mux.HandleFunc("/home", handler.Home)
}
