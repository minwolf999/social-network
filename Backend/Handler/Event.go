package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"slices"

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
		uuid, err := uuid.NewV7()
		if err != nil {
			// Handle UUID generation error
			nw.Error("There is a problem with the generation of the uuid")
			log.Printf("[%s] [CreateEvent] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}
		event.Id = uuid.String()

		// Insert the new event in the db
		if err = event.InsertIntoDb(db); err != nil {
			nw.Error("There is an error during the push in the db")
			log.Printf("[%s] [CreateEvent] There is an error during the push in the db: %v", r.RemoteAddr, err)
			return
		}

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
			log.Printf("[%s] [JoinOrDeclineEvent] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptAuthorId, err := utils.DecryptJWT(joinEvent.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [JoinOrDeclineEvent] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted User ID
		joinEvent.UserId = decryptAuthorId

		if err = utils.IfNotExistsInDB("Event", db, map[string]any{"Id": joinEvent.EventId}); err != nil {
			nw.Error("Invalid event id")
			log.Printf("[%s] [JoinOrDeclineEvent] Invalid event id : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("JoinEvent", db, map[string]any{"EventId": joinEvent.EventId, "UserId": joinEvent.UserId}); err != nil {
			nw.Error("Event already joined")
			log.Printf("[%s] [JoinOrDeclineEvent] Event already joined : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("DeclineEvent", db, map[string]any{"EventId": joinEvent.EventId, "UserId": joinEvent.UserId}); err != nil {
			nw.Error("Event already declined")
			log.Printf("[%s] [JoinOrDeclineEvent] Event already declined : %v", r.RemoteAddr, err)
			return
		}

		if err = joinEvent.InsertIntoDb(db); err != nil {
			nw.Error("Impossible to insert in the db")
			log.Printf("[%s] [JoinOrDeclineEvent] Impossible to insert in the db : %v", r.RemoteAddr, err)
			return
		}

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Event joined successfully",
		})
		if err != nil {
			log.Printf("[%s] [JoinOrDeclineEvent] %s", r.RemoteAddr, err.Error())
		}
	}
}

func DeclineEvent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a Comment object to hold the decoded request body
		var joinEvent model.DeclineEvent
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&joinEvent); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [JoinOrDeclineEvent] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		decryptAuthorId, err := utils.DecryptJWT(joinEvent.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [JoinOrDeclineEvent] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Set the decrypted User ID
		joinEvent.UserId = decryptAuthorId

		if err = utils.IfNotExistsInDB("Event", db, map[string]any{"Id": joinEvent.EventId}); err != nil {
			nw.Error("Invalid event id")
			log.Printf("[%s] [JoinOrDeclineEvent] Invalid event id : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("JoinEvent", db, map[string]any{"EventId": joinEvent.EventId, "UserId": joinEvent.UserId}); err != nil {
			nw.Error("Event already joined")
			log.Printf("[%s] [JoinOrDeclineEvent] Event already joined : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("DeclineEvent", db, map[string]any{"EventId": joinEvent.EventId, "UserId": joinEvent.UserId}); err != nil {
			nw.Error("Event already declined")
			log.Printf("[%s] [JoinOrDeclineEvent] Event already declined : %v", r.RemoteAddr, err)
			return
		}

		if err = joinEvent.InsertIntoDb(db); err != nil {
			nw.Error("Impossible to insert in the db")
			log.Printf("[%s] [JoinOrDeclineEvent] Impossible to insert in the db : %v", r.RemoteAddr, err)
			return
		}

		// Send a success response in JSON format
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Event joined successfully",
		})
		if err != nil {
			log.Printf("[%s] [JoinOrDeclineEvent] %s", r.RemoteAddr, err.Error())
		}
	}
}