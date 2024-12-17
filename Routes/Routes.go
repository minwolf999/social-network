package routes

import (
	"fmt"
	"net/http"

	handler "social-network/Handler"
	model "social-network/Model"
	utils "social-network/Utils"
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
	mux.Handle("/getAllUsers", handler.GetAllUsers(db))

	// Posts routes
	mux.Handle("/createPost", handler.CreatePost(db))
	mux.Handle("/getPost", handler.GetPost(db))

	// Comments routes
	mux.Handle("/createComment", handler.CreateComment(db))
	mux.Handle("/getComment", handler.GetComment(db))

	// Followers routes
	mux.Handle("/addFollowed", handler.AddFollower(db))
	mux.Handle("/getFollowed", handler.GetFollowed(db))
	mux.Handle("/getFollower", handler.GetFollower(db))

	mux.Handle("/getFollowerAndFollowed", handler.GetFollowerAndFollowed(db))

	mux.Handle("/getFollowedRequest", handler.GetFollowedRequest(db))
	mux.Handle("/getFollowRequest", handler.GetFollowRequest(db))
	mux.Handle("/declineFollowedRequest", handler.DeclineFollowedRequest(db))
	mux.Handle("/acceptFollowedRequest", handler.AcceptFollowedRequest(db))

	mux.Handle("/removeFollowed", handler.RemoveFollowed(db))
	mux.Handle("/removeFollower", handler.RemoveFollower(db))

	// Like routes
	mux.Handle("/like", handler.HandleLike(db))
	mux.Handle("/getLikePost", handler.GetLike(db))

	// Setting route
	mux.Handle("/updateUserInfo", handler.HandleChangeUserData(db))

	// Group routes
	mux.Handle("/createGroup", handler.CreateGroup(db))
	mux.Handle("/leaveGroup", handler.LeaveGroup(db))
	mux.Handle("/getGroup", handler.GetGroup(db))
	mux.Handle("/getAllGroups", handler.GetAllGroups(db))
	mux.Handle("/getGroupsJoined", handler.GetGroupsJoined(db))
	mux.Handle("/getGroupsPosts", handler.GetGroupsPosts(db))
	mux.Handle("/deleteGroup", handler.DeleteGroup(db))

	mux.Handle("/joinGroup", handler.JoinGroup(db))
	mux.Handle("/getSendJoinRequest", handler.GetSendJoinRequest(db))
	mux.Handle("/getJoinRequest", handler.GetJoinRequest(db))
	mux.Handle("/declineJoinRequest", handler.DeclineJoinRequest(db))
	mux.Handle("/acceptJoinRequest", handler.AcceptJoinRequest(db))

	mux.Handle("/inviteGroup", handler.InviteGroup(db))
	mux.Handle("/getInvitationGroup", handler.GetInvitationGroup(db))
	mux.Handle("/getInvitationUserInGroup", handler.GetInvitationUserInGroup(db))
	mux.Handle("/declineInvitationGroup", handler.DeclineInvitationGroup(db))
	mux.Handle("/acceptInvitationGroup", handler.AcceptInvitationGroup(db))

	// Event routes
	mux.Handle("/createEvent", handler.CreateEvent(db))
	mux.Handle("/joinEvent", handler.JoinEvent(db))
	mux.Handle("/declineEvent", handler.DeclineEvent(db))
	mux.Handle("/getEvent", handler.GetEvent(db))
	mux.Handle("/getAllGroupEvents", handler.GetAllGroupEvents(db))
	mux.Handle("/getJoinedEvent", handler.GetJoinedEvent(db))

	// Notification routes
	mux.Handle("/getAllNotifications", handler.GetAllNotifications(db))
	mux.Handle("/getGroupNotification", handler.GetGroupNotification(db))
	mux.Handle("/getUserNotification", handler.GetUserNotification(db))
	mux.Handle("/deleteAllNotifications", handler.DeleteAllNotifications(db))
	mux.Handle("/deleteAllGroupNotifications", handler.DeleteAllGroupNotifications(db))
	mux.Handle("/deleteAllUserNotifications", handler.DeleteAllUserNotifications(db))

	// Chat routes
	mux.Handle("/addMessage", handler.AddMessage(db))
	mux.Handle("/getMessage", handler.GetMessage(db))

	// Websocket route
	mux.Handle("/websocket/", handler.Websocket(db))

	go utils.AutoDeleteEvent(db)
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
