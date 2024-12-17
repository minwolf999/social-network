package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	model "social-network/Model"
	utils "social-network/Utils"

	"github.com/gofrs/uuid"
)

// AddFollower handles the addition of a follower for a user.
// It expects the request body to contain the follower details in JSON format.
func AddFollower(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to handle error responses
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a Follower object to hold the decoded request body
		var follower model.Follower
		// Decode the JSON request body into the follower object
		if err := json.NewDecoder(r.Body).Decode(&follower); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [AddFollower] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(follower.FollowerId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [AddFollower] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		follower.FollowerId = decryptAuthorId // Set the decrypted user ID

		// Check if FollowerId is provided
		if follower.FollowerId == "" {
			nw.Error("There is no id for the user to follow")
			log.Printf("[%s] [AddFollower] There is no id for the user to follow", r.RemoteAddr)
			return
		}

		// Generate a new UUID for the follower relationship
		uid, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
			log.Printf("[%s] [AddFollower] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}
		follower.Id = uid.String() // Set the generated UUID as the follower ID

		// Check if the user exists in the Auth table
		if err = utils.IfExistsInDB("Auth", db, map[string]any{"Id": follower.FollowerId}); err != nil {
			nw.Error("There is no user with the id of the JWT : " + follower.FollowerId)
			log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, "There is no user with the id of the JWT")
			return
		}

		// Check if the follower user exists in the Auth table
		if err = utils.IfExistsInDB("Auth", db, map[string]any{"Id": follower.FollowerId}); err != nil {
			nw.Error("The Id of the person you want to follow doesn't exist")
			log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, "The Id of the person you want to follow doesn't exist")
			return
		}

		var followerData model.Register
		followerData.Id = follower.FollowerId

		if err = followerData.SelectFromDb(db, map[string]any{"Id": followerData.Id}); err != nil {
			nw.Error("There is a probleme during the fetching of the other user's data")
			log.Printf("[%s] [AddFollower] There is a probleme during the fetching of the other user's data : %s", r.RemoteAddr, err)
			return
		}

		// Check if the follower relationship already exists
		if followerData.Email == "" {
			nw.Error("The user already follows this user")
			log.Printf("[%s] [AddFollower] The user already follows this user", r.RemoteAddr)
			return
		}

		notifMessage := ""

		var followedData model.Register
		followedData.Id = follower.FollowedId
		
		if err = followedData.SelectFromDb(db, map[string]any{"Id": followedData.Id}); err != nil {
			nw.Error("There is a probleme during the fetching of the other user's data")
			log.Printf("[%s] [AddFollower] There is a probleme during the fetching of the other user's data : %s", r.RemoteAddr, err)
			return
		}

		if followedData.Status == "public" {
			notifMessage = "You have been followed"

			// Insert the follower relationship into the database
			if err := follower.InsertIntoDb(db); err != nil {
				nw.Error("Internal Error: There is a problem during the push in the DB: " + err.Error())
				log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, err.Error())
				return
			}
		} else if followedData.Status == "private" {
			notifMessage = "You have receive a followed request"

			// Check if the follower user exists in the Auth table
			if err = utils.IfNotExistsInDB("FollowingRequest", db, map[string]any{"FollowerId": follower.FollowerId, "FollowedId": follower.FollowedId}); err != nil {
				nw.Error("The request has already been send")
				log.Printf("[%s] [AddFollower] The request has already been send : %s", r.RemoteAddr, err)
				return
			}

			var followRequest = model.FollowRequest{
				FollowerId: follower.FollowerId,
				FollowedId: follower.FollowedId,
			}

			if err = followRequest.InsertIntoDb(db); err != nil {
				nw.Error("Impossible to send the request")
				log.Printf("[%s] [AddFollower] Impossible to send the request : %s", r.RemoteAddr, err)
			}
		}

		notifId, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
			log.Printf("[%s] [AddFollower] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}

		var userData model.Register
		if err = userData.SelectFromDb(db, map[string]any{"Id": follower.FollowerId}); err != nil {
			nw.Error("There is a problem during the fetching of the user") // Handle UUID generation error
			log.Printf("[%s] [AddFollower] There is a problem during the fetching of the user : %s", r.RemoteAddr, err)
			return
		}

		var userDataName string
		if userData.Username == "" {
			userDataName = userData.FirstName + " " + userData.LastName
		} else {
			userDataName = userData.Username
		}

		notification := model.Notification{
			Id:          notifId.String(),
			UserId:      follower.FollowedId,
			Status:      "Follow",
			Description: fmt.Sprintf("%s by %s", notifMessage, userDataName),
			GroupId:     "",
			OtherUserId: "",
		}

		if err = notification.InsertIntoDb(db); err != nil {
			nw.Error("There is a probleme during the sending of a notification")
			log.Printf("[%s] [AddFollower] There is a probleme during the sending of a notification : %s", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[follower.FollowerId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				Description string
				Value       model.Follower
				Notification model.Notification
			}

			WebsocketMessage.Type = "Follow"
			WebsocketMessage.Description = notifMessage
			WebsocketMessage.Value = follower
			WebsocketMessage.Notification = notification

			if err = model.ConnectedWebSocket.Conn[follower.FollowerId].WriteJSON(WebsocketMessage); err != nil {
				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [AddFollower] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}

			_, isOk2 := model.ConnectedWebSocket.Conn[follower.FollowedId] 
			if isOk2 {
				if err = model.ConnectedWebSocket.Conn[follower.FollowedId].WriteJSON(WebsocketMessage); err != nil {
					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [AddFollower] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Add follower successfully",
		})
		if err != nil {
			log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, err.Error())
		}
	}
}

