package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	model "social-network/Model"
	utils "social-network/Utils"

	"golang.org/x/crypto/bcrypt"
)

func HandleChangeUserData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var request struct {
			SessionId string `json:"SessionId"`
			NewName   string `json:"NewName"`
			NewPass   string `json:"NewPass"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [ChangeUserData] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		uid, err := utils.DecryptJWT(request.SessionId, db)
		if err != nil {
			nw.Error("Error when decrypt the JWT")
			log.Printf("[%s] [ChangeUserData] %s", r.RemoteAddr, err.Error())
			return
		}

		if request.NewName != "" {
			if err := ChangeUserName(db, request.NewName, uid); err != nil {
				nw.Error("Error changing username")
				log.Printf("[%s] [ChangeUserData] Error changing username: %v", r.RemoteAddr, err)
				return
			}
		}

		if request.NewPass != "" {
			if err := ChangePass(db, request.NewPass, uid); err != nil {
				nw.Error("Error changing password")
				log.Printf("[%s] [ChangeUserData] Error changing password: %v", r.RemoteAddr, err)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "User data updated successfully",
		})
		if err != nil {
			log.Printf("[%s] [ChangeUserData] %s", r.RemoteAddr, err.Error())
		}
	}
}

func ChangeUserName(db *sql.DB, name, uid string) error {
	actualname, err := utils.SelectFromDb("UserInfo", db, map[string]any{"Id": uid})
	if err != nil {
		return err
	}

	userdata, err := utils.ParseRegisterData(actualname[0])
	if err != nil {
		log.Println("Error Parsing Data", err)
		return err
	}

	if userdata.Username == name {
		log.Println("New Username and current Username are the same", err)
		return err
	} else {
		utils.UpdateDb("UserInfo", db, map[string]any{"Username": name}, map[string]any{"Id": uid})
		fmt.Println("Change username Succes")
	}

	return nil
}

func ChangePass(db *sql.DB, newpass, uid string) error {
	actualpass, err := utils.SelectFromDb("Auth", db, map[string]any{"Id": uid})
	if err != nil {
		return err
	}

	userdata, err := utils.ParseAuthData(actualpass[0])
	if err != nil {
		log.Println("Error Parsing Data", err)
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userdata.Password), []byte(newpass)); err == nil {
		log.Println("This Password is already used", err)
		return err
	} else {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		utils.UpdateDb("Auth", db, map[string]any{"Password": hashedPass}, map[string]any{"Id": uid})
		fmt.Println("Change password Succes")
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
