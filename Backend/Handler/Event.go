package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	model "social-network/Model"
	utils "social-network/Utils"

	"github.com/gofrs/uuid"
)

func CreateEvent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a Comment object to hold the decoded request body
		var event model.Event
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [CreateEvent] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptAuthorId, err := utils.DecryptJWT(event.OrganisatorId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [CreateEvent] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted Organisator ID
		event.OrganisatorId = decryptAuthorId

		// Check if the datas are not empty
		if event.DateOfTheEvent == "" || event.Description == "" || event.GroupId == "" || event.OrganisatorId == "" || event.Title == "" {
			nw.Error("There is an empty field")
			log.Printf("[%s] [CreateEvent] There is no text for the comment", r.RemoteAddr)
			return
		}

		var group model.Group
		if err = group.SelectFromDb(db, map[string]any{"Id": event.GroupId}); err != nil || group.Id == "" {
			nw.Error("The groupId given correspond to no group")
			log.Printf("[%s] [CreateEvent] The groupId given correspond to no group : %s", r.RemoteAddr, err)
			return
		}

		group.SplitMembers()
		if !slices.Contains(group.SplitMemberIds, event.OrganisatorId) {
			nw.Error("You are not in this group")
			log.Printf("[%s] [CreateEvent] You are not in this group", r.RemoteAddr)
			return
		}

		// Generate a new UUID for the event
		uid, err := uuid.NewV7()
		if err != nil {
			// Handle UUID generation error
			nw.Error("There is a problem with the generation of the uuid")
			log.Printf("[%s] [CreateEvent] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}
		event.Id = uid.String()

		// Insert the new event in the db
		if err = event.InsertIntoDb(db); err != nil {
			nw.Error("There is an error during the push in the db")
			log.Printf("[%s] [CreateEvent] There is an error during the push in the db: %v", r.RemoteAddr, err)
			return
		}

		group.SplitMembers()

		model.ConnectedWebSocket.Mu.Lock()
		for i := range group.SplitMemberIds {
			notifId, err := uuid.NewV7()
			if err != nil {
				nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
				log.Printf("[%s] [CreateEvent] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
				return
			}

			var userData model.Register
			if err = userData.SelectFromDb(db, map[string]any{"Id": event.OrganisatorId}); err != nil {
				nw.Error("There is a problem during the fetching of the user") // Handle UUID generation error
				log.Printf("[%s] [CreateEvent] There is a problem during the fetching of the user : %s", r.RemoteAddr, err)
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
				UserId:      group.SplitMemberIds[i],
				Status:      "Event",
				Description: fmt.Sprintf("An Event \"%s\" as been posted by %s for the group %s", event.Title, userDataName, group.GroupName),
				GroupId:     group.Id,
				OtherUserId: "",
			}

			if err = notification.InsertIntoDb(db); err != nil {
				nw.Error("There is a probleme during the sending of a notification")
				log.Printf("[%s] [CreateEvent] There is a probleme during the sending of a notification : %s", r.RemoteAddr, err)
				return
			}

			_, isOk := model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]]
			if isOk {
				var WebsocketMessage struct {
					Type        string
					GroupId     string
					Description string
					Value       model.Event
				}

				WebsocketMessage.Type = "Event"
				WebsocketMessage.GroupId = event.GroupId
				WebsocketMessage.Description = "An event have created for the group"
				WebsocketMessage.Value = event

				if err = model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]].WriteJSON(WebsocketMessage); err != nil {

					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [CreateEvent] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Event created successfully",
		})
		if err != nil {
			log.Printf("[%s] [CreateEvent] %s", r.RemoteAddr, err.Error())
		}
	}
}

