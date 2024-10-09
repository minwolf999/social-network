package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	model "social-network/Model"
	utils "social-network/Utils"

	"golang.org/x/crypto/bcrypt"
)

/*
This function takes 1 argument:
  - a pointer to an sql.DB instance named db, which represents the database connection.

The purpose of this function is to handle the HTTP request for changing user data, including updating the username and password.

The function returns an http.HandlerFunc that performs the following steps:
  - Initializes a custom ResponseWriter for sending error responses.
  - Defines a struct to decode the JSON request body containing the session ID, new username, and new password.
  - Decodes the request body into the defined struct and handles any errors.
  - Decrypts the session ID to retrieve the user ID using the DecryptJWT function.
  - Updates the username if a new name is provided, handling errors accordingly.
  - Updates the password if a new password is provided, ensuring it is not the same as the current password.
  - Sends a success response in JSON format if both operations succeed.
*/
func HandleChangeUserData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to handle error responses
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Define a struct for the request body
		var request struct {
			SessionId string `json:"SessionId"` // User session ID
			NewName   string `json:"NewName"`   // New username
			NewPass   string `json:"NewPass"`   // New password
		}

		// Decode the JSON request body into the request struct
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [ChangeUserData] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the session ID to get the user ID
		uid, err := utils.DecryptJWT(request.SessionId, db)
		if err != nil {
			nw.Error("Error when decrypt the JWT") // Handle JWT decryption error
			log.Printf("[%s] [ChangeUserData] %s", r.RemoteAddr, err.Error())
			return
		}

		var userInfo model.Register // Variable to hold user information

		// Update username if a new name is provided
		if request.NewName != "" {
			userInfo.Id = uid
			if err := ChangeUserName(db, request.NewName, userInfo); err != nil {
				nw.Error("Error changing username") // Handle username change error
				log.Printf("[%s] [ChangeUserData] Error changing username: %v", r.RemoteAddr, err)
				return
			}
		}

		// Update password if a new password is provided
		if request.NewPass != "" {
			userInfo.Auth.Id = uid
			if err := ChangePass(db, request.NewPass, userInfo.Auth); err != nil {
				nw.Error("Error changing password") // Handle password change error
				log.Printf("[%s] [ChangeUserData] Error changing password: %v", r.RemoteAddr, err)
				return
			}
		}

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "User data updated successfully",
		})
		if err != nil {
			log.Printf("[%s] [ChangeUserData] %s", r.RemoteAddr, err.Error())
		}
	}
}

/*
This function takes 3 arguments:
  - a pointer to an sql.DB instance named db, representing the database connection.
  - a string named name, representing the new username.
  - a model.Register instance named userdata, representing the user data.

The purpose of this function is to change the user's username in the database.

The function performs the following steps:
  - Selects the user data from the database using the user ID.
  - Checks if the new username is the same as the current one, returning an error if they are the same.
  - If they are different, it updates the database with the new username.
*/
func ChangeUserName(db *sql.DB, name string, userdata model.Register) error {
	err := userdata.SelectFromDb(db, map[string]any{"Id": userdata.Id}) // Retrieve user data
	if err != nil {
		return err // Return any errors encountered
	}

	// Check if the new username is the same as the current username
	if userdata.Username == name {
		return errors.New("new username and current username are the same")
	} else {
		// Update the database with the new username
		return userdata.UpdateDb(db, map[string]any{"Username": name}, map[string]any{"Id": userdata.Id})
	}
}

/*
This function takes 3 arguments:
  - a pointer to an sql.DB instance named db, representing the database connection.
  - a string named newpass, representing the new password.
  - a model.Auth instance named userData, representing the user's authentication data.

The purpose of this function is to change the user's password in the database.

The function performs the following steps:
  - Selects the user data from the database using the user ID.
  - Compares the new password with the current password, returning an error if they are the same.
  - If they are different, it hashes the new password and updates the database with the new hashed password.
*/
func ChangePass(db *sql.DB, newpass string, userData model.Auth) error {
	err := userData.SelectFromDb(db, map[string]any{"Id": userData.Id}) // Retrieve user data
	if err != nil {
		return err // Return any errors encountered
	}

	// Compare the new password with the current password
	if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(newpass)); err == nil {
		return errors.New("new password and current password are the same") // Return error if they are the same
	} else {
		// Hash the new password and update the database
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
		if err != nil {
			return err // Return any errors encountered during hashing
		}

		return userData.UpdateDb(db, map[string]any{"Password": string(hashedPass)}, map[string]any{"Id": userData.Id}) // Update password in database
	}
}


