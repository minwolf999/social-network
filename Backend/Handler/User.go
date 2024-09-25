package handler

import (
	"database/sql"
	"encoding/json"
	"io"
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

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var sessionId string
		json.Unmarshal(body, &sessionId)

		uid, err := utils.DecryptJWT(sessionId, db)
		if err != nil {
			nw.Error("Error when decrypt the JWT")
			log.Printf("[%s] [Settings] %s", r.RemoteAddr, err.Error())
			return
		}
		userInfos, err := displayInfos(db, uid)
		if err != nil {
			nw.Error("Error when get infos")
			log.Printf("[%s] [Settings] %s", r.RemoteAddr, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success":   true,
			"Message":   "Sending Infos",
			"userInfos": userInfos,
		})
		if err != nil {
			log.Printf("[%s] [Settings] %s", r.RemoteAddr, err.Error())
		}
	}
}

func displayInfos(db *sql.DB, uid string) ([]map[string]any, error) {
	infos, err := utils.SelectFromDb("UserInfo", db, map[string]any{"Id": uid})
	if err != nil {
		return nil, err
	}
	return infos, nil
}
