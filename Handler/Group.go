package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"

	model "social-network/Model"
	utils "social-network/Utils"

	"github.com/gofrs/uuid"
)

/*
This function takes 1 argument:
  - a pointer to an SQL database object

The purpose of this function is to handle the creation of a group, processing the request to create a new group in the database.

The function returns an http.HandlerFunc that can be used as a handler for HTTP requests.
*/
func CreateGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to handle errors and responses.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var group model.Group
		// Decode the JSON request body into the group struct.
		if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [CreateGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the LeaderId from the JWT to obtain the actual user ID.
		decryptAuthorId, err := utils.DecryptJWT(group.LeaderId, db)
		if err != nil {
			// Return error if the JWT is invalid.
			nw.Error("Invalid JWT")
			log.Printf("[%s] [CreateGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		// Assign the decrypted LeaderId back to the group.
		group.LeaderId = decryptAuthorId

		// Initialize MemberIds with the LeaderId.
		group.MemberIds = decryptAuthorId

		// Validate that required fields are provided.
		if group.GroupName == "" || group.CreationDate == "" || group.LeaderId == "" {
			// Return error if any required field is missing.
			nw.Error("There is an empty field")
			log.Printf("[%s] [CreateGroup] There is an empty field", r.RemoteAddr)
			return
		}

		// Generate a new UUID for the group.
		uuid, err := uuid.NewV7()
		if err != nil {
			// Return error if UUID generation fails.
			nw.Error("There is a problem with the generation of the uuid")
			log.Printf("[%s] [CreateGroup] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}

		// Set the generated UUID as the group's ID.
		group.Id = uuid.String()

		// Check if a group with the same name already exists in the database.
		if err = utils.IfNotExistsInDB("Groups", db, map[string]any{"GroupName": group.GroupName}); err != nil {
			// Return error if group name conflict occurs.
			nw.Error("There is already a group with the name : " + group.GroupName)
			log.Printf("[%s] [CreateGroup] %s", r.RemoteAddr, err)
			return
		}

		// Insert the new group into the database.
		if err := group.InsertIntoDb(db); err != nil {
			// Return error if database insertion fails.
			nw.Error("Internal Error: There is a problem during the push in the DB: " + err.Error())
			log.Printf("[%s] [CreateGroup] %s", r.RemoteAddr, err.Error())
			return
		}

		// Set the response header to indicate JSON content and respond with success message.
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group created successfully",

			// Include the ID of the newly created group in the response.
			"GroupId": group.Id,
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [CreateGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

/*
This function takes 1 argument:
  - a pointer to an SQL database object

The purpose of this function is to handle user requests to join or leave a group, updating the group's member list in the database accordingly.

The function returns an http.HandlerFunc that can be used as a handler for HTTP requests.
*/
func LeaveGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to handle errors and responses.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Struct to hold the data from the request body.
		var datas struct {
			// ID of the user joining or leaving the group.
			UserId string `json:"UserId"`
			// ID of the group.
			GroupId string `json:"GroupId"`
		}
		// Decode the JSON request body into the datas struct.
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [LeaveGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to obtain the actual user ID.
		decryptAuthorId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			// Return error if the JWT is invalid.
			nw.Error("Invalid JWT")
			log.Printf("[%s] [LeaveGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Assign the decrypted UserId back to the datas struct.
		datas.UserId = decryptAuthorId

		var user model.Register
		if err = user.SelectFromDb(db, map[string]any{"Id": datas.UserId}); err != nil {
			// Return error if there is a problem during the database query.
			nw.Error("Internal error: Problem during database query")
			log.Printf("[%s] [LeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		var group model.Group
		// Query the database for the group using the provided GroupId.
		if err = group.SelectFromDb(db, map[string]any{"Id": datas.GroupId}); err != nil {
			// Return error if there is a problem during the database query.
			nw.Error("Internal error: Problem during database query")
			log.Printf("[%s] [LeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		// Create a copy of the group data for modification.
		DetailGroup := group

		// Split the members into a manageable format.
		user.SplitGroups()
		DetailGroup.SplitMembers()

		index := slices.Index(DetailGroup.SplitMemberIds, datas.UserId)
		if index != -1 {
			if index < len(DetailGroup.SplitMemberIds)-1 {
				DetailGroup.SplitMemberIds = append(DetailGroup.SplitMemberIds[:index], DetailGroup.SplitMemberIds[index+1:]...)
			} else {
				DetailGroup.SplitMemberIds = DetailGroup.SplitMemberIds[:index]
			}
		}

		index = slices.Index(user.SplitGroupsJoined, datas.GroupId)
		if index != -1 {
			if index < len(user.SplitGroupsJoined)-1 {
				user.SplitGroupsJoined = append(user.SplitGroupsJoined[:index], user.SplitGroupsJoined[index+1:]...)
			} else {
				user.SplitGroupsJoined = user.SplitGroupsJoined[:index]
			}
		}

		if len(DetailGroup.SplitMemberIds) == 0 {
			if err = group.DeleteFromDb(db, map[string]any{"Id": group.Id}); err != nil {
				nw.Error("Internal error: Error during the delete of the group : " + err.Error())
				log.Printf("[%s] [LeaveGroup] Error during the delete of the group : %v", r.RemoteAddr, err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(map[string]any{
				"Success": true,
				"Message": "Group delete successfully",
			})
			if err != nil {
				// Log any error that occurs while encoding the response.
				log.Printf("[%s] [LeaveGroup] %s", r.RemoteAddr, err.Error())
			}

			return
		}

		// Update the LeaderId to the first member's ID after a user leaves.
		DetailGroup.LeaderId = DetailGroup.SplitMemberIds[0]

		// Update the member list format.
		user.JoinGroups()
		DetailGroup.JoinMembers()

		if err = user.UpdateDb(db, map[string]any{"GroupsJoined": user.GroupsJoined}, map[string]any{"Id": user.Id}); err != nil {
			// Return error if there is a problem during database update.
			nw.Error("Internal error: Problem during database update : " + err.Error())
			log.Printf("[%s] [LeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		// Update the group's member list in the database.
		if err = DetailGroup.UpdateDb(db, map[string]any{"LeaderId": DetailGroup.LeaderId, "MemberIds": DetailGroup.MemberIds}, map[string]any{"Id": DetailGroup.Id}); err != nil {
			// Return error if there is a problem during database update.
			nw.Error("Internal error: Problem during database update : " + err.Error())
			log.Printf("[%s] [LeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[datas.UserId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				GroupId     string
				Group       model.Group
				Description string
			}

			WebsocketMessage.Type = "LeaveGroup"
			WebsocketMessage.GroupId = group.Id
			WebsocketMessage.Description = "You leave the group"

			if err = group.SelectFromDb(db, map[string]any{"Id": group.Id}); err != nil {
				nw.Error("Error during the fetch of the group updated datas")
				log.Printf("[%s] [JoinGroup] Error during the fetch of the group updated datas : %s", r.RemoteAddr, err)
				return
			}
			group.SplitMembers()

			WebsocketMessage.Group = group

			if err = model.ConnectedWebSocket.Conn[datas.UserId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [JoinGroup] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Set the response header to indicate JSON content and respond with success message.
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,

			// Success message for joining the group.
			"Message": "Group leaved successfully",
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [LeaveGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

/*
This function takes 1 argument:
  - a pointer to an SQL database object

The purpose of this function is to handle user requests to retrieve a group’s information from the database.

The function returns an http.HandlerFunc that can be used as a handler for HTTP requests.
*/
func GetGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to handle errors and responses.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Struct to hold the data from the request body.
		var datas struct {
			// ID of the user making the request.
			UserId string `json:"UserId"`
			// ID of the group to retrieve.
			GroupId string `json:"GroupId"`
		}
		// Decode the JSON request body into the datas struct.
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to obtain the actual user ID.
		decryptAuthorId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			// Return error if the JWT is invalid.
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Assign the decrypted UserId back to the datas struct.
		datas.UserId = decryptAuthorId

		// Create a Group instance to hold the retrieved group data.
		var group = model.Group{Id: datas.GroupId}
		// Query the database for the group using the provided GroupId.
		err = group.SelectFromDb(db, map[string]any{"Id": datas.GroupId})
		if err != nil {
			// Return error if there is a problem during the database query.
			nw.Error("Internal error: Problem during database query")
			log.Printf("[%s] [GetGroup] %v", r.RemoteAddr, err)
			return
		}

		group.SplitMembers()

		// Set the response header to indicate JSON content and respond with the group data.
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			// Success message for retrieving the group.
			"Message": "Group obtained successfully",
			// The retrieved group data.
			"Group": group,
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [GetGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetAllGroups(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var userJWT string
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&userJWT); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetAllGroups] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		_, err := utils.DecryptJWT(userJWT, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetAllGroups] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		var groups model.Groups
		if err = groups.SelectFromDb(db, map[string]any{}); err != nil {
			nw.Error("Error during fetching the groups")
			log.Printf("[%s] [GetAllGroups] Error during fetching the groups : %v", r.RemoteAddr, err)
			return
		}

		// Set the response header to indicate JSON content and respond with success message.
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Groups getted successfully",

			"Groups": groups,
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [GetAllGroups] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetGroupsJoined(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var userJWT string
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&userJWT); err != nil {
			// Send error if decoding fails
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetGroupsJoined] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		userId, err := utils.DecryptJWT(userJWT, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetGroupsJoined] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		var groups model.Groups
		if err = groups.SelectFromDb(db, map[string]any{}); err != nil {
			nw.Error("Error during fetching the groups")
			log.Printf("[%s] [GetGroupsJoined] Error during fetching the groups : %v", r.RemoteAddr, err)
			return
		}

		for i := 0; i < len(groups); i++ {
			groups[i].SplitMembers()

			if !slices.Contains(groups[i].SplitMemberIds, userId) {
				if i < len(groups)-1 {
					groups = append(groups[:i], groups[i+1:]...)
				} else {
					groups = groups[:i]
				}
				i--
			} else {
				var notifications model.Notifications
				if err = notifications.SelectFromDb(db, map[string]any{"GroupId": groups[i].Id, "UserId": userId}); err != nil {
					nw.Error("Error during fetching the groups notifications quantity")
					log.Printf("[%s] [GetGroupsJoined] Error during fetching the groups notifications : %v", r.RemoteAddr, err)
					return
				}

				groups[i].NotificationQuantity = len(notifications)
			}
		}

		// Set the response header to indicate JSON content and respond with success message.
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Groups getted successfully",

			"Groups": groups,
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [GetGroupsJoined] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetGroupsPosts(db *sql.DB) http.HandlerFunc {
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
			log.Printf("[%s] [GetGroupsPosts] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the OrganisatorId from the JWT to get the actual Organisator ID
		_, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetGroupsPosts] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		if err := utils.IfExistsInDB("Groups", db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("Invalid group id")
			log.Printf("[%s] [GetGroupsPosts] Invalid group id: %v", r.RemoteAddr, err)
			return
		}

		var posts model.Posts
		if err = posts.SelectFromDb(db, map[string]any{"IsGroup": datas.GroupId}); err != nil {
			nw.Error("Error during the fetch of the DB")
			log.Printf("[%s] [GetGroupsPosts] Error during the fetch of the DB: %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group's posts getted successfully",

			"Posts": posts,
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [GetGroupsPosts] %s", r.RemoteAddr, err.Error())
		}
	}
}

/*
This function takes 1 argument:
  - a pointer to an SQL database object

The purpose of this function is to handle user requests to delete a group from the database.

The function returns an http.HandlerFunc that can be used as a handler for HTTP requests.
*/
func DeleteGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to handle errors and responses.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Struct to hold the data from the request body.
		var datas struct {
			// ID of the user making the request.
			UserId string `json:"UserId"`
			// ID of the group to be deleted.
			GroupId string `json:"GroupId"`
		}

		// Decode the JSON request body into the datas struct.
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [DeleteGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to obtain the actual user ID.
		decryptAuthorId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			// Return error if the JWT is invalid.
			nw.Error("Invalid JWT")
			log.Printf("[%s] [DeleteGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		// Assign the decrypted UserId back to the datas struct.
		datas.UserId = decryptAuthorId
		// Create a Group instance with the provided GroupId.
		var group = model.Group{Id: datas.GroupId}

		// Query the database to select the group by its ID.
		err = group.SelectFromDb(db, map[string]any{"Id": group.Id})
		if err != nil {
			// Return error if there is a problem during the database select operation.
			nw.Error("Error during the select in the db")
			log.Printf("[%s] [DeleteGroup] Error during the select or the parse: %v", r.RemoteAddr, err)
			return
		}

		// Check if the user is the leader of the group before allowing deletion.
		if group.LeaderId != datas.UserId {
			// Return error if the user is not the leader.
			nw.Error("You can't delete this group")
			log.Printf("[%s] [DeleteGroup] You can't delete this group", r.RemoteAddr)
			return
		}

		// Delete the group from the database.
		if err = group.DeleteFromDb(db, map[string]any{"Id": group.Id}); err != nil {
			// Return error if there is a problem during the database delete operation.
			nw.Error("Error during the remove of the db")
			log.Printf("[%s] [DeleteGroup] Error during the remove in the db: %v", r.RemoteAddr, err)
			return
		}

		// Set the response header to indicate JSON content and respond with success message.
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			// Success message for deleting the group.
			"Message": "Group deleted successfully",
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [DeleteGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func JoinGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas model.JoinGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [JoinGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [JoinGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.UserId = decryptUserId

		var group model.Group
		if err = group.SelectFromDb(db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("There is a problem during the fetch of the DB")
			log.Printf("[%s] [JoinGroup] There is a problem during the fetch of the DB : %v", r.RemoteAddr, err)
			return
		}

		if group.GroupName == "" {
			nw.Error("There is no group with this id")
			log.Printf("[%s] [JoinGroup] There is no group with this id : %v", r.RemoteAddr, err)
			return
		}

		group.SplitMembers()
		if slices.Contains(group.SplitMemberIds, datas.UserId) {
			nw.Error("You are already in the group")
			log.Printf("[%s] [JoinGroup] This user is already in the group : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfNotExistsInDB("JoinGroupRequest", db, map[string]any{"UserId": datas.UserId, "GroupId": datas.GroupId}); err != nil {
			nw.Error("You are already send a request to join this group")
			log.Printf("[%s] [JoinGroup] You are already send a request to join this group : %v", r.RemoteAddr, err)
			return
		}

		if err = datas.InsertIntoDb(db); err != nil {
			nw.Error("There is an error storing the query")
			log.Printf("[%s] [JoinGroup] There is an error storing the query : %v", r.RemoteAddr, err)
			return
		}

		notifId, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
			log.Printf("[%s] [JoinGroup] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}

		var userData model.Register
		if err = userData.SelectFromDb(db, map[string]any{"Id": datas.UserId}); err != nil {
			nw.Error("There is a problem during the fetching of the user") // Handle UUID generation error
			log.Printf("[%s] [JoinGroup] There is a problem during the fetching of the user : %s", r.RemoteAddr, err)
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
			UserId:      group.LeaderId,
			Status:      "Group",
			Description: fmt.Sprintf("An join request as been send to join the group \"%s\" by %s", group.GroupName, userDataName),
			GroupId:     group.Id,
			OtherUserId: "",
		}

		if err = notification.InsertIntoDb(db); err != nil {
			nw.Error("There is a probleme during the sending of a notification")
			log.Printf("[%s] [JoinGroup] There is a probleme during the sending of a notification : %s", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[group.LeaderId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				GroupId     string
				Description string
				Value       model.Group
				JoinRequest model.JoinGroupRequest
			}

			WebsocketMessage.Type = "JoinGroup"
			WebsocketMessage.GroupId = group.Id
			WebsocketMessage.Description = "A join request has been send to your group"
			WebsocketMessage.Value = group

			var tmp model.JoinGroupRequests
			if err = tmp.SelectFromDb(db, map[string]any{"UserId": datas.UserId, "GroupId": datas.GroupId}); err != nil {
				nw.Error("Error during the fetch of the request")
				log.Printf("[%s] [JoinGroup] Error during the fetch of the request : %s", r.RemoteAddr, err)
				return
			}

			WebsocketMessage.JoinRequest = tmp[0]

			if err = model.ConnectedWebSocket.Conn[group.LeaderId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [JoinGroup] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Join Request successfully send",
		})
		if err != nil {
			log.Printf("[%s] [JoinGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetJoinRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas model.JoinGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetJoinRequest] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetJoinRequest] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.UserId = decryptUserId

		if err = utils.IfExistsInDB("Groups", db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("There is no group with this id")
			log.Printf("[%s] [GetJoinRequest] There is no group with this id : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("Groups", db, map[string]any{"Id": datas.GroupId, "LeaderId": datas.UserId}); err != nil {
			nw.Error("The current user isn't the leader of this goup")
			log.Printf("[%s] [GetJoinRequest] The current user isn't the leader of this goup : %v", r.RemoteAddr, err)
			return
		}

		var requests model.JoinGroupRequests
		if err = requests.SelectFromDb(db, map[string]any{"GroupId": datas.GroupId}); err != nil {
			nw.Error("There is an error during the fetching of the DB")
			log.Printf("[%s] [GetJoinRequest] There is an error during the fetching of the DB : %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group join request getted successfully",
			"Value":   requests,
		})
		if err != nil {
			log.Printf("[%s] [GetJoinRequest] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetSendJoinRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas struct {
			UserId  string `json:"UserId"`
			GroupId string `json:"GroupId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetSendJoinRequest] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetSendJoinRequest] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.UserId = decryptUserId

		var requests model.JoinGroupRequests
		if err = requests.SelectFromDb(db, map[string]any{"UserId": datas.UserId, "GroupId": datas.GroupId}); err != nil {
			nw.Error("Error during the fetch of the DB")
			log.Printf("[%s] [GetSendJoinRequest] Error during the fetch of the DB : %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group join request getted successfully",
			"Value":   requests,
		})
		if err != nil {
			log.Printf("[%s] [GetSendJoinRequest] %s", r.RemoteAddr, err.Error())
		}
	}
}

func DeclineJoinRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas struct {
			UserId     string `json:"UserId"`
			GroupId    string `json:"GroupId"`
			JoinUserId string `json:"JoinUserId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [DeclineJoinRequest] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [DeclineJoinRequest] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.UserId = decryptUserId

		if err = utils.IfExistsInDB("Groups", db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("There is no group with this id")
			log.Printf("[%s] [DeclineJoinRequest] There is no group with this id : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("Groups", db, map[string]any{"Id": datas.GroupId, "LeaderId": datas.UserId}); err != nil {
			nw.Error("The current user isn't the leader of this goup")
			log.Printf("[%s] [DeclineJoinRequest] The current user isn't the leader of this goup : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("JoinGroupRequest", db, map[string]any{"UserId": datas.JoinUserId, "GroupId": datas.GroupId}); err != nil {
			nw.Error("There is no request to join the group")
			log.Printf("[%s] [DeclineJoinRequest] There is no request to join the group : %v", r.RemoteAddr, err)
			return
		}

		if err = model.RemoveFromDB("JoinGroupRequest", db, map[string]any{"UserId": datas.JoinUserId, "GroupId": datas.GroupId}); err != nil {
			nw.Error("there is an error during the delete of the request")
			log.Printf("[%s] [DeclineJoinRequest] There is an error during the delete of the request : %s", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[datas.JoinUserId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				GroupId     string
				UserId      string
				Description string
			}

			WebsocketMessage.Type = "DeclineJoinRequest"
			WebsocketMessage.GroupId = datas.GroupId
			WebsocketMessage.UserId = datas.JoinUserId
			WebsocketMessage.Description = "A join request has been send to your group"

			if err = model.ConnectedWebSocket.Conn[datas.JoinUserId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [JoinGroup] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Join request successfully denied",
		})
		if err != nil {
			log.Printf("[%s] [DeclineJoinRequest] %s", r.RemoteAddr, err.Error())
		}
	}
}

func AcceptJoinRequest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas struct {
			UserId     string `json:"UserId"`
			GroupId    string `json:"GroupId"`
			JoinUserId string `json:"JoinUserId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [AcceptJoinRequest] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [AcceptJoinRequest] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.UserId = decryptUserId

		if err = utils.IfExistsInDB("Groups", db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("There is no group with this id")
			log.Printf("[%s] [AcceptJoinRequest] There is no group with this id : %v", r.RemoteAddr, err)
			return
		}

		var group model.Group
		if err = group.SelectFromDb(db, map[string]any{"Id": datas.GroupId, "LeaderId": datas.UserId}); err != nil {
			nw.Error("There is an error during the fetch of the group data")
			log.Printf("[%s] [AcceptJoinRequest] There is an error during the fetch of the group data : %v", r.RemoteAddr, err)
			return
		}

		if group.GroupName == "" {
			nw.Error("The current user isn't the leader of this goup")
			log.Printf("[%s] [AcceptJoinRequest] The current user isn't the leader of this goup", r.RemoteAddr)
			return
		}

		if err = utils.IfExistsInDB("JoinGroupRequest", db, map[string]any{"UserId": datas.JoinUserId, "GroupId": datas.GroupId}); err != nil {
			nw.Error("There is no request to join the group")
			log.Printf("[%s] [AcceptJoinRequest] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		group.SplitMembers()
		group.SplitMemberIds = append(group.SplitMemberIds, datas.JoinUserId)
		group.JoinMembers()

		if err = group.UpdateDb(db, map[string]any{"MemberIds": group.MemberIds}, map[string]any{"Id": group.Id}); err != nil {
			nw.Error("There is an error during the update of the group data")
			log.Printf("[%s] [AcceptJoinRequest] There is an error during the update of the group data : %s", r.RemoteAddr, err)
			return
		}

		if err = model.RemoveFromDB("JoinGroupRequest", db, map[string]any{"UserId": datas.JoinUserId, "GroupId": datas.GroupId}); err != nil {
			nw.Error("There is an error during the delete of the request")
			log.Printf("[%s] [AcceptJoinRequest] There is an error during the delete of the request : %s", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[datas.JoinUserId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				GroupId     string
				Description string
			}

			WebsocketMessage.Type = "AcceptJoinRequest"
			WebsocketMessage.GroupId = datas.GroupId
			WebsocketMessage.Description = "A join request has been send to your group"

			if err = model.ConnectedWebSocket.Conn[datas.JoinUserId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [JoinGroup] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Join request successfully accepted",
		})
		if err != nil {
			log.Printf("[%s] [AcceptJoinRequest] %s", r.RemoteAddr, err.Error())
		}
	}
}

func InviteGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas model.InviteGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [InviteGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(datas.SenderId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [InviteGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.SenderId = decryptUserId

		var group model.Group
		if err = group.SelectFromDb(db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("There is a problem during the fetch of the DB")
			log.Printf("[%s] [InviteGroup] There is a problem during the fetch of the DB : %v", r.RemoteAddr, err)
			return
		}

		if group.GroupName == "" {
			nw.Error("There is no group with this id")
			log.Printf("[%s] [InviteGroup] There is no group with this id : %v", r.RemoteAddr, err)
			return
		}

		group.SplitMembers()
		if slices.Contains(group.SplitMemberIds, datas.ReceiverId) {
			nw.Error("This user is already in the group")
			log.Printf("[%s] [InviteGroup] This user is already in the group : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("UserInfo", db, map[string]any{"Id": datas.ReceiverId}); err != nil {
			nw.Error("There is no user with this id")
			log.Printf("[%s] [InviteGroup] There is no user with the id %s : %v", r.RemoteAddr, datas.ReceiverId, err)
			return
		}

		if err = utils.IfNotExistsInDB("InviteGroupRequest", db, map[string]any{"GroupId": datas.GroupId, "ReceiverId": datas.ReceiverId}); err != nil {
			nw.Error("This user has already receive an invitation")
			log.Printf("[%s] [InviteGroup] This user has already receive an invitation %s : %v", r.RemoteAddr, datas.ReceiverId, err)
			return
		}

		if err = datas.InsertIntoDb(db); err != nil {
			nw.Error("There is an error during the store of the invitation")
			log.Printf("[%s] [InviteGroup] There is an error during the store of the invitation : %v", r.RemoteAddr, err)
			return
		}

		notifId, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
			log.Printf("[%s] [InviteGroup] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}

		var userData model.Register
		if err = userData.SelectFromDb(db, map[string]any{"Id": datas.SenderId}); err != nil {
			nw.Error("There is a problem during the fetching of the user") // Handle UUID generation error
			log.Printf("[%s] [InviteGroup] There is a problem during the fetching of the user : %s", r.RemoteAddr, err)
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
			UserId:      datas.ReceiverId,
			Status:      "Group",
			Description: fmt.Sprintf("An invitation to join the group \"%s\" as been send by %s", group.GroupName, userDataName),
			GroupId:     group.Id,
			OtherUserId: "",
		}

		if err = notification.InsertIntoDb(db); err != nil {
			nw.Error("There is a probleme during the sending of a notification")
			log.Printf("[%s] [InviteGroup] There is a probleme during the sending of a notification : %s", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		_, isOk := model.ConnectedWebSocket.Conn[datas.ReceiverId]
		if isOk {
			var WebsocketMessage struct {
				Type        string
				GroupId     string
				Description string
				Value       model.Group
				Invite      model.InviteGroupRequest
			}

			WebsocketMessage.Type = "InviteGroup"
			WebsocketMessage.GroupId = group.Id
			WebsocketMessage.Description = "A invite request has been send to join a group"
			WebsocketMessage.Value = group

			var tmp model.InviteGroupRequests
			if err = tmp.SelectFromDb(db, map[string]any{"GroupId": datas.GroupId, "ReceiverId": datas.ReceiverId, "SenderId": datas.SenderId}); err != nil {
				nw.Error("Error during the fetch of the db")
				log.Printf("[%s] [InviteGroup] Error during the fetch of the db : %s", r.RemoteAddr, err)
				return
			}

			WebsocketMessage.Invite = tmp[0]

			if err = model.ConnectedWebSocket.Conn[datas.ReceiverId].WriteJSON(WebsocketMessage); err != nil {

				nw.Error("Error during the communication with the websocket")
				log.Printf("[%s] [InviteGroup] Error during the communication with the websocket : %s", r.RemoteAddr, err)
				return
			}
		}

		group.SplitMembers()
		for i := range group.SplitMemberIds {
			_, isOk := model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]]
			if isOk {
				var WebsocketMessage struct {
					Type        string
					UserId      string
					Description string
				}

				WebsocketMessage.Type = "InvitePeopleInGroup"
				WebsocketMessage.Description = "A invite request has been send to join a group"
				WebsocketMessage.UserId = datas.ReceiverId

				if err = model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]].WriteJSON(WebsocketMessage); err != nil {

					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [InviteGroup] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}
			}
		}

		model.ConnectedWebSocket.Mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Invitation request successfully send",
		})
		if err != nil {
			log.Printf("[%s] [InviteGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetInvitationGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var userId string
		if err := json.NewDecoder(r.Body).Decode(&userId); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetInvitationGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(userId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetInvitationGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		userId = decryptUserId

		var invitations model.InviteGroupRequests
		if err = invitations.SelectFromDb(db, map[string]any{"ReceiverId": userId}); err != nil {
			nw.Error("Error during the fetching of the DB")
			log.Printf("[%s] [GetInvitationGroup] Error during the fetching of the DB : %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Invitations successfully getted",
			"Value":   invitations,
		})
		if err != nil {
			log.Printf("[%s] [InviteGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetInvitationUserInGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas struct {
			UserId  string `json:"UserId"`
			GroupId string `json:"GroupId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetInvitationUserInGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetInvitationUserInGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.UserId = decryptUserId

		var group model.Group
		if err = group.SelectFromDb(db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("Error during the fetch of the group data")
			log.Printf("[%s] [GetInvitationUserInGroup] Error during the fetch of the group data : %v", r.RemoteAddr, err)
			return
		}

		group.SplitMembers()
		if !slices.Contains(group.SplitMemberIds, datas.UserId) {
			nw.Error("You are not in this group")
			log.Printf("[%s] [GetInvitationUserInGroup] You are not in this group", r.RemoteAddr)
			return
		}

		var invitations model.InviteGroupRequests
		if err = invitations.SelectFromDb(db, map[string]any{"GroupId": datas.GroupId}); err != nil {
			nw.Error("Error during the fetching of the DB")
			log.Printf("[%s] [GetInvitationUserInGroup] Error during the fetching of the DB : %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Invitations successfully getted",
			"Value":   invitations,
		})
		if err != nil {
			log.Printf("[%s] [GetInvitationUserInGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func DeclineInvitationGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas model.InviteGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [DeclineInvitationGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(datas.ReceiverId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [DeclineInvitationGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.ReceiverId = decryptUserId

		if err = utils.IfExistsInDB("Groups", db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("There is no group with this id")
			log.Printf("[%s] [DeclineInvitationGroup] There is no group with this id : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("InviteGroupRequest", db, map[string]any{"GroupId": datas.GroupId, "ReceiverId": datas.ReceiverId}); err != nil {
			nw.Error("This user hasn't received any invitation for this group")
			log.Printf("[%s] [DeclineInvitationGroup] This user hasn't received any invitation for this group %s : %v", r.RemoteAddr, datas.ReceiverId, err)
			return
		}

		if err = model.RemoveFromDB("InviteGroupRequest", db, map[string]any{"GroupId": datas.GroupId, "ReceiverId": datas.ReceiverId}); err != nil {
			nw.Error("There is an error during the delete of the invitation")
			log.Printf("[%s] [DeclineInvitationGroup] There is an error during the delete of the invitation : %s", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Invitations successfully denied",
		})
		if err != nil {
			log.Printf("[%s] [DeclineInvitationGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func AcceptInvitationGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var datas model.InviteGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [AcceptInvitationGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		decryptUserId, err := utils.DecryptJWT(datas.ReceiverId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [AcceptInvitationGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		datas.ReceiverId = decryptUserId

		var group model.Group
		if err = group.SelectFromDb(db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("There is a problem during the fetch of the DB")
			log.Printf("[%s] [AcceptInvitationGroup] There is a problem during the fetch of the DB : %v", r.RemoteAddr, err)
			return
		}

		if group.GroupName == "" {
			nw.Error("There is no group with this id")
			log.Printf("[%s] [AcceptInvitationGroup] There is no group with this id : %v", r.RemoteAddr, err)
			return
		}

		if err = utils.IfExistsInDB("InviteGroupRequest", db, map[string]any{"GroupId": datas.GroupId, "ReceiverId": datas.ReceiverId}); err != nil {
			nw.Error("This user hasn't received any invitation for this group")
			log.Printf("[%s] [AcceptInvitationGroup] This user hasn't received any invitation for this group %s : %v", r.RemoteAddr, datas.ReceiverId, err)
			return
		}

		group.SplitMembers()
		group.SplitMemberIds = append(group.SplitMemberIds, datas.ReceiverId)
		group.JoinMembers()

		if err = group.UpdateDb(db, map[string]any{"MemberIds": group.MemberIds}, map[string]any{"Id": group.Id}); err != nil {
			nw.Error("There is an error during the update of the group's data")
			log.Printf("[%s] [AcceptInvitationGroup] There is an error during the update of the group's data : %s", r.RemoteAddr, err)
			return
		}

		if err = model.RemoveFromDB("InviteGroupRequest", db, map[string]any{"GroupId": datas.GroupId, "ReceiverId": datas.ReceiverId}); err != nil {
			nw.Error("There is an error during the delete of the invitation")
			log.Printf("[%s] [AcceptInvitationGroup] There is an error during the delete of the invitation : %s", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		for i := range group.SplitMemberIds {
			_, isOk := model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]]
			if isOk {
				var WebsocketMessage struct {
					Type        string
					GroupId     string
					Description string
					Value       model.Group
				}

				WebsocketMessage.Type = "AcceptInviteGroup"
				WebsocketMessage.GroupId = group.Id
				WebsocketMessage.Description = "A invite request has been accepted"
				WebsocketMessage.Value = group

				if err = model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]].WriteJSON(WebsocketMessage); err != nil {

					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [AcceptInvitationGroup] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Invitations successfully denied",
		})
		if err != nil {
			log.Printf("[%s] [AcceptInvitationGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}
