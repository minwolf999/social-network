package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"
)

func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var tmp struct {
			SessionId     string `json:"SessionId"`
			OtherPeopleId string `json:"OtherPeopleId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&tmp); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [Login] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		uid := tmp.OtherPeopleId

		valueJWT, err := utils.DecryptJWT(tmp.SessionId, db)
		if err != nil {
			nw.Error("Error when decrypt the JWT")
			log.Printf("[%s] [Settings] %s", r.RemoteAddr, err.Error())
			return
		}

		if uid == "" {
			uid = valueJWT
		}

		var userData model.Register
		userData.Id = uid

		if err = displayInfos(db, userData); err != nil {
			nw.Error("Error when get infos")
			log.Printf("[%s] [Settings] %s", r.RemoteAddr, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success":   true,
			"Message":   "Sending Infos",
			"userInfos": userData,
		})
		if err != nil {
			log.Printf("[%s] [Settings] %s", r.RemoteAddr, err.Error())
		}
	}
}

func displayInfos(db *sql.DB, userData model.Register) (error) {
	return userData.SelectFromDb(db, map[string]any{"Id": userData.Id})
}
