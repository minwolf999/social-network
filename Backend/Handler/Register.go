package handler

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"
)

/*
This function takes 2 arguments:
  - an http.ResponseWriter
  - an *http.Request

The purpose of this function is to handle the register endpoint.

The function return no value
*/
func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var register model.Register
		json.Unmarshal(body, &register)
		json.Unmarshal(body, &register.Auth)

		// We look if all is good in the datas send in the body of the request
		if err := utils.RegisterVerification(register); err != nil {
			nw.Error(err.Error())
			return
		}

		// We generate an UUID and crypt the password
		if err := utils.CreateUuidAndCrypt(&register); err != nil {
			nw.Error(err.Error())
			return
		}

		// We get the row in the db where the email is equal to the email send
		authData, err := utils.SelectFromDb("Auth", db, map[string]any{"Email": register.Auth.Email})
		if err != nil {
			nw.Error("Internal error: Problem during database query: " + err.Error())
			return
		}

		if len(authData) != 0 {
			nw.Error("Email is already used")
			return
		}

		// We insert in the table Auth of the db the id, email and password of the people trying to register
		if err := utils.InsertIntoDb("Auth", db, register.Auth.Id, register.Auth.Email, register.Auth.Password); err != nil {
			fmt.Println(err)
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			return
		}

		// We insert in the table UserInfo of the db the rest of the values
		if err := utils.InsertIntoDb("UserInfo", db, register.Auth.Id, register.Auth.Email, register.FirstName, register.LastName, register.BirthDate, register.ProfilePicture, register.Username, register.AboutMe); err != nil {
			fmt.Println(err)
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			return
		}

		// Set a cookie with the id of the people but converted to base64
		utils.SetCookie(w, base64.StdEncoding.EncodeToString([]byte(register.Auth.Id)))

		// We send a success response to the request
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"Success":   true,
			"Message":   "Register successfully",
			"sessionId": base64.StdEncoding.EncodeToString([]byte(register.Auth.Id)),
		})
	}
}
