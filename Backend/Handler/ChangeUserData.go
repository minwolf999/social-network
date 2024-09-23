package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	model "social-network/Model"
	utils "social-network/Utils"

	"golang.org/x/crypto/bcrypt"
)

func ChangeUserName(db *sql.DB, name, uid string) error {
	actualname, err := utils.SelectFromDb("UserInfo", db, map[string]any{"Id": uid})
	if err != nil {
		return err
	}

	userdata, err := ParseUserDataInfos(actualname[0])
	if err != nil {
		log.Println("Error Parsing Data", err)
		return err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userdata.Username), []byte(name)); err != nil {
		log.Println("New Username and current Username are the same", err)
		return err
	} else {

	}

	return nil
}






func ParseUserDataInfos(userData map[string]any) (model.Register, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return model.Register{}, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var registerResult model.Register
	err = json.Unmarshal(serializedData, &registerResult)

	return registerResult, err
}
