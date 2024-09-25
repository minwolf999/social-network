package routes

import (
	"fmt"
	"net/http"

	handler "social-network/Handler"
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

	mux.Handle("/", handler.Redirect())

	// Log routes
	mux.Handle("/login", handler.Login(db))
	mux.Handle("/register", handler.Register(db))

	// Cookie route
	mux.Handle("/verificationSessionId", handler.VerificationSessionId(db))
	
	// UserDatas routes
	mux.Handle("/getUser", handler.GetUser(db))

	// Posts routes
	mux.Handle("/createPost", handler.CreatePost(db))
	mux.Handle("/getPost", handler.GetPost(db))

	// Comments routes
	mux.Handle("/createComment/{postId}", handler.CreateComment(db))
	mux.Handle("/getComment/{postId}", handler.GetComment(db))

	// Followers routes
	mux.Handle("/addFollowed", handler.AddFollower(db))
	mux.Handle("/removeFollowed", handler.RemoveFollower(db))
	mux.Handle("/getFollowed", handler.GetFollowed(db))
	mux.Handle("/getFollower", handler.GetFollower(db))

	mux.Handle("/like", handler.HandleLike(db))
	mux.Handle("/settings", handler.HandleChangeUserData(db))
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
