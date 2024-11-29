package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	model "social-network/Model"
	utils "social-network/Utils"

	"github.com/gofrs/uuid"
)

func AddMessage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var message model.Message

		// Decode the JSON request body into the request struct
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [AddMessage] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		uid, err := utils.DecryptJWT(message.Id, db)
		if err != nil {
			nw.Error("Error when decrypt the JWT") // Handle JWT decryption error
			log.Printf("[%s] [AddMessage] Error when decrypt the JWT : %s", r.RemoteAddr, err.Error())
			return
		}
		message.Id = uid

		if message.ReceiverId == "" && message.GroupId == "" {
			nw.Error("There is no recipient") // Handle JWT decryption error
			log.Printf("[%s] [AddMessage] There is no recipient", r.RemoteAddr)
			return
		}

		if message.ReceiverId != "" && message.GroupId != "" {
			nw.Error("Is this a private message or a group message ?") // Handle JWT decryption error
			log.Printf("[%s] [AddMessage] Is this a private message or a group message ?", r.RemoteAddr)
			return
		}

		if err = utils.IfExistsInDB("UserInfo", db, map[string]any{"Id": message.ReceiverId}); message.ReceiverId != "" && err != nil {
			nw.Error("The receiver user didn't exist") // Handle JWT decryption error
			log.Printf("[%s] [AddMessage] The receiver user didn't exist: %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("Groups", db, map[string]any{"Id": message.GroupId}); message.GroupId != "" && err != nil {
			nw.Error("The receiver group didn't exist") // Handle JWT decryption error
			log.Printf("[%s] [AddMessage] The receiver group didn't exist: %v", r.RemoteAddr, err)
			return
		}

		messageId, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
			log.Printf("[%s] [AddMessage] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}

		message.Id = messageId.String()

		if err = message.InsertIntoDb(db); err != nil {
			nw.Error("Error during the push in the DB") // Handle JWT decryption error
			log.Printf("[%s] [AddMessage] Error during the push in the DB: %v", r.RemoteAddr, err)
			return
		}

		if message.ReceiverId != "" {

			notifId, err := uuid.NewV7()
			if err != nil {
				nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
				log.Printf("[%s] [AddMessage] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
				return
			}

			var userData model.Register
			if err = userData.SelectFromDb(db, map[string]any{"Id": message.ReceiverId}); err != nil {
				nw.Error("There is a problem during the fetching of the user") // Handle UUID generation error
				log.Printf("[%s] [CreateComment] There is a problem during the fetching of the user : %s", r.RemoteAddr, err)
				return
			}

			var userDataName string
			if userData.Username == "" {
				userDataName = userData.FirstName + " " + userData.LastName
			} else {
				userDataName = userData.Username
			}

			notification := model.Notification{
				Id:          notifId.String(),
				UserId:      message.ReceiverId,
				Status:      "Chat",
				Description: fmt.Sprintf("A message as been send by %s", userDataName),
				GroupId:     "",
				OtherUserId: message.SenderId,
			}

			if err = notification.InsertIntoDb(db); err != nil {
				nw.Error("There is a probleme during the sending of a notification")
				log.Printf("[%s] [CreateComment] There is a probleme during the sending of a notification : %s", r.RemoteAddr, err)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Message send successfully",
		})
		if err != nil {
			log.Printf("[%s] [AddMessage] %s", r.RemoteAddr, err.Error())
		}
	}
}
