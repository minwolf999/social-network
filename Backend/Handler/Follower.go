package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	model "social-network/Model"
	utils "social-network/Utils"

	"github.com/gofrs/uuid"
)

func AddFollower(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var follower model.Follower
		json.Unmarshal(body, &follower)

		// We decrypt the Id of the user make the request to follow someone
		decryptAuthorId, err := utils.DecryptJWT(follower.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [AddFollower] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		follower.UserId = decryptAuthorId

		// We check if the Id of the people want"ed to follow has been forgotten
		if follower.FollowerId == "" {
			nw.Error("There is no id for the user to follow")
			log.Printf("[%s] [AddFollower] There is no id for the user to follow", r.RemoteAddr)
			return
		}

		// We create a UID for the following link
		uuid, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a probleme with the generation of the uuid")
			log.Printf("[%s] [AddFollower] There is a probleme with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}
		follower.Id = uuid.String()

		// We look if the 2 Ids exists in the Db
		if err = utils.IfExistsInDB("Auth", db, map[string]any{"Id": follower.UserId}); err != nil {
			nw.Error("There is no user with the id of the JWT")
			log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, "There is no user with the id of the JWT")
			return
		}

		if err = utils.IfExistsInDB("Auth", db, map[string]any{"Id": follower.FollowerId}); err != nil {
			nw.Error("There Id of the people you want to follow didn't exist")
			log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, "There Id of the people you want to follow didn't exist")
			return
		}

		// We insert in the table Follower of the db the structure created
		if err := utils.InsertIntoDb("Follower", db, follower.Id, follower.UserId, follower.FollowerId); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, err.Error())
			return
		}

		// We send a success response to the request
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Add follower successfully",
		})
		if err != nil {
			log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, err.Error())
		}
	}
}

func RemoveFollower(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var follower model.Follower
		json.Unmarshal(body, &follower)

		// We decrypt the Id of the user make the request to follow someone
		decryptAuthorId, err := utils.DecryptJWT(follower.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [RemoveFollower] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		follower.UserId = decryptAuthorId

		// We check if the Id of the people want"ed to follow has been forgotten
		if follower.FollowerId == "" {
			nw.Error("There is no id for the user to follow")
			log.Printf("[%s] [RemoveFollower] There is no id for the user to follow", r.RemoteAddr)
			return
		}

		if err = utils.RemoveFromDB("Follower", db, map[string]any{"UserId": follower.UserId, "FollowerId": follower.FollowerId}); err != nil {
			nw.Error("Internal Error: There is a probleme during the delete in the DB: " + err.Error())
			log.Printf("[%s] [RemoveFollower] %s", r.RemoteAddr, err.Error())
			return
		}

		// We send a success response to the request
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Add follower successfully",
		})
		if err != nil {
			log.Printf("[%s] [RemoveFollower] %s", r.RemoteAddr, err.Error())
		}
	}
}
