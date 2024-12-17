package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"
)

/*
VerificationSessionId handles the verification of a provided session ID.

It checks if the session ID is valid and retrieves the associated user information from the database.
*/
func VerificationSessionId(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize a new ResponseWriter for structured error handling.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Create a variable to hold the session ID from the request body.
		var sessionId string

		// Decode the incoming JSON request body to extract the session ID.
		if err := json.NewDecoder(r.Body).Decode(&sessionId); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [VerificationSessionId] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the session ID to retrieve the associated user ID.
		decryptId, err := utils.DecryptJWT(sessionId, db)
		if err != nil {
			// Return error if JWT decryption fails.
			nw.Error("Invalid JWT")
			log.Printf("[%s] [VerificationSessionId] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		// Initialize a new Auth struct to hold user authentication data.
		var auth model.Auth
		// Set the ID to the decrypted value.
		auth.Id = decryptId

		// Query the database to retrieve user authentication data based on the user ID.
		if err = auth.SelectFromDb(db, map[string]any{"Id": auth.Id}); err != nil {
			// Return error if there is a problem during the database query.
			nw.Error("Internal error: Problem during database query: " + err.Error())
			log.Printf("[%s] [VerificationSessionId] %s", r.RemoteAddr, err.Error())
			return
		}

		// Check if the authentication data is valid for cookie verification.
		if err := CheckDatasForCookie(auth); err != nil {
			// Return error if cookie verification fails.
			nw.Error(err.Error())
			log.Printf("[%s] [VerificationSessionId] %s", r.RemoteAddr, err.Error())
			return
		}

		// Set response headers for JSON content.
		w.Header().Set("Content-Type", "application/json")

		// Encode the response JSON indicating successful verification.
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Valid cookie",
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [VerificationSessionId] %s", r.RemoteAddr, err.Error())
		}
	}
}


/*
This function takes 1 argument:
  - a map who contain the value of the select and the name of the colum in the db selected

The purpose of this function is to check if the datas are empty or not.

The function return 1 value:
  - an error
*/
func CheckDatasForCookie(authData model.Auth) error {
	// We check if the datas get are empty or not
	if authData.Id == "" || authData.Email == "" || authData.Password == "" {
		return errors.New("nobody have this Id")
	}

	return nil
}
