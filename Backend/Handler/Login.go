package handler

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"

	"golang.org/x/crypto/bcrypt"
)

/*
This function takes 2 arguments:
  - an http.ResponseWriter
  - an *http.Request

The purpose of this function is to handle the login endpoint.

The function return no value
*/
func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r * http.Request) {
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
			return
		}
	
		// We get the row in the db where the email is equal to the email send
		authData, err := utils.SelectFromDb("Auth", db, map[string]any{"Email": loginData.Email})
		if err != nil {
			nw.Error("Internal error: Problem during database query")
			return
		}
	
		// We check if there is no result
		if len(authData) != 1 {
			nw.Error("Incorrect email")
			return
		}
	
		// We parse the result into a good structure
		userData, err := parseUserData(authData[0])
		if err != nil {
			nw.Error(err.Error())
			return
		}
	
		// We compare the password give and the crypted password
		if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginData.Password)); err != nil {
			nw.Error("Invalid password")
			return
		}
	
		// We set a cookie
		utils.SetCookie(w, base64.StdEncoding.EncodeToString([]byte(userData.Id)))
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"Success":   true,
			"Message":   "Login successfully",
			"sessionId": base64.StdEncoding.EncodeToString([]byte(userData.Id)),
		})
	}
}

/*
This function takes 1 argument:
  - a map who contain the value of the select and the name of the colum in the db selected

The purpose of this function is to parse the datas into a good structure.

The function return 2 values:
  - an variable of type Auth
  - an error
*/
func parseUserData(userData map[string]any) (model.Auth, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return model.Auth{}, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var authResult model.Auth
	err = json.Unmarshal(serializedData, &authResult)

	return authResult, err
}
