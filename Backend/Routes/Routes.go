package routes

import (
	"net/http"
	handler "social-network/Handler"
	middleware "social-network/Middleware"
)

func Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", handler.Redirect)

	mux.HandleFunc("POST /login", handler.Login)
	mux.HandleFunc("POST /register", middleware.RegisterMiddleware(handler.Register))

	mux.HandleFunc("/home", handler.Home)
}
