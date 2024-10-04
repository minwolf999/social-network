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

func CreateGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var group model.Group
		if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [CreateGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We decrypt the Id of the user make the request to follow someone
		decryptAuthorId, err := utils.DecryptJWT(group.LeaderId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [CreateGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		group.LeaderId = decryptAuthorId
		group.MemberIds = decryptAuthorId

		if group.GroupName == "" || group.CreationDate == "" || group.LeaderId == "" {
			nw.Error("There is no id for the user to follow")
			log.Printf("[%s] [CreateGroup] There is no id for the user to follow", r.RemoteAddr)
			return
		}

		uuid, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a probleme with the generation of the uuid")
			log.Printf("[%s] [CreateGroup] There is a probleme with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}
		group.Id = uuid.String()

		// We look if the 2 Ids exists in the Db
		if err = utils.IfNotExistsInDB("Groups", db, map[string]any{"GroupName": group.GroupName}); err != nil {
			nw.Error("There is already a group with the name : " + group.GroupName)
			log.Printf("[%s] [CreateGroup] %s", r.RemoteAddr, err)
			return
		}

		// We insert in the table Follower of the db the structure created
		if err := utils.InsertIntoDb("Groups", db, group.Id, group.LeaderId, group.MemberIds, group.GroupName, group.CreationDate); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [CreateGroup] %s", r.RemoteAddr, err.Error())
			return
		}

		// We send a success response to the request
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group created successfully",
		})
		if err != nil {
			log.Printf("[%s] [CreateGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func JoinAndLeaveGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var datas struct {
			UserId  string `json:"UserId"`
			GroupId string `json:"GroupId"`

			JoinOrLeave string `json:"JoinOrLeave"`
		}
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [JoinAndLeaveGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We decrypt the Id of the user make the request to follow someone
		decryptAuthorId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [JoinAndLeaveGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		datas.UserId = decryptAuthorId

		groupDatas, err := utils.SelectFromDb("Groups", db, map[string]any{"Id": datas.GroupId})
		if err != nil {
			nw.Error("Internal error: Problem during database query")
			log.Printf("[%s] [JoinAndLeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		group, err := utils.ParseGroupData(groupDatas)
		if err != nil {
			nw.Error("Internal Error: There is a probleme during the parse of the structure : " + err.Error())
			log.Printf("[%s] [JoinAndLeaveGroup] %s", r.RemoteAddr, err.Error())
			return
		}

		if len(group) != 1 {
			nw.Error("Internal Error: There is no group with this id")
			log.Printf("[%s] [JoinAndLeaveGroup] There is no group with this id", r.RemoteAddr)
			return
		}

		datas.JoinOrLeave = strings.ToLower(datas.JoinOrLeave)

		if datas.JoinOrLeave != "join" && datas.JoinOrLeave != "leave" {
			nw.Error("Internal Error: You can only join or leave a group")
			log.Printf("[%s] [JoinAndLeaveGroup] You can only join or leave a group", r.RemoteAddr)
			return
		}

		DetailGroup := group[0]
		DetailGroup.SplitMembers()

		if datas.JoinOrLeave == "join" {
			if slices.Contains(DetailGroup.SplitMemberIds, datas.UserId) {
				nw.Error("You are already in this group")
				log.Printf("[%s] [JoinAndLeaveGroup] You are already in this group", r.RemoteAddr)
				return
			}

			DetailGroup.SplitMemberIds = append(DetailGroup.SplitMemberIds, datas.UserId)
		} else {
			for i := range DetailGroup.SplitMemberIds {
				if DetailGroup.SplitMemberIds[i] == datas.UserId {
					if i < len(DetailGroup.SplitMemberIds)-1 {
						DetailGroup.SplitMemberIds = append(DetailGroup.SplitMemberIds[:i], DetailGroup.SplitMemberIds[i+1:]...)
					} else {
						DetailGroup.SplitMemberIds = DetailGroup.SplitMemberIds[:i]
					}

					break
				}
			}

			DetailGroup.LeaderId = DetailGroup.SplitMemberIds[0]
		}

		DetailGroup.JoinMembers()

		if err = utils.UpdateDb("Groups", db, map[string]any{"MemberIds": DetailGroup.MemberIds}, map[string]any{"Id": DetailGroup.Id}); err != nil {
			nw.Error("Internal error: Problem during database update")
			log.Printf("[%s] [JoinAndLeaveGroup] %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group joined successfuly",
		})
		if err != nil {
			log.Printf("[%s] [JoinAndLeaveGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var datas struct {
			UserId  string `json:"UserId"`
			GroupId string `json:"GroupId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We decrypt the Id of the user make the request to follow someone
		decryptAuthorId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		datas.UserId = decryptAuthorId

		groupDatas, err := utils.SelectFromDb("Groups", db, map[string]any{"Id": datas.GroupId})
		if err != nil {
			nw.Error("Internal error: Problem during database query")
			log.Printf("[%s] [GetGroup] %v", r.RemoteAddr, err)
			return
		}

		group, err := utils.ParseGroupData(groupDatas)
		if err != nil {
			nw.Error("Internal Error: There is a probleme during the parse of the structure : " + err.Error())
			log.Printf("[%s] [GetGroup] %s", r.RemoteAddr, err.Error())
			return
		}

		if len(group) != 1 {
			nw.Error("Internal Error: There is no group with this id")
			log.Printf("[%s] [GetGroup] There is no group with this id", r.RemoteAddr)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group obtained successfully",
			"Group": group[0],
		})
		if err != nil {
			log.Printf("[%s] [GetGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}

func DeleteGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var datas struct {
			UserId  string `json:"UserId"`
			GroupId string `json:"GroupId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&datas); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [DeleteGroup] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We decrypt the Id of the user make the request to follow someone
		decryptAuthorId, err := utils.DecryptJWT(datas.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [DeleteGroup] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		datas.UserId = decryptAuthorId

		groupDatas, err := utils.SelectFromDb("Groups", db, map[string]any{"Id": datas.GroupId})
		if err != nil {
			nw.Error("Internal error: Problem during database query")
			log.Printf("[%s] [DeleteGroup] %v", r.RemoteAddr, err)
			return
		}

		group, err := utils.ParseGroupData(groupDatas)
		if err != nil {
			nw.Error("Internal Error: There is a probleme during the parse of the structure : " + err.Error())
			log.Printf("[%s] [DeleteGroup] %s", r.RemoteAddr, err.Error())
			return
		}

		if len(group) != 1 {
			nw.Error("Internal Error: There is no group with this id")
			log.Printf("[%s] [DeleteGroup] There is no group with this id", r.RemoteAddr)
			return
		}

		if group[0].LeaderId != datas.UserId {
			nw.Error("You can't delete this group")
			log.Printf("[%s] [DeleteGroup] You can't delete this group", r.RemoteAddr)
			return
		}

		if err = utils.RemoveFromDB("Groups", db, map[string]any{"Id": datas.GroupId}); err != nil {
			nw.Error("Error during the remove of the db")
			log.Printf("[%s] [DeleteGroup] Error during the remove in the db: %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Group delete successfully",
		})
		if err != nil {
			log.Printf("[%s] [DeleteGroup] %s", r.RemoteAddr, err.Error())
		}
	}
}