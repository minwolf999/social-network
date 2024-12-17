package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"
)

/*
GetUser handles the retrieval of user information from the database.

It can return information for the requesting user or another specified user.
The function uses a JWT for authorization.
*/
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize a new ResponseWriter for structured error handling.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Create a temporary struct to hold the incoming request data.
		var tmp struct {
			// The JWT session ID for the user.
			SessionId string `json:"SessionId"`
			// The ID of another user, if requested.
			OtherPeopleId string `json:"OtherPeopleId"`
		}

		// Decode the incoming JSON request body into the temporary struct.
		if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetUser] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Start with the ID of the user to retrieve.
		uid := tmp.OtherPeopleId

		// Decrypt the Session ID from the JWT to get the user ID of the requester.
		valueJWT, err := utils.DecryptJWT(tmp.SessionId, db)
		if err != nil {
			// Return error if JWT decryption fails.
			nw.Error("Error when decrypting the JWT")
			log.Printf("[%s] [GetUser] %s", r.RemoteAddr, err.Error())
			return
		}

		// If no specific user ID was provided, use the requester’s ID.
		if uid == "" {
			uid = valueJWT
		}

		// Create a variable to hold the user data.
		var userData model.Register

		// Set the ID to the user we want to retrieve.
		userData.Id = uid

		// Retrieve the user data from the database based on the user ID.
		if err = userData.SelectFromDb(db, map[string]any{"Id": userData.Id}); err != nil {
			// Return error if there is a problem retrieving user info.
			nw.Error("Error when getting user info")
			log.Printf("[%s] [GetUser] %s", r.RemoteAddr, err.Error())
			return
		}

		// Set response headers for JSON content.
		w.Header().Set("Content-Type", "application/json")

		// Encode the response JSON with user data.
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Sending user info",
			// Return the retrieved user data.
			"userInfos": userData,
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [GetUser] %s", r.RemoteAddr, err.Error())
		}
	}
}

/*
This function takes 1 argument:
  - a pointer to an SQL database object

The purpose of this function is to handle the retrieval of all the users in the database.

The function returns an http.HandlerFunc that can be used as a handler for HTTP requests.
*/
func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var userJWT string
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&userJWT); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetAllUsers] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		_, err := utils.DecryptJWT(userJWT, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetAllUsers] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		var users model.Users
		if err = users.SelectFromDb(db, map[string]any{}); err != nil {
			fmt.Println(err)
		}

		// Set the response header to indicate JSON content and respond with success message.
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Users getted successfully",

			"Users": users,
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [GetAllUsers] %s", r.RemoteAddr, err.Error())
		}
	}
}
