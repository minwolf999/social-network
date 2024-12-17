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

		uid, err := utils.DecryptJWT(message.SenderId, db)
		if err != nil {
			nw.Error("Error when decrypt the JWT") // Handle JWT decryption error
			log.Printf("[%s] [AddMessage] Error when decrypt the JWT : %s", r.RemoteAddr, err.Error())
			return
		}
		message.SenderId = uid

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

		if err = utils.IfExistsInDB("UserInfo", db, map[string]any{"Id": message.SenderId}); message.ReceiverId != "" && err != nil {
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
			if err = userData.SelectFromDb(db, map[string]any{"Id": message.SenderId}); err != nil {
				nw.Error("There is a problem during the fetching of the user") // Handle UUID generation error
				log.Printf("[%s] [AddMessage] There is a problem during the fetching of the user : %s", r.RemoteAddr, err)
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
				log.Printf("[%s] [AddMessage] There is a probleme during the sending of a notification : %s", r.RemoteAddr, err)
				return
			}

			model.ConnectedWebSocket.Mu.Lock()
			_, isOk := model.ConnectedWebSocket.Conn[message.SenderId]
			if isOk {
				var WebsocketMessage struct {
					Type         string
					Sender       string
					Description  string
					Value        model.Message
					Notification model.Notification
				}

				WebsocketMessage.Type = "Private Chat"
				WebsocketMessage.Sender = message.SenderId
				WebsocketMessage.Description = "A private message have been send"
				WebsocketMessage.Value = message
				WebsocketMessage.Notification = notification

				if err = model.ConnectedWebSocket.Conn[message.SenderId].WriteJSON(WebsocketMessage); err != nil {

					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [AddMessage] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}

				_, isOk2 := model.ConnectedWebSocket.Conn[message.ReceiverId]
				if isOk2 {
					if err = model.ConnectedWebSocket.Conn[message.ReceiverId].WriteJSON(WebsocketMessage); err != nil {

						nw.Error("Error during the communication with the websocket")
						log.Printf("[%s] [AddMessage] Error during the communication with the websocket : %s", r.RemoteAddr, err)
						return
					}
				}
			}
			model.ConnectedWebSocket.Mu.Unlock()

		} else {
			var group model.Group
			if err = group.SelectFromDb(db, map[string]any{"Id": message.GroupId}); err != nil {
				nw.Error("There is a problem during the fetching of the group") // Handle UUID generation error
				log.Printf("[%s] [AddMessage] There is a problem during the fetching of the group : %s", r.RemoteAddr, err)
				return
			}

			group.SplitMembers()

			var userData model.Register
			if err = userData.SelectFromDb(db, map[string]any{"Id": message.SenderId}); err != nil {
				nw.Error("There is a problem during the fetching of the user") // Handle UUID generation error
				log.Printf("[%s] [AddMessage] There is a problem during the fetching of the user : %s", r.RemoteAddr, err)
				return
			}

			var userDataName string
			if userData.Username == "" {
				userDataName = userData.FirstName + " " + userData.LastName
			} else {
				userDataName = userData.Username
			}

			for i := range group.SplitMemberIds {
				notifId, err := uuid.NewV7()
				if err != nil {
					nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
					log.Printf("[%s] [AddMessage] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
					return
				}

				notification := model.Notification{
					Id:          notifId.String(),
					UserId:      group.SplitMemberIds[i],
					Status:      "Chat",
					Description: fmt.Sprintf("A message as been send by %s", userDataName),
					GroupId:     message.GroupId,
					OtherUserId: "",
				}

				if message.SenderId != group.SplitMemberIds[i] {
					if err = notification.InsertIntoDb(db); err != nil {
						nw.Error("There is a probleme during the sending of a notification")
						log.Printf("[%s] [AddMessage] There is a probleme during the sending of a notification : %s", r.RemoteAddr, err)
						return
					}
				}

				model.ConnectedWebSocket.Mu.Lock()
				_, isOk := model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]]
				if isOk {
					var WebsocketMessage struct {
						Type         string
						Sender       string
						Description  string
						GroupId      string
						Value        model.Message
						Notification model.Notification
					}

					WebsocketMessage.Type = "Group Chat"
					WebsocketMessage.Sender = message.SenderId
					WebsocketMessage.Description = "A group message have been send"
					WebsocketMessage.GroupId = message.GroupId
					WebsocketMessage.Value = message
					WebsocketMessage.Notification = notification

					if err = model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]].WriteJSON(WebsocketMessage); err != nil {

						nw.Error("Error during the communication with the websocket")
						log.Printf("[%s] [AddMessage] Error during the communication with the websocket : %s", r.RemoteAddr, err)
						return
					}
				}
				model.ConnectedWebSocket.Mu.Unlock()

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

func GetMessage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var message model.Message

		// Decode the JSON request body into the request struct
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [GetMessage] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		uid, err := utils.DecryptJWT(message.SenderId, db)
		if err != nil {
			nw.Error("Error when decrypt the JWT") // Handle JWT decryption error
			log.Printf("[%s] [GetMessage] Error when decrypt the JWT : %s", r.RemoteAddr, err.Error())
			return
		}
		message.SenderId = uid

		if message.ReceiverId == "" && message.GroupId == "" {
			nw.Error("There is no recipient") // Handle JWT decryption error
			log.Printf("[%s] [GetMessage] There is no recipient", r.RemoteAddr)
			return
		}

		if message.ReceiverId != "" && message.GroupId != "" {
			nw.Error("Is this a private message or a group message ?") // Handle JWT decryption error
			log.Printf("[%s] [GetMessage] Is this a private message or a group message ?", r.RemoteAddr)
			return
		}

		var messages model.Messages
		if message.GroupId != "" {
			if err = messages.GetGroupsMessages(db, message); err != nil {
				nw.Error("Error during the fetch of the group messages") // Handle JWT decryption error
				log.Printf("[%s] [GetMessage] Error during the fetch of the group messages: %v", r.RemoteAddr, err)
				return
			}
		} else {
			if err = messages.GetPrivateMessages(db, message); err != nil {
				nw.Error("Error during the fetch of the messages") // Handle JWT decryption error
				log.Printf("[%s] [GetMessage] Error during the fetch of the messages: %v", r.RemoteAddr, err)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Message retreived successfully",
			"Value":   messages,
		})
		if err != nil {
			log.Printf("[%s] [GetMessage] %s", r.RemoteAddr, err.Error())
		}
	}
}
