package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"
)

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var register model.Register
		if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [Register] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We look if all is good in the datas send in the body of the request
		if err := utils.RegisterVerification(register); err != nil {
			nw.Error(err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
			return
		}

		// We generate an UUID and crypt the password
		if err := utils.CreateUuidAndCrypt(&register); err != nil {
			nw.Error(err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
			return
		}

		if len(register.ProfilePicture) > 400000 {
			nw.Error("To big image")
			log.Printf("[%s] [Register] To big image", r.RemoteAddr)
			return
		}

		// We get the row in the db where the email is equal to the email send
		if err := utils.IfExistsInDB("Auth", db, map[string]any{"Email": register.Auth.Email}); err != nil && err.Error() != "there is no match" {
			nw.Error("Email is already used : " + err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, "Email is already used")
			return
		}

		// We insert in the table Auth of the db the id, email and password of the people trying to register
		if err := utils.InsertIntoDb("Auth", db, register.Auth.Id, register.Auth.Email, register.Auth.Password, 0); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
			return
		}

		// We insert in the table UserInfo of the db the rest of the values
		if err := utils.InsertIntoDb("UserInfo", db, register.Auth.Id, register.Auth.Email, register.FirstName, register.LastName, register.BirthDate, register.ProfilePicture, register.Username, register.AboutMe); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
			return
		}

		// We send a success response to the request
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"Success":   true,
			"Message":   "Login successfully",
			"sessionId": utils.GenerateJWT(register.Auth.Id),
		})
		if err != nil {
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
		}
	}
}
