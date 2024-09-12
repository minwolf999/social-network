package routes

import (
	"fmt"
	"net/http"
	handler "social-network/Handler"
	middleware "social-network/Middleware"
	utils "social-network/Utils"
)

func Routes(mux *http.ServeMux) {

	db, err := utils.OpenDb("sqlite3", "./Database/Database.sqlite")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	mux.HandleFunc("/", handler.Redirect)

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) { handler.Login(w, r) })
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) { middleware.RegisterMiddleware(handler.Register, db) })

	mux.HandleFunc("/home", handler.Home)
}
