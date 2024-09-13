package handler

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
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
func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	responseWriter := model.ResponseWriter{
		ResponseWriter: w,
	}

	// We get the datas set in the context and Unmarshal them
	loginData, err := getLoginDataFromContext(r)
	if err != nil {
		responseWriter.Error("Internal error: Unmarshal error")
		return
	}

	// We get the row in the db where the email is equal to the email send
	authData, err := utils.SelectFromDb("Auth", db, map[string]any{"Email": loginData.Email})
	if err != nil {
		responseWriter.Error("Internal error: Problem during database query")
		return
	}

	// We check if there is no result
	if len(authData) != 1 {
		responseWriter.Error("Incorrect email")
		return
	}

	// We parse the result into a good structure
	userData, err := parseUserData(authData[0])
	if err != nil {
		responseWriter.Error(err.Error())
		return
	}

	// We compare the password give and the crypted password
	if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginData.Password)); err != nil {
		responseWriter.Error("Invalid password")
		return
	}

	// We set a cookie
	utils.SetCookie(w, base64.StdEncoding.EncodeToString([]byte(userData.Id)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Registered successfully")
}


/*
This function takes 1 argument:
  - an *http.Request

The purpose of this function is to get the value set in the context.

The function return 2 values:
	- an variable of type Auth
	- an error
*/
func getLoginDataFromContext(req *http.Request) (model.Auth, error) {
	// We decrypt the context with the key and stock the result in a []byte
	contextData := req.Context().Value(model.LoginCtx).([]byte)

	// We Unmarshall the result and return it
	var loginData model.Auth
	err := json.Unmarshal(contextData, &loginData)
		
	return loginData, err
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
