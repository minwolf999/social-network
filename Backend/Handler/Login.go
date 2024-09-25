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

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var loginData model.Auth
		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [Login] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We look if all is good in the datas send in the body of the request
		if loginData.Email == "" || loginData.Password == "" {
			nw.Error("There is an empty field")
			log.Printf("[%s] [Login] %s", r.RemoteAddr, "There is an empty field")
			return
		}

		// We get the row in the db where the email is equal to the email send
		authData, err := utils.SelectFromDb("Auth", db, map[string]any{"Email": loginData.Email})
		if err != nil {
			nw.Error("Internal error: Problem during database query: " + err.Error())
			log.Printf("[%s] [Login] %s", r.RemoteAddr, err.Error())
			return
		}

		// We check if there is no result
		if len(authData) != 1 {
			nw.Error("Incorrect email")
			log.Printf("[%s] [Login] %s", r.RemoteAddr, "Incorrect email")
			return
		}

		// We parse the result into a good structure
		userData, err := utils.ParseAuthData(authData[0])
		if err != nil {
			nw.Error(err.Error())
			log.Printf("[%s] [Login] %s", r.RemoteAddr, err.Error())
			return
		}

		if userData.ConnectionAttempt >= 10 {
			//-------------------------------------------------------------------------------------------------------------------------------------
			//-------------------------------------------------------------------------------------------------------------------------------------
			// 								Send an email to reset the account
			//-------------------------------------------------------------------------------------------------------------------------------------
			//-------------------------------------------------------------------------------------------------------------------------------------

			nw.Error("Your account has been locked due to multiple unsuccessful logins, an email has been sent to you to reset your password and unlock your account")
			log.Printf("[%s] [Login] Your account have been locked for to many connection", r.RemoteAddr)
			return
		}

		// We compare the password give and the crypted password
		if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginData.Password)); err != nil {
			nw.Error("Invalid password")
			log.Printf("[%s] [Login] %s", r.RemoteAddr, err.Error())
			
			userData.ConnectionAttempt++
			if err = utils.UpdateDb("Auth", db, map[string]any{"ConnectionAttempt": userData.ConnectionAttempt}, map[string]any{"Id": userData.Id}); err != nil {
				log.Printf("Error during the update in the Db: %v", err)
			}

			return
		}

		userData.ConnectionAttempt = 0
		if err = utils.UpdateDb("Auth", db, map[string]any{"ConnectionAttempt": userData.ConnectionAttempt}, map[string]any{"Id": userData.Id}); err != nil {
			log.Printf("Error during the update in the Db: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success":   true,
			"Message":   "Login successfully",
			"sessionId": utils.GenerateJWT(userData.Id),
		})
		if err != nil {
			log.Printf("[%s] [Login] %s", r.RemoteAddr, err.Error())
		}
	}
}