// RemoveFollower handles the removal of a follower from a user.
// It expects the request body to contain the follower details in JSON format.
func RemoveFollower(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a Follower object to hold the decoded request body
		var follower model.Follower
		// Decode the JSON request body into the follower object
		if err := json.NewDecoder(r.Body).Decode(&follower); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [RemoveFollower] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(follower.FollowerId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [RemoveFollower] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		follower.FollowerId = decryptAuthorId // Set the decrypted user ID

		// Check if FollowerId is provided
		if follower.FollowerId == "" {
			nw.Error("There is no id for the user to unfollow")
			log.Printf("[%s] [RemoveFollower] There is no id for the user to unfollow", r.RemoteAddr)
			return
		}

		// Attempt to delete the follower relationship from the database
		if err = follower.DeleteFromDb(db, map[string]any{"FollowerId": follower.FollowerId, "FollowedId": follower.FollowedId}); err != nil {
			nw.Error("Internal Error: There is a problem during the delete in the DB: " + err.Error())
			log.Printf("[%s] [RemoveFollower] %s", r.RemoteAddr, err.Error())
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[follower.FollowerId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				Description string
				Value       model.Follower
			}

			WebsocketMessage.Type = "Follow"
			WebsocketMessage.Description = "You have unfollow"
			WebsocketMessage.Value = follower

			if err = model.ConnectedWebSocket.Conn[follower.FollowerId].WriteJSON(WebsocketMessage); err != nil {
				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [AddFollower] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}

			_, isOk2 := model.ConnectedWebSocket.Conn[follower.FollowedId]
			if isOk2 {
				if err = model.ConnectedWebSocket.Conn[follower.FollowedId].WriteJSON(WebsocketMessage); err != nil {
					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [AddFollower] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Remove follower successfully",
		})
		if err != nil {
			log.Printf("[%s] [RemoveFollower] %s", r.RemoteAddr, err.Error())
		}
	}
}

// RemoveFollower handles the removal of a follower from a user.
// It expects the request body to contain the follower details in JSON format.
func RemoveFollowed(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a Follower object to hold the decoded request body
		var follower model.Follower
		// Decode the JSON request body into the follower object
		if err := json.NewDecoder(r.Body).Decode(&follower); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [RemoveFollower] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(follower.FollowerId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [RemoveFollower] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		follower.FollowerId = decryptAuthorId // Set the decrypted user ID

		// Check if FollowerId is provided
		if follower.FollowerId == "" {
			nw.Error("There is no id for the user to unfollow")
			log.Printf("[%s] [RemoveFollower] There is no id for the user to unfollow", r.RemoteAddr)
			return
		}

		// Attempt to delete the follower relationship from the database
		if err = follower.DeleteFromDb(db, map[string]any{"FollowerId": follower.FollowedId, "FollowedId": follower.FollowerId}); err != nil {
			nw.Error("Internal Error: There is a problem during the delete in the DB: " + err.Error())
			log.Printf("[%s] [RemoveFollower] %s", r.RemoteAddr, err.Error())
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[follower.FollowerId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				Description string
				Value       model.Follower
			}

			WebsocketMessage.Type = "Follow"
			WebsocketMessage.Description = "You have unfollow"
			WebsocketMessage.Value = follower

			if err = model.ConnectedWebSocket.Conn[follower.FollowerId].WriteJSON(WebsocketMessage); err != nil {
				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [AddFollower] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}

			_, isOk2 := model.ConnectedWebSocket.Conn[follower.FollowedId]
			if isOk2 {
				if err = model.ConnectedWebSocket.Conn[follower.FollowedId].WriteJSON(WebsocketMessage); err != nil {
					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [AddFollower] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Remove follower successfully",
		})
		if err != nil {
			log.Printf("[%s] [RemoveFollower] %s", r.RemoteAddr, err.Error())
		}
	}
}

// GetFollowed retrieves the list of users followed by the authenticated user.
// It expects the request body to contain the user ID in JSON format.
func GetFollowed(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a Follower object to hold the decoded request body
		var follower struct {
			UserId      string `json:"UserId"`
			OtherUserId string `json:"OtherUserId"`
		}
		// Decode the JSON request body into the follower object
		if err := json.NewDecoder(r.Body).Decode(&follower); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [GetFollowed] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(follower.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [GetFollowed] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		follower.UserId = decryptAuthorId // Set the decrypted user ID

		// Check if the user exists in the Auth table
		if err = utils.IfExistsInDB("Auth", db, map[string]any{"Id": follower.UserId}); err != nil {
			nw.Error("Invalid Id in the JWT")
			log.Printf("[%s] [GetFollowed] Invalid Id in the JWT : %v", r.RemoteAddr, err)
			return
		}

		// Prepare a list to hold the followed users
		var follows model.Followers

		if follower.OtherUserId != "" {
			follows = model.Followers{
				{
					FollowedId: follower.OtherUserId,
				},
			}
		} else {
			follows = model.Followers{
				{
					FollowedId: follower.UserId,
				},
			}
		}

		// Retrieve the list of users followed by the authenticated user from the database
		if err := follows.SelectFromDb(db, map[string]any{"FollowedId": follows[0].FollowedId}); err != nil {
			nw.Error("Internal Error: There is a problem during the select in the DB: " + err.Error())
			log.Printf("[%s] [GetFollowed] %s", r.RemoteAddr, err.Error())
			return
		}

		// Send the list of followed users as a response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get followed successfully",
			"Follow":  follows,
		})
		if err != nil {
			log.Printf("[%s] [GetFollowed] %s", r.RemoteAddr, err.Error())
		}
	}
}

