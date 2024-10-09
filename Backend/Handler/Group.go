package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

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
			nw.Error("There is no id for the user to follow")
			log.Printf("[%s] [CreateGroup] There is no id for the user to follow", r.RemoteAddr)
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
		if err := model.InsertIntoDb("Groups", db, group.Id, group.LeaderId, group.MemberIds, group.GroupName, group.CreationDate); err != nil {
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
func JoinAndLeaveGroup(db *sql.DB) http.HandlerFunc {
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
			// Action to either join or leave the group.
			JoinOrLeave string `json:"JoinOrLeave"`
		}
		// Decode the JSON request body into the datas struct.
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [JoinAndLeaveGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the UserId from the JWT to obtain the actual user ID.
		decryptAuthorId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			// Return error if the JWT is invalid.
			nw.Error("Invalid JWT")
			log.Printf("[%s] [JoinAndLeaveGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		// Assign the decrypted UserId back to the datas struct.
		datas.UserId = decryptAuthorId

		var user model.Register
		if err = user.SelectFromDb(db, map[string]any{"Id": datas.UserId}); err != nil {
			// Return error if there is a problem during the database query.
			nw.Error("Internal error: Problem during database query")
			log.Printf("[%s] [JoinAndLeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		var group model.Group
		// Query the database for the group using the provided GroupId.
		if err = group.SelectFromDb(db, map[string]any{"Id": datas.GroupId}); err != nil {
			// Return error if there is a problem during the database query.
			nw.Error("Internal error: Problem during database query")
			log.Printf("[%s] [JoinAndLeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		// Convert the JoinOrLeave action to lowercase for consistency.
		datas.JoinOrLeave = strings.ToLower(datas.JoinOrLeave)

		// Validate that the action is either "join" or "leave".
		if datas.JoinOrLeave != "join" && datas.JoinOrLeave != "leave" {
			// Return error if the action is not valid.
			nw.Error("Internal Error: You can only join or leave a group")
			log.Printf("[%s] [JoinAndLeaveGroup] You can only join or leave a group", r.RemoteAddr)
			return
		}

		// Create a copy of the group data for modification.
		DetailGroup := group

		// Split the members into a manageable format.
		user.SplitGroups()
		DetailGroup.SplitMembers()

		if datas.JoinOrLeave == "join" {
			// Check if the user is already a member of the group.
			if slices.Contains(DetailGroup.SplitMemberIds, datas.UserId) {
				// Return error if the user is already in the group.
				nw.Error("You are already in this group")
				log.Printf("[%s] [JoinAndLeaveGroup] You are already in this group", r.RemoteAddr)
				return
			}

			user.SplitGroupsJoined = append(user.SplitGroupsJoined, datas.GroupId)

			// Add the UserId to the group’s member list.
			DetailGroup.SplitMemberIds = append(DetailGroup.SplitMemberIds, datas.UserId)
		} else {
			// If the action is to leave, remove the UserId from the member list.
			// for i := range DetailGroup.SplitMemberIds {
			// 	if DetailGroup.SplitMemberIds[i] == datas.UserId {
			// 		// Remove the user from the list of members.
			// 		if i < len(DetailGroup.SplitMemberIds)-1 {
			// 			DetailGroup.SplitMemberIds = append(DetailGroup.SplitMemberIds[:i], DetailGroup.SplitMemberIds[i+1:]...)
			// 		} else {
			// 			DetailGroup.SplitMemberIds = DetailGroup.SplitMemberIds[:i]
			// 		}
			// 		break
			// 	}
			// }

			index := slices.Index(DetailGroup.SplitMemberIds, datas.UserId)
			if index < len(DetailGroup.SplitMemberIds)-1 {
				DetailGroup.SplitMemberIds = append(DetailGroup.SplitMemberIds[:index], DetailGroup.SplitMemberIds[index+1:]...)
			} else {
				DetailGroup.SplitMemberIds = DetailGroup.SplitMemberIds[:index]
			}

			index = slices.Index(user.SplitGroupsJoined, datas.GroupId)
			if index < len(user.SplitGroupsJoined)-1 {
				user.SplitGroupsJoined = append(user.SplitGroupsJoined[:index], user.SplitGroupsJoined[index+1:]...)
			} else {
				user.SplitGroupsJoined = user.SplitGroupsJoined[:index]
			}

			// Update the LeaderId to the first member's ID after a user leaves.
			DetailGroup.LeaderId = DetailGroup.SplitMemberIds[0]
		}

		// Update the member list format.
		user.JoinGroups()
		DetailGroup.JoinMembers()

		if err = user.UpdateDb(db, map[string]any{"GroupsJoined": user.GroupsJoined}, map[string]any{"Id": user.Id}); err != nil {
			// Return error if there is a problem during database update.
			nw.Error("Internal error: Problem during database update : " + err.Error())
			log.Printf("[%s] [JoinAndLeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		// Update the group's member list in the database.
		if err = DetailGroup.UpdateDb(db, map[string]any{"MemberIds": DetailGroup.MemberIds}, map[string]any{"Id": DetailGroup.Id}); err != nil {
			// Return error if there is a problem during database update.
			nw.Error("Internal error: Problem during database update : " + err.Error())
			log.Printf("[%s] [JoinAndLeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		// Set the response header to indicate JSON content and respond with success message.
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,

			// Success message for joining the group.
			"Message": "Group joined successfully",
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [JoinAndLeaveGroup] %s", r.RemoteAddr, err.Error())
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
