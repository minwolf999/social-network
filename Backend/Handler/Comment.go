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

func CreateComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var comment model.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [CreateComent] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		if comment.Text == "" {
			nw.Error("There is no text for the comment")
			log.Printf("[%s] [CreateComment] There is no text for the comment", r.RemoteAddr)
			return
		}

		// We decrypt the comment author Id
		decryptAuthorId, err := utils.DecryptJWT(comment.AuthorId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [CreateComment] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		comment.AuthorId = decryptAuthorId

		var post = model.Post{
			Id: comment.PostId,
		}

		// We check if the id given for the parent post fit with a real post id in the db
		if err := post.SelectFromDbById(db); err != nil {
			nw.Error("Internal error: Problem during database query: " + err.Error())
			log.Printf("[%s] [CreateComment] %s", r.RemoteAddr, err.Error())
			return
		}

		if post.Text == "" {
			nw.Error("There is no post with the Id : " + comment.PostId)
			log.Printf("[%s] [CreateComment] There is no post with the Id : %s", r.RemoteAddr, comment.PostId)
			return
		}

		// We create a UID for the comment
		uuid, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a probleme with the generation of the uuid")
			log.Printf("[%s] [CreateComment] There is a probleme with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}
		comment.Id = uuid.String()

		// We insert the comment in the db
		if err = comment.InsertIntoDb(db); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [CreateComment] %s", r.RemoteAddr, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Comment created successfully",
		})
		if err != nil {
			log.Printf("[%s] [CreateComment] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var comment model.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetComment] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We decrypt the post author ID
		_, err := utils.DecryptJWT(comment.AuthorId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetComment] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		var comments model.Comments
		// We check if there is a precise Comment to get and make the request
		if comment.PostId != "" {
			err = comment.SelectFromDbById(db)
			comments[0] = comment
		} else {
			err = comments.SelectAllFromDb(db)
		}
		if err != nil {
			nw.Error("Error during the select in the db")
			log.Printf("[%s] [GetComment] Error during the select in the db : %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get comments successfuly",
			"Posts":   comments,
		})
		if err != nil {
			log.Printf("[%s] [GetComment] %s", r.RemoteAddr, err.Error())
		}
	}
}
