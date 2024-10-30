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
		decryptAuthorId, err := utils.DecryptJWT(userJWT, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetAllUsers] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		
		var users model.Users
		if err = users.SelectFromDb(db, map[string]any{}); err != nil {
			fmt.Println(err)
		}

		for i := len(users)-1; i > 0; i-- {
			if users[i].Auth.Id == decryptAuthorId {
				if i < len(users) {
					users = append(users[:i], users[i+1:]...)
				} else {
					users = users[:i]
				}
			}
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
