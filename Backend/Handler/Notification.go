package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"
)

func GetAllNotifications(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var userId string
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&userId); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetAllNotifications] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptUserId, err := utils.DecryptJWT(userId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [GetAllNotifications] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted Organisator ID
		userId = decryptUserId

		var notifications model.Notifications
		if err = notifications.SelectFromDb(db, map[string]any{"UserId": userId}); err != nil {
			nw.Error("Error during the fetch of the DB") // Handle invalid JWT error
			log.Printf("[%s] [GetAllNotifications] Error during the fetch of the DB : %v", r.RemoteAddr, err)
			return
		}

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Notifications getted successfully",
			"Value": notifications,
		})
		if err != nil {
			log.Printf("[%s] [GetAllNotifications] %s", r.RemoteAddr, err.Error())
		}
	}
}

