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
			"Value":   notifications,
		})
		if err != nil {
			log.Printf("[%s] [GetAllNotifications] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetGroupNotification(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas struct {
			UserId  string `json:"UserId"`
			GroupId string `json:"GroupId"`
		}

		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetGroupNotification] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [GetGroupNotification] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted Organisator ID
		datas.UserId = decryptUserId

		var notifications model.Notifications
		if err = notifications.SelectFromDb(db, map[string]any{"UserId": datas.UserId, "GroupId": datas.GroupId}); err != nil {
			nw.Error("Error during the fetch of the DB") // Handle invalid JWT error
			log.Printf("[%s] [GetGroupNotification] Error during the fetch of the DB : %v", r.RemoteAddr, err)
			return
		}

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group notifications getted successfully",
			"Value":   notifications,
		})
		if err != nil {
			log.Printf("[%s] [GetGroupNotification] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetUserNotification(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas struct {
			UserId      string `json:"UserId"`
			OtherUserId string `json:"OtherUserId"`
		}
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetGroupNotification] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [GetGroupNotification] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted Organisator ID
		datas.UserId = decryptUserId

		var notifications model.Notifications
		if err = notifications.SelectFromDb(db, map[string]any{"UserId": datas.UserId, "OtherUserId": datas.OtherUserId}); err != nil {
			nw.Error("Error during the fetch of the DB") // Handle invalid JWT error
			log.Printf("[%s] [GetGroupNotification] Error during the fetch of the DB : %v", r.RemoteAddr, err)
			return
		}

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "User notifications getted successfully",
			"Value":   notifications,
		})
		if err != nil {
			log.Printf("[%s] [GetGroupNotification] %s", r.RemoteAddr, err.Error())
		}
	}
}

func DeleteAllNotifications(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var userId string
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&userId); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [DeleteAllNotifications] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptUserId, err := utils.DecryptJWT(userId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [DeleteAllNotifications] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted Organisator ID
		userId = decryptUserId

		var notification model.Notification
		if err = notification.DeleteFromDb(db, map[string]any{"UserId": userId}); err != nil {
			nw.Error("Error during the update of the DB") // Handle invalid JWT error
			log.Printf("[%s] [DeleteAllNotifications] Error during the update of the DB : %v", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[userId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				Description string
			}

			WebsocketMessage.Type = "DeleteAllNotification"
			WebsocketMessage.Description = "All notification have been removed"
			
			if err = model.ConnectedWebSocket.Conn[userId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket") // Handle invalid JWT error
				log.Printf("[%s] [DeleteAllNotifications] Error during the communication with the websocket : %v", r.RemoteAddr, err)
				return
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "User notifications deleted successfully",
		})
		if err != nil {
			log.Printf("[%s] [DeleteAllNotifications] %s", r.RemoteAddr, err.Error())
		}
	}
}

func DeleteAllGroupNotifications(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas struct {
			UserId  string `json:"UserId"`
			GroupId string `json:"GroupId"`
		}
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [DeleteAllGroupNotifications] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [DeleteAllGroupNotifications] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted Organisator ID
		datas.UserId = decryptUserId

		var notification model.Notification
		if err = notification.DeleteFromDb(db, map[string]any{"UserId": datas.UserId, "GroupId": datas.GroupId}); err != nil {
			nw.Error("Error during the update of the DB") // Handle invalid JWT error
			log.Printf("[%s] [DeleteAllGroupNotifications] Error during the update of the DB : %v", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[datas.UserId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				GroupId     string
				Description string
			}

			WebsocketMessage.Type = "DeleteGroupNotification"
			WebsocketMessage.GroupId = datas.GroupId
			WebsocketMessage.Description = "All the notification of a group have been removed"

			if err = model.ConnectedWebSocket.Conn[datas.UserId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket") // Handle invalid JWT error
				log.Printf("[%s] [DeleteAllGroupNotifications] Error during the communication with the websocket : %v", r.RemoteAddr, err)
				return
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group notifications deleted successfully",
		})
		if err != nil {
			log.Printf("[%s] [DeleteAllGroupNotifications] %s", r.RemoteAddr, err.Error())
		}
	}
}

func DeleteAllUserNotifications(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas struct {
			UserId      string `json:"UserId"`
			OtherUserId string `json:"OtherUserId"`
		}
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [DeleteAllUserNotifications] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [DeleteAllUserNotifications] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted Organisator ID
		datas.UserId = decryptUserId

		var notification model.Notification
		if err = notification.DeleteFromDb(db, map[string]any{"UserId": datas.UserId, "OtherUserId": datas.OtherUserId}); err != nil {
			nw.Error("Error during the update of the DB") // Handle invalid JWT error
			log.Printf("[%s] [DeleteAllUserNotifications] Error during the update of the DB : %v", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[datas.UserId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				UserId      string
				Description string
			}

			WebsocketMessage.Type = "DeleteUserNotification"
			WebsocketMessage.UserId = datas.OtherUserId
			WebsocketMessage.Description = "All the notification with a user have been removed"

			if err = model.ConnectedWebSocket.Conn[datas.UserId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket") // Handle invalid JWT error
				log.Printf("[%s] [DeleteAllUserNotifications] Error during the communication with the websocket : %v", r.RemoteAddr, err)
				return
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "User notifications deleted successfully",
		})
		if err != nil {
			log.Printf("[%s] [DeleteAllUserNotifications] %s", r.RemoteAddr, err.Error())
		}
	}
}