// GetFollower retrieves the list of followers for a user.
// It expects the request body to contain the user ID in JSON format.
func GetFollower(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a struct to hold the decoded request body
		var follower struct {
			UserId      string `json:"UserId"`
			OtherUserId string `json:"OtherUserId"`
		}
		// Decode the JSON request body into the follower struct
		if err := json.NewDecoder(r.Body).Decode(&follower); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [GetFollower] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(follower.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [GetFollower] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		follower.UserId = decryptAuthorId // Set the decrypted user ID

		// Check if the user exists in the Auth table
		if err = utils.IfExistsInDB("Auth", db, map[string]any{"Id": follower.UserId}); err != nil {
			nw.Error("Invalid Id in the JWT")
			log.Printf("[%s] [GetFollower] Invalid Id in the JWT : %v", r.RemoteAddr, err)
			return
		}

		// Prepare a list to hold the followers of the user
		var follows model.Followers

		if follower.OtherUserId != "" {
			follows = model.Followers{
				{
					FollowerId: follower.OtherUserId,
				},
			}
		} else {
			follows = model.Followers{
				{
					FollowerId: follower.UserId,
				},
			}
		}

		// Retrieve the list of followers from the database
		if err := follows.SelectFromDb(db, map[string]any{"FollowerId": follows[0].FollowerId}); err != nil {
			nw.Error("Internal Error: There is a problem during the select in the DB: " + err.Error())
			log.Printf("[%s] [GetFollower] %s", r.RemoteAddr, err.Error())
			return
		}

		// Send the list of followers as a response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get followed successfully",
			"Follow":  follows,
		})
		if err != nil {
			log.Printf("[%s] [GetFollower] %s", r.RemoteAddr, err.Error())
		}
	}
}

// IsFollowedBy checks if a user is followed by another user.
// It returns true if the follower relationship exists; false otherwise.
func IsFollowedBy(follower, followed string, db *sql.DB) bool {
	// Query the database to check if the follower relationship exists
	follow, err := model.SelectFromDb("Follower", db, map[string]any{"FollowerId": follower, "FollowedId": followed})
	if err != nil {
		return false // Return false if an error occurs during the query
	}

	// Parse the results to check if the follower relationship is valid
	follows, err := follow.ParseFollowersData()
	if err != nil {
		return false // Return false if there is an error in parsing
	}

	// Return true if exactly one follower relationship is found
	return len(follows) == 1
}

func GetFollowedRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a struct to hold the decoded request body
		var userId string
		// Decode the JSON request body into the follower struct
		if err := json.NewDecoder(r.Body).Decode(&userId); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetFollowedRequest] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(userId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetFollowedRequest] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		userId = decryptAuthorId

		var followRequests model.FollowRequests
		if err = followRequests.SelectFromDb(db, map[string]any{"FollowedId": userId}); err != nil {
			nw.Error("Error during the fetch of the database")
			log.Printf("[%s] [GetFollowedRequest] Error during the fetch of the database : %v", r.RemoteAddr, err)
			return
		}

		// Send the list of followers as a response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get followed successfully",
			"Follow":  followRequests,
		})
		if err != nil {
			log.Printf("[%s] [GetFollowedRequest] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetFollowRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a struct to hold the decoded request body
		var FollowerId string
		// Decode the JSON request body into the follower struct
		if err := json.NewDecoder(r.Body).Decode(&FollowerId); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetFollowRequest] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(FollowerId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetFollowRequest] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		FollowerId = decryptAuthorId

		var followRequests model.FollowRequests
		if err = followRequests.SelectFromDb(db, map[string]any{"FollowerId": FollowerId}); err != nil {
			nw.Error("Error during the fetch of the database")
			log.Printf("[%s] [GetFollowRequest] Error during the fetch of the database : %v", r.RemoteAddr, err)
			return
		}

		// Send the list of followers as a response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get followed successfully",
			"Follow":  followRequests,
		})
		if err != nil {
			log.Printf("[%s] [GetFollowRequest] %s", r.RemoteAddr, err.Error())
		}
	}
}

func DeclineFollowedRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a struct to hold the decoded request body
		var followedRequest model.FollowRequest
		// Decode the JSON request body into the follower struct
		if err := json.NewDecoder(r.Body).Decode(&followedRequest); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [DeclineFollowedRequest] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(followedRequest.FollowedId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [DeclineFollowedRequest] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		followedRequest.FollowedId = decryptAuthorId

		if err := utils.IfExistsInDB("FollowingRequest", db, map[string]any{"FollowerId": followedRequest.FollowerId, "FollowedId": followedRequest.FollowedId}); err != nil {
			nw.Error("There is no request for following this user")
			log.Printf("[%s] [DeclineFollowedRequest] There is no request for following this user : %s", r.RemoteAddr, err)
			return
		}

		if err := followedRequest.DeleteFromDb(db, map[string]any{"FollowerId": followedRequest.FollowerId, "FollowedId": followedRequest.FollowedId}); err != nil {
			nw.Error("Error during the delete of the request")
			log.Printf("[%s] [DeclineFollowedRequest] Error during the delete of the request: %v", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[followedRequest.FollowedId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				UserId      string
				Description string
				Value       model.FollowRequest
			}

			WebsocketMessage.Type = "DeclineFollow"
			WebsocketMessage.UserId = followedRequest.FollowerId
			WebsocketMessage.Description = "Your follow request have been denied"
			WebsocketMessage.Value = followedRequest

			if err = model.ConnectedWebSocket.Conn[followedRequest.FollowedId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [DeclineFollowedRequest] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}

			WebsocketMessage.UserId = followedRequest.FollowedId

			_, isOk2 := model.ConnectedWebSocket.Conn[followedRequest.FollowerId]
			if isOk2 {
				if err = model.ConnectedWebSocket.Conn[followedRequest.FollowerId].WriteJSON(WebsocketMessage); err != nil {
	
					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [DeclineFollowedRequest] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send the list of followers as a response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Request successfully denied",
		})
		if err != nil {
			log.Printf("[%s] [DeclineFollowedRequest] %s", r.RemoteAddr, err.Error())
		}
	}
}

func AcceptFollowedRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a struct to hold the decoded request body
		var followedRequest model.FollowRequest
		// Decode the JSON request body into the follower struct
		if err := json.NewDecoder(r.Body).Decode(&followedRequest); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [AcceptFollowedRequest] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(followedRequest.FollowedId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [DeclineFollowedRequest] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		followedRequest.FollowedId = decryptAuthorId

		if err := utils.IfExistsInDB("FollowingRequest", db, map[string]any{"FollowerId": followedRequest.FollowerId, "FollowedId": followedRequest.FollowedId}); err != nil {
			nw.Error("There is no request for following this user")
			log.Printf("[%s] [AcceptFollowedRequest] There is no request for following this user : %s", r.RemoteAddr, err)
			return
		}

		if err := utils.IfNotExistsInDB("Follower", db, map[string]any{"FollowerId": followedRequest.FollowerId, "FollowedId": followedRequest.FollowedId}); err != nil {
			nw.Error("There is already a follow between this users")
			log.Printf("[%s] [AcceptFollowedRequest] There is already a follow between this users : %s", r.RemoteAddr, err)
			return
		}

		// Generate a new UUID for the follower relationship
		uuid, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
			log.Printf("[%s] [AddFollower] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}

		var follow = model.Follower{
			Id:         uuid.String(),
			FollowerId: followedRequest.FollowerId,
			FollowedId: followedRequest.FollowedId,
		}

		if err = follow.InsertIntoDb(db); err != nil {
			nw.Error("Error during the insert in the DB")
			log.Printf("[%s] [AcceptFollowedRequest] Error during the insert in the DB: %v", r.RemoteAddr, err)
			return
		}

		if err := followedRequest.DeleteFromDb(db, map[string]any{"FollowerId": followedRequest.FollowerId, "FollowedId": followedRequest.FollowedId}); err != nil {
			nw.Error("Error during the delete of the request")
			log.Printf("[%s] [AcceptFollowedRequest] Error during the delete of the request: %v", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[followedRequest.FollowedId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				UserId      string
				Description string
				Value       model.Follower
			}

			var followDetail model.Followers
			if err = followDetail.SelectFromDb(db, map[string]any{"Id": follow.Id}); err != nil {
				nw.Error("Error during the fetch of the DB")
				log.Printf("[%s] [AcceptFollowedRequest] Error during the fetch of the DB: %v", r.RemoteAddr, err)
				return
			}

			WebsocketMessage.Type = "AcceptFollow"
			WebsocketMessage.UserId = followedRequest.FollowedId
			WebsocketMessage.Description = "Your follow request have been accepted"

			if len(followDetail) >= 1 {
				WebsocketMessage.Value = followDetail[0]
			} else {
				WebsocketMessage.Value = model.Follower{}
			}

			if err = model.ConnectedWebSocket.Conn[followedRequest.FollowedId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [AcceptFollowedRequest] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}

			_,  isOk2 := model.ConnectedWebSocket.Conn[followedRequest.FollowerId]
			if isOk2 {
				if err = model.ConnectedWebSocket.Conn[followedRequest.FollowerId].WriteJSON(WebsocketMessage); err != nil {
	
					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [DeclineFollowedRequest] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send the list of followers as a response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Request successfully accepted",
		})
		if err != nil {
			log.Printf("[%s] [AcceptFollowedRequest] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetFollowerAndFollowed(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a struct to hold the decoded request body
		var FollowerId string
		// Decode the JSON request body into the follower struct
		if err := json.NewDecoder(r.Body).Decode(&FollowerId); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetFollowerAndFollowed] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to get the actual user ID
		decryptAuthorId, err := utils.DecryptJWT(FollowerId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetFollowerAndFollowed] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		FollowerId = decryptAuthorId

		var followeds model.Followers
		if err = followeds.SelectFromDb(db, map[string]any{"FollowerId": FollowerId}); err != nil {
			nw.Error("Error during the fetch of the database")
			log.Printf("[%s] [GetFollowerAndFollowed] Error during the fetch of the database : %v", r.RemoteAddr, err)
			return
		}

		var followers model.Followers
		if err = followers.SelectFromDb(db, map[string]any{"FollowedId": FollowerId}); err != nil {
			nw.Error("Error during the fetch of the database")
			log.Printf("[%s] [GetFollowerAndFollowed] Error during the fetch of the database : %v", r.RemoteAddr, err)
			return
		}

		var followerAndFollowed []struct {
			UserId              string `json:"UserId"`
			User_Name           string `json:"User_Name"`
			User_Username       string `json:"User_Username"`
			User_ProfilePicture string `json:"User_ProfilePicture"`
		}

		for _, v := range followeds {
			for _, v2 := range followers {
				if v.FollowedId == v2.FollowerId {
					var data struct {
						UserId              string `json:"UserId"`
						User_Name           string `json:"User_Name"`
						User_Username       string `json:"User_Username"`
						User_ProfilePicture string `json:"User_ProfilePicture"`
					}

					data.UserId = v.FollowedId
					data.User_Name = v.Followed_Name
					data.User_Username = v.Followed_Username
					data.User_ProfilePicture = v.Followed_Picture

					followerAndFollowed = append(followerAndFollowed, data)
				}
			}
		}

		// Send the list of followers as a response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get followed and follower successfully",
			"Follow":  followerAndFollowed,
		})
		if err != nil {
			log.Printf("[%s] [GetFollowerAndFollowed] %s", r.RemoteAddr, err.Error())
		}
	}
}
