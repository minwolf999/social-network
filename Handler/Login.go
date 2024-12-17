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
This function handles user login by verifying the provided credentials.

It takes a pointer to an SQL database as an argument and returns an http.HandlerFunc.

The function checks for valid login data, validates credentials, and manages connection attempts.
*/
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize a new ResponseWriter for structured error handling.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Create a variable to hold the login data.
		var loginData model.Auth

		// Decode the incoming JSON request body into the loginData structure.
		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [Login] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Store the password for later comparison.
		password := loginData.Password

		// Check for empty email or password fields.
		if loginData.Email == "" || loginData.Password == "" {
			// Return error if there is an empty field.
			nw.Error("There is an empty field")
			log.Printf("[%s] [Login] %s", r.RemoteAddr, "There is an empty field")
			return
		}

		// Attempt to retrieve the user data from the database.
		if err := loginData.SelectFromDb(db, map[string]any{"Email": loginData.Email}); err != nil {
			// Return error if there is a problem during the database query.
			nw.Error("Internal error: Problem during database query: " + err.Error())
			log.Printf("[%s] [Login] %s", r.RemoteAddr, err.Error())
			return
		}

		// Check if the password is empty after selection.
		if loginData.Password == "" {
			// Return error if the email is incorrect.
			nw.Error("Incorrect email")
			log.Printf("[%s] [Login] %s", r.RemoteAddr, "Incorrect email")
			return
		}

		// Lock account after 10 unsuccessful login attempts.
		if loginData.ConnectionAttempt == 10 {
			//-------------------------------------------------------------------------------------------------------------------------------------
			//-------------------------------------------------------------------------------------------------------------------------------------
			// 								Send an email to reset the account
			//-------------------------------------------------------------------------------------------------------------------------------------
			//-------------------------------------------------------------------------------------------------------------------------------------

			// Return error notifying the user that their account is locked.
			nw.Error("Your account has been locked due to multiple unsuccessful logins, an email has been sent to you to reset your password and unlock your account")
			log.Printf("[%s] [Login] Your account has been locked for too many connection attempts", r.RemoteAddr)
			return
		} else if loginData.ConnectionAttempt > 10 {
			// Return error notifying the user that their account is banned.
			nw.Error("Your account has been banned")
			log.Printf("[%s] [Login] Your account has been banned", r.RemoteAddr)
			return
		}

		// Compare the hashed password with the provided password.
		if err := bcrypt.CompareHashAndPassword([]byte(loginData.Password), []byte(password)); err != nil {
			// Return error if the password is invalid.
			nw.Error("Invalid password")
			log.Printf("[%s] [Login] %s", r.RemoteAddr, err.Error())

			// Increment the connection attempt count.
			loginData.ConnectionAttempt++
			if err = loginData.UpdateDb(db, map[string]any{"ConnectionAttempt": loginData.ConnectionAttempt}, map[string]any{"Id": loginData.Id}); err != nil {
				log.Printf("Error during the update in the Db: %v", err)
			}

			return
		}

		// Reset connection attempts upon successful login.
		loginData.ConnectionAttempt = 0
		if err := loginData.UpdateDb(db, map[string]any{"ConnectionAttempt": loginData.ConnectionAttempt}, map[string]any{"Id": loginData.Id}); err != nil {
			log.Printf("Error during the update in the Db: %v", err)
		}

		// Set response headers for JSON content.
		w.Header().Set("Content-Type", "application/json")

		// Encode the response JSON for successful login.
		err := json.NewEncoder(w).Encode(map[string]any{
			"Success":   true,
			"Message":   "Login successfully",
			"sessionId": utils.GenerateJWT(loginData.Id),
		})
		if err != nil {
			log.Printf("[%s] [Login] %s", r.RemoteAddr, err.Error())
		}
	}
}

