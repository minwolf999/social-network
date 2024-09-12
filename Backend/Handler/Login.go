package handler

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	model "social-network/Model"
	utils "social-network/Utils"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	responseWriter := model.ResponseWriter{
		ResponseWriter: w,
	}

	loginData, err := getLoginDataFromContext(r)
	if err != nil {
		responseWriter.Error("Internal error: Unmarshal error")
		return
	}

	authData, err := utils.SelectFromDb("Auth", db, map[string]any{"Email": loginData.Email})
	if err != nil {
		responseWriter.Error("Internal error: Problem during database query")
		return
	}

	if len(authData) != 1 {
		responseWriter.Error("Incorrect email")
		return
	}

	userData, err := parseUserData(authData)
	if err != nil {
		responseWriter.Error(err.Error())
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginData.Password)); err != nil {
		responseWriter.Error("Invalid password")
		return
	}

	utils.SetCookie(w, base64.StdEncoding.EncodeToString([]byte(userData.Id)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Registered successfully")
}

func getLoginDataFromContext(req *http.Request) (model.Auth, error) {
	contextData := req.Context().Value(model.LoginCtx).([]byte)
	var loginData model.Auth
	if err := json.Unmarshal(contextData, &loginData); err != nil {
		return model.Auth{}, err
	}
	return loginData, nil
}

func parseUserData(userData [][]interface{}) (model.Auth, error) {
	serializedData, err := json.Marshal(userData[0])
	if err != nil {
		return model.Auth{}, errors.New("internal error: conversion problem")
	}

	userMap := map[string]interface{}{
		"id":       strings.Split(strings.ReplaceAll(string(serializedData)[1:len(serializedData)-1], "\"", ""), ",")[0],
		"email":    strings.Split(strings.ReplaceAll(string(serializedData)[1:len(serializedData)-1], "\"", ""), ",")[1],
		"password": strings.Split(strings.ReplaceAll(string(serializedData)[1:len(serializedData)-1], "\"", ""), ",")[2],
	}

	serializedUserMap, err := json.Marshal(userMap)
	if err != nil {
		return model.Auth{}, errors.New("internal error: conversion problem")
	}

	var authResult model.Auth
	err = json.Unmarshal(serializedUserMap, &authResult)
	return authResult, err
}
