package handler

import (
	"database/sql"
	"encoding/json"
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
		var userInfo model.Register

		// Decode the JSON request body into the request struct
		if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [ChangeUserData] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the session ID to get the user ID
		uid, err := utils.DecryptJWT(userInfo.Id, db)
		if err != nil {
			nw.Error("Error when decrypt the JWT") // Handle JWT decryption error
			log.Printf("[%s] [ChangeUserData] Error when decrypt the JWT : %s", r.RemoteAddr, err.Error())
			return
		}
		userInfo.Id = uid

		var userPreviousData model.Register
		if err = userPreviousData.SelectFromDb(db, map[string]any{"Id": userInfo.Id}); err != nil {
			nw.Error("Error during the fetching of the DB") // Handle JWT decryption error
			log.Printf("[%s] [ChangeUserData] Error during the fetching of the DB : %s", r.RemoteAddr, err.Error())
			return
		}

		if userPreviousData.Email != userInfo.Email && userInfo.Email != ""{
			if err = userPreviousData.Auth.UpdateDb(db, map[string]any{"Email": userInfo.Email}, map[string]any{"Id": userInfo.Id}); err != nil {
				nw.Error("Error during the updating of the DB") // Handle JWT decryption error
				log.Printf("[%s] [ChangeUserData] Error during the updating of the DB : %s", r.RemoteAddr, err.Error())
				return
			}
		}

		if err = bcrypt.CompareHashAndPassword([]byte(userPreviousData.Password), []byte(userInfo.Password)); err != nil && userInfo.Password != "" {
			if !IsValidPassword(userInfo.Password) {
				nw.Error("Incorrect password ! the password must contain 8 characters, 1 uppercase letter, 1 special character, 1 number") // Handle JWT decryption error
				log.Printf("[%s] [ChangeUserData] incorrect password ! the password must contain 8 characters, 1 uppercase letter, 1 special character, 1 number", r.RemoteAddr)
				return
			}

			hashedPass, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
			if err != nil {
				nw.Error("Error during the hashing of the password") // Handle JWT decryption error
				log.Printf("[%s] [ChangeUserData] Error during the hashing of the password : %s", r.RemoteAddr, err.Error())
				return
			}
	
			if err = userPreviousData.Auth.UpdateDb(db, map[string]any{"Password": string(hashedPass)}, map[string]any{"Id": userInfo.Id}); err != nil {
				nw.Error("Error during the hashing of the password") // Handle JWT decryption error
				log.Printf("[%s] [ChangeUserData] Error during the hashing of the password : %s", r.RemoteAddr, err.Error())
				return
			}
		}

		if userPreviousData.Username != userInfo.Username && userInfo.Username != "" {
			if err = userPreviousData.UpdateDb(db, map[string]any{"Username": userInfo.Username}, map[string]any{"Id": userInfo.Id}); err != nil {
				nw.Error("Error during the updating of the DB") // Handle JWT decryption error
				log.Printf("[%s] [ChangeUserData] Error during the updating of the DB : %s", r.RemoteAddr, err.Error())
				return
			}
		}

		if userPreviousData.AboutMe != userInfo.AboutMe {
			if err = userPreviousData.UpdateDb(db, map[string]any{"AboutMe": userInfo.AboutMe}, map[string]any{"Id": userInfo.Id}); err != nil {
				nw.Error("Error during the updating of the DB") // Handle JWT decryption error
				log.Printf("[%s] [ChangeUserData] Error during the updating of the DB : %s", r.RemoteAddr, err.Error())
				return
			}
		}

		if userPreviousData.Status != userInfo.Status && (userInfo.Status == "private" || userInfo.Status == "public") {
			if err = userPreviousData.UpdateDb(db, map[string]any{"Status": userInfo.Status}, map[string]any{"Id": userInfo.Id}); err != nil {
				nw.Error("Error during the updating of the DB") // Handle JWT decryption error
				log.Printf("[%s] [ChangeUserData] Error during the updating of the DB : %s", r.RemoteAddr, err.Error())
				return
			}
		}

		if userPreviousData.ProfilePicture != userInfo.ProfilePicture && userInfo.ProfilePicture != "" {
			if err = userPreviousData.UpdateDb(db, map[string]any{"ProfilePicture": userInfo.ProfilePicture}, map[string]any{"Id": userInfo.Id}); err != nil {
				nw.Error("Error during the updating of the DB") // Handle JWT decryption error
				log.Printf("[%s] [ChangeUserData] Error during the updating of the DB : %s", r.RemoteAddr, err.Error())
				return
			}
		}

		if userPreviousData.Banner != userInfo.Banner && userInfo.Banner != "" {
			if err = userPreviousData.UpdateDb(db, map[string]any{"Banner": userInfo.Banner}, map[string]any{"Id": userInfo.Id}); err != nil {
				nw.Error("Error during the updating of the DB") // Handle JWT decryption error
				log.Printf("[%s] [ChangeUserData] Error during the updating of the DB : %s", r.RemoteAddr, err.Error())
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
