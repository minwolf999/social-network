package handler

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
		userData, err := ParseUserData(authData[0])
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
			"sessionId": GenerateJWT(userData.Id),
		})
		if err != nil {
			log.Printf("[%s] [Login] %s", r.RemoteAddr, err.Error())
		}
	}
}


/*
This function takes 1 argument:
  - a string who contain the value to set in the JWT

The purpose of this function is to create a JWT with a value crypted inside him.

The function return 1 value:
  - a string who is the JWT
*/
func GenerateJWT(str string) string {
	// We convert the header object in base 64 for the first part
	header := base64.StdEncoding.EncodeToString([]byte(`{
		"typ": "JWT"
	}`))

	// We convert the value in base 64 for the second part
	content := base64.StdEncoding.EncodeToString([]byte(str))

retry:
	// We hash the key for the last part od the JWT
	key, err := bcrypt.GenerateFromPassword([]byte(model.SecretKey), 12)
	if err != nil {
		fmt.Println(err)
	}
	if strings.Contains(string(key), ".") {
		goto retry
	}

	// We assemble the 3 part with a . between each part
	return header + "." + content + "." + string(key)
}


/*
This function takes 1 argument:
  - a map who contain the value of the select and the name of the colum in the db selected

The purpose of this function is to parse the datas into a good structure.

The function return 2 values:
  - an variable of type Auth
  - an error
*/
func ParseUserData(userData map[string]any) (model.Auth, error) {
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
