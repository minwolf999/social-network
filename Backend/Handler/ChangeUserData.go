package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
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

		var userInfo model.Register
		userInfo.Id = uid

		if request.NewName != "" {
			if err := ChangeUserName(db, request.NewName, userInfo); err != nil {
				nw.Error("Error changing username")
				log.Printf("[%s] [ChangeUserData] Error changing username: %v", r.RemoteAddr, err)
				return
			}
		}

		if request.NewPass != "" {
			if err := ChangePass(db, request.NewPass, userInfo); err != nil {
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

func ChangeUserName(db *sql.DB, name string, userdata model.Register) error {
	err := userdata.SelectFromDbById(db)
	if err != nil {
		return err
	}

	if userdata.Username == name {
		return errors.New("new username and current username are the same")
	} else {
		return model.UpdateDb("UserInfo", db, map[string]any{"Username": name}, map[string]any{"Id": userdata.Id})
	}
}

func ChangePass(db *sql.DB, newpass string, userData model.Register) error {
	err := userData.SelectFromDbById(db)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(newpass)); err == nil {
		return errors.New("new password and current password are the same")
	} else {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		return model.UpdateDb("Auth", db, map[string]any{"Password": string(hashedPass)}, map[string]any{"Id": userData.Id})
	}
}