func JoinEvent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a Comment object to hold the decoded request body
		var joinEvent model.JoinEvent
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&joinEvent); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [JoinEvent] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptAuthorId, err := utils.DecryptJWT(joinEvent.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [JoinEvent] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted User ID
		joinEvent.UserId = decryptAuthorId

		if err = utils.IfExistsInDB("Event", db, map[string]any{"Id": joinEvent.EventId}); err != nil {
			nw.Error("Invalid event id")
			log.Printf("[%s] [JoinEvent] Invalid event id : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfNotExistsInDB("JoinEvent", db, map[string]any{"EventId": joinEvent.EventId, "UserId": joinEvent.UserId}); err != nil {
			nw.Error("Event already joined")
			log.Printf("[%s] [JoinEvent] Event already joined : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfNotExistsInDB("DeclineEvent", db, map[string]any{"EventId": joinEvent.EventId, "UserId": joinEvent.UserId}); err != nil {
			var declineEvent = model.DeclineEvent(joinEvent)

			if err = declineEvent.DeleteFromDb(db, map[string]any{"EventId": declineEvent.EventId, "UserId": declineEvent.UserId}); err != nil {
				nw.Error("Error during the remove of the previous decline")
				log.Printf("[%s] [JoinEvent] Error during the remove of the previous decline: %v", r.RemoteAddr, err)
			}
		}

		if err = joinEvent.InsertIntoDb(db); err != nil {
			nw.Error("Impossible to insert in the db")
			log.Printf("[%s] [JoinEvent] Impossible to insert in the db : %v", r.RemoteAddr, err)
			return
		}

		var declineEvent model.DeclineEvent
		if err = declineEvent.DeleteFromDb(db, map[string]any{"EventId": joinEvent.EventId, "UserId": joinEvent.UserId}); err != nil {
			nw.Error("Error during the delete of the previous decline for this event")
			log.Printf("[%s] [DeclineEvent] Error during the delete of the previous decline for this event : %v", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[joinEvent.UserId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				Description string
				EventId     string
			}

			WebsocketMessage.Type = "Join Event"
			WebsocketMessage.Description = "An event have been accepted"
			WebsocketMessage.EventId = joinEvent.EventId

			if err = model.ConnectedWebSocket.Conn[joinEvent.UserId].WriteJSON(WebsocketMessage); err != nil {
				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [CreateEvent] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Event joined successfully",
		})
		if err != nil {
			log.Printf("[%s] [JoinEvent] %s", r.RemoteAddr, err.Error())
		}
	}
}

func DeclineEvent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a DeclineEvent object to hold the decoded request body
		var declineEvent model.DeclineEvent
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&declineEvent); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [DeclineEvent] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptAuthorId, err := utils.DecryptJWT(declineEvent.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [DeclineEvent] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted User ID
		declineEvent.UserId = decryptAuthorId

		if err = utils.IfExistsInDB("Event", db, map[string]any{"Id": declineEvent.EventId}); err != nil {
			nw.Error("Invalid event id")
			log.Printf("[%s] [DeclineEvent] Invalid event id : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfNotExistsInDB("JoinEvent", db, map[string]any{"EventId": declineEvent.EventId, "UserId": declineEvent.UserId}); err != nil {
			var joinEvent = model.JoinEvent(declineEvent)

			if err = declineEvent.DeleteFromDb(db, map[string]any{"EventId": joinEvent.EventId, "UserId": joinEvent.UserId}); err != nil {
				nw.Error("Error during the remove of the previous decline")
				log.Printf("[%s] [DeclineEvent] Error during the remove of the previous decline: %v", r.RemoteAddr, err)
			}
		}

		if err = utils.IfNotExistsInDB("DeclineEvent", db, map[string]any{"EventId": declineEvent.EventId, "UserId": declineEvent.UserId}); err != nil {
			nw.Error("Event already declined")
			log.Printf("[%s] [DeclineEvent] Event already declined : %v", r.RemoteAddr, err)
			return
		}

		if err = declineEvent.InsertIntoDb(db); err != nil {
			nw.Error("Impossible to insert in the db")
			log.Printf("[%s] [DeclineEvent] Impossible to insert in the db : %v", r.RemoteAddr, err)
			return
		}

		var acceptEvent model.JoinEvent
		if err = acceptEvent.DeleteFromDb(db, map[string]any{"EventId": declineEvent.EventId, "UserId": declineEvent.UserId}); err != nil {
			nw.Error("Error during the delete of the previous accept for this event")
			log.Printf("[%s] [DeclineEvent] Error during the delete of the previous accept for this event : %v", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[declineEvent.UserId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				Description string
				EventId     string
			}

			WebsocketMessage.Type = "Decline Event"
			WebsocketMessage.Description = "An event have been decline"
			WebsocketMessage.EventId = declineEvent.EventId

			if err = model.ConnectedWebSocket.Conn[declineEvent.UserId].WriteJSON(WebsocketMessage); err != nil {
				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [CreateEvent] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Event declined successfully",
		})
		if err != nil {
			log.Printf("[%s] [DeclineEvent] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetEvent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a DeclineEvent object to hold the decoded request body
		var event struct {
			UserId  string `json:"UserId"`
			EventId string `json:"EventId"`
		}
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetEvent] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptAuthorId, err := utils.DecryptJWT(event.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetEvent] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted User ID
		event.UserId = decryptAuthorId

		if err = utils.IfExistsInDB("Event", db, map[string]any{"Id": event.EventId}); err != nil {
			nw.Error("The event Id given didn't exist")
			log.Printf("[%s] [GetEvent] The event Id given didn't exist : %v", r.RemoteAddr, err)
			return
		}

		var tmp model.Event
		eventDetail, err := tmp.SelectFromDb(db, map[string]any{"Id": event.EventId})
		if err != nil {
			nw.Error("There is an error during the select")
			log.Printf("[%s] [GetEvent] There is an error during the select : %v", r.RemoteAddr, err)
			return
		}

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Event detail get successfully",
			"Value":   eventDetail,
		})
		if err != nil {
			log.Printf("[%s] [GetEvent] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetAllGroupEvents(db *sql.DB) http.HandlerFunc {
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
			log.Printf("[%s] [GetAllGroupEvents] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		_, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetAllGroupEvents] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		var event model.Event
		events, err := event.SelectFromDb(db, map[string]any{"GroupId": datas.GroupId})
		if err != nil {
			nw.Error("Error during the fetching of the DB")
			log.Printf("[%s] [GetAllGroupEvents] Error during the fetching of the DB : %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Events getted successfully",
			"Value":   events,
		})
		if err != nil {
			log.Printf("[%s] [GetEvent] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetJoinedEvent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas struct {
			UserId string `json:"UserId"`
		}

		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetJoinedEvent] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		jwt, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetJoinedEvent] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.UserId = jwt

		var event model.Event
		events, err := event.SelectFromDb(db, map[string]any{})
		if err != nil {
			nw.Error("Error during the fetch of the DB")
			log.Printf("[%s] [GetJoinedEvent] Error during the fetch of the DB : %v", r.RemoteAddr, err)
			return
		}

		for i := 0; i < len(events); i++ {
			if !strings.Contains(events[i].JoinUsers, datas.UserId) {
				if i < len(events)-1 {
					events = append(events[:i], events[i+1:]...)
				} else {
					events = events[:i]
				}

				i--
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Events getted successfully",
			"Value":   events,
		})
		if err != nil {
			log.Printf("[%s] [GetJoinedEvent] %s", r.RemoteAddr, err.Error())
		}
	}
}
