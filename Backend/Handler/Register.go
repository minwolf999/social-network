package handler

import (
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
func Register(w http.ResponseWriter, r *http.Request) {
	// We create a variable of http.responseWriter personnalised
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

	// We insert in the table Auth of the db the id, email and password of the people trying to register
	if err := utils.InsertIntoDb("Auth", register.Auth.Id, register.Auth.Email, register.Auth.Password); err != nil {
		nw.Error("Internal Error: There is a probleme during the push in the DB")
		return
	}

	// We insert in the table UserInfo of the db the rest of the values
	if err := utils.InsertIntoDb("UserInfo", register.Auth.Id, register.Auth.Email, register.FirstName, register.LastName, register.BirthDate, register.ProfilePicture, register.Username, register.AboutMe); err != nil {
		fmt.Println(err)
		nw.Error("Internal Error: There is a probleme during the push in the DB")
		return
	}

	// We send a success response to the request
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Register successfully")
}
