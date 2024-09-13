package routes

import (
	"fmt"
	"net/http"

	handler "social-network/Handler"
	middleware "social-network/Middleware"
	utils "social-network/Utils"
)

/*
This function takes 1 argument:
  - an *http.ServeMux who is used to create the server

The purpose of this function is to create all the server endpoints and define the function call for each endpoint.

The function have no return
*/
func Routes(mux *http.ServeMux) {
	db, err := utils.OpenDb("sqlite3", "./Database/Database.sqlite")
	if err != nil {
		fmt.Println(err)
	}

	mux.HandleFunc("/", handler.Redirect)

	mux.HandleFunc("/login", middleware.LoginMiddleware(handler.Login, db))
	mux.HandleFunc("/register",  middleware.RegisterMiddleware(handler.Register, db))

	mux.HandleFunc("/home", handler.Home)
}

// Mock Login handler for testing
func mockLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate an OK response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// Mock Register handler for testing
func mockRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate an OK response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Register successful"))
}

// Routes function for testing with mocks
func RoutesForTest(mux *http.ServeMux) {
	// We replace real handlers with mocks
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Redirected"))
	})

	mux.HandleFunc("/login", mockLoginHandler)
	mux.HandleFunc("/register", mockRegisterHandler)
	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome Home"))
	})
}
