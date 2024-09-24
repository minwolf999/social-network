package handler

import (
	"database/sql"
	"encoding/json"
	"io"
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
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var loginData model.Auth
		json.Unmarshal(body, &loginData)

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

		// We compare the password give and the crypted password
		if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginData.Password)); err != nil {
			nw.Error("Invalid password")
			log.Printf("[%s] [Login] %s", r.RemoteAddr, err.Error())
			return
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
