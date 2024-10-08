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

func VerificationSessionId(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var sessionId string
		if err := json.NewDecoder(r.Body).Decode(&sessionId); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [Login] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptId, err := utils.DecryptJWT(sessionId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [VerificationSessionId] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		var auth model.Auth
		auth.Id = decryptId

		// We get the People who have this id in the db
		if err = auth.SelectFromDb(db, map[string]any{"Id": auth.Id}); err != nil {
			nw.Error("Internal error: Problem during database query: " + err.Error())
			log.Printf("[%s] [VerificationSessionId] %s", r.RemoteAddr, err.Error())
			return
		}

		if err := CheckDatasForCookie(auth); err != nil {
			nw.Error(err.Error())
			log.Printf("[%s] [VerificationSessionId] %s", r.RemoteAddr, err.Error())
			return
		}

		// We send a success response to the request
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Valid cookie",
		})
		if err != nil {
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
