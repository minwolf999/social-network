package handler

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	nw := model.ResponseWriter{
		ResponseWriter: w,
	}

	// We get the datas set in the context and Unmarshal them
	contextValue := r.Context().Value(model.RegisterCtx).([]byte)

	var register model.Register
	if err := json.Unmarshal(contextValue, &register); err != nil {
		nw.Error("Internal Error: There is an Unmarshal error")
		return
	}

	// We get the row in the db where the email is equal to the email send
	authData, err := utils.SelectFromDb("Auth", db, map[string]any{"Email": register.Auth.Email})
	if err != nil {
		nw.Error("Internal error: Problem during database query")
		return
	}

	if len(authData) != 0 {
		nw.Error("Email is already used")
		return
	}

	// We insert in the table Auth of the db the id, email and password of the people trying to register
	if err := utils.InsertIntoDb("Auth", db, register.Auth.Id, register.Auth.Email, register.Auth.Password); err != nil {
		fmt.Println(err)
		nw.Error("Internal Error: There is a probleme during the push in the DB")
		return
	}

	// We insert in the table UserInfo of the db the rest of the values
	if err := utils.InsertIntoDb("UserInfo", db, register.Auth.Id, register.Auth.Email, register.FirstName, register.LastName, register.BirthDate, register.ProfilePicture, register.Username, register.AboutMe); err != nil {
		fmt.Println(err)
		nw.Error("Internal Error: There is a probleme during the push in the DB")
		return
	}

	// Set a cookie with the id of the people but converted to base64
	utils.SetCookie(w, base64.StdEncoding.EncodeToString([]byte(register.Auth.Id)))

	// We send a success response to the request
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"Success": true,
		"Message": "Register successfully",
		"sessionId": base64.StdEncoding.EncodeToString([]byte(register.Auth.Id)),
	})
}
