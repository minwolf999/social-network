package handler

import (
	"database/sql"
	"encoding/json"
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
		var follower model.Follower
		if err := json.NewDecoder(r.Body).Decode(&follower); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [AddFollower] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

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
			nw.Error("There is no user with the id of the JWT : " + follower.UserId)
			log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, "There is no user with the id of the JWT")
			return
		}

		if err = utils.IfExistsInDB("Auth", db, map[string]any{"Id": follower.FollowerId}); err != nil {
			nw.Error("There Id of the people you want to follow didn't exist")
			log.Printf("[%s] [AddFollower] %s", r.RemoteAddr, "There Id of the people you want to follow didn't exist")
			return
		}

		if err = utils.IfNotExistsInDB("Follower", db, map[string]any{"UserId": follower.UserId, "FollowerId": follower.FollowerId}); err != nil {
			nw.Error("The user already follows this user")
			log.Printf("[%s] [AddFollower] The user already follows this user", r.RemoteAddr)
			return
		}

		// We insert in the table Follower of the db the structure created
		if err := follower.InsertIntoDb(db); err != nil {
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
		var follower model.Follower
		if err := json.NewDecoder(r.Body).Decode(&follower); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [RemoveFollower] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

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

		if err = model.RemoveFromDB("Follower", db, map[string]any{"UserId": follower.UserId, "FollowerId": follower.FollowerId}); err != nil {
			nw.Error("Internal Error: There is a probleme during the delete in the DB: " + err.Error())
			log.Printf("[%s] [RemoveFollower] %s", r.RemoteAddr, err.Error())
			return
		}

		// We send a success response to the request
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Remove follower successfully",
		})
		if err != nil {
			log.Printf("[%s] [RemoveFollower] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetFollowed(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var follower model.Follower
		if err := json.NewDecoder(r.Body).Decode(&follower); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetFollower] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We decrypt the Id of the user make the request to follow someone
		decryptAuthorId, err := utils.DecryptJWT(follower.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetFollower] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		follower.UserId = decryptAuthorId

		if err = utils.IfExistsInDB("Auth", db, map[string]any{"Id": follower.UserId}); err != nil {
			nw.Error("Invalid Id in the JWT")
			log.Printf("[%s] [GetFollower] Invalid Id in the JWT : %v", r.RemoteAddr, err)
			return
		}

		var follows = model.Followers{
			{
				UserId: follower.UserId,
			},
		}

		// var follows model.Followers
		// follows[0].UserId = follower.UserId
		if err := follows.SelectFromDbByUserId(db); err != nil {
			nw.Error("Internal Error: There is a probleme during the selecte in the DB: " + err.Error())
			log.Printf("[%s] [GetFollower] %s", r.RemoteAddr, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get followed successfuly",
			"Follow":  follows,
		})
		if err != nil {
			log.Printf("[%s] [GetFollower] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetFollower(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var follower struct {
			UserId string `json:"UserId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&follower); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetFollower] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We decrypt the Id of the user make the request to follow someone
		decryptAuthorId, err := utils.DecryptJWT(follower.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetFollower] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		follower.UserId = decryptAuthorId

		if err = utils.IfExistsInDB("Auth", db, map[string]any{"Id": follower.UserId}); err != nil {
			nw.Error("Invalid Id in the JWT")
			log.Printf("[%s] [GetFollower] Invalid Id in the JWT : %v", r.RemoteAddr, err)
			return
		}

		var follows = model.Followers{
			{
				FollowerId: follower.UserId,
			},
		}

		// var follows model.Followers
		// follows[0].FollowerId = follower.UserId
		if err := follows.SelectFromDbByFollowerId(db); err != nil {
			nw.Error("Internal Error: There is a probleme during the selecte in the DB: " + err.Error())
			log.Printf("[%s] [GetFollower] %s", r.RemoteAddr, err.Error())
			return
		}


		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get followed successfuly",
			"Follow":  follows,
		})
		if err != nil {
			log.Printf("[%s] [GetFollower] %s", r.RemoteAddr, err.Error())
		}
	}
}

func IsFollowedBy(follower, followed string, db *sql.DB) bool {
	follow, err := model.SelectFromDb("Follower", db, map[string]any{"UserId": follower, "FollowedId": followed})
	if err != nil {
		return false
	}

	follows, err := follow.ParseFollowersData()
	if err != nil {
		return false
	}

	return len(follows) == 1
}
