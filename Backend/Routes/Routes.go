package routes

import (
	"fmt"
	"net/http"

	handler "social-network/Handler"
	model "social-network/Model"
)

/*
This function takes 1 argument:
  - an *http.ServeMux who is used to create the server

The purpose of this function is to create all the server endpoints and define the function call for each endpoint.

The function have no return
*/
func Routes(mux *http.ServeMux) {
	db, err := model.OpenDb("sqlite3", "./Database/Database.sqlite")
	if err != nil {
		fmt.Println(err)
	}

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
	mux.Handle("/createComment", handler.CreateComment(db))
	mux.Handle("/getComment", handler.GetComment(db))

	// Followers routes
	mux.Handle("/addFollowed", handler.AddFollower(db))
	mux.Handle("/removeFollowed", handler.RemoveFollower(db))
	mux.Handle("/getFollowed", handler.GetFollowed(db))
	mux.Handle("/getFollower", handler.GetFollower(db))

	// Like route
	mux.Handle("/like", handler.HandleLike(db))

	// Setting route
	mux.Handle("/settings", handler.HandleChangeUserData(db))

	// Group routes
	mux.Handle("/createGroup", handler.CreateGroup(db))
	mux.Handle("/joinOrLeaveGroup", handler.JoinAndLeaveGroup(db))
	mux.Handle("/getGroup", handler.GetGroup(db))
	mux.Handle("/deleteGroup", handler.DeleteGroup(db))

	// Event routes
	mux.Handle("/createEvent", handler.CreateEvent(db))
	mux.Handle("/joinEvent", handler.JoinEvent(db))
	mux.Handle("/declineEvent", handler.DeclineEvent(db))
	mux.Handle("/getEvent", handler.GetEvent(db))


	// Websocket route
	mux.Handle("/websocket", handler.Websocket(db))
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
