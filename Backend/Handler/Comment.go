package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
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
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var comment model.Comment
		json.Unmarshal(body, &comment)

		// We take the post id from the url used
		comment.PostId = r.PathValue("postId")

		if err := Verification(&comment, db); err != nil {
			nw.Error("Internal Error: There is a probleme during the verification: " + err.Error())
			log.Printf("[%s] [CreateComment] %s", r.RemoteAddr, err.Error())
			return
		}

		// We insert the comment in the db
		if err := utils.InsertIntoDb("Comment", db, comment.Id, comment.AuthorId, comment.Text, comment.CreationDate, comment.PostId); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [CreateComment] %s", r.RemoteAddr, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Comment created successfully",
		})
		if err != nil {
			log.Printf("[%s] [CreateComment] %s", r.RemoteAddr, err.Error())
		}
	}
}

func Verification(comment *model.Comment, db *sql.DB) error {
	if comment.Text == "" {
		return errors.New("there is no text for the comment")
	}

	// We decrypt the comment author Id
	decryptAuthorId, err := utils.DecryptJWT(comment.AuthorId, db)
	if err != nil {
		return err
	}
	comment.AuthorId = decryptAuthorId

	// We check if the id given for the parent post fit with a real post id in the db
	post, err := utils.SelectFromDb("Post", db, map[string]any{"Id": comment.PostId})
	if err != nil {
		return err
	}

	if len(post) != 1 {
		return errors.New("There is no post with the Id : " + comment.PostId)
	}

	// We create a UID for the comment
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	comment.Id = uuid.String()

	return nil
}

func GetComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var comment model.Comment
		json.Unmarshal(body, &comment)

		// We take the post id from the url used
		comment.PostId = r.PathValue("postId")

		// We decrypt the post author ID
		_, err := utils.DecryptJWT(comment.AuthorId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetPost] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		// We check if there is a precise Comment to get and make the request
		comments, err := SwitchCaseRequest("Comment", db, comment)
		if err != nil {
			nw.Error("Error during the select in the db")
			log.Printf("[%s] [GetComment] Error during the select in the db : %v", r.RemoteAddr, err)
			return
		}

		// We parse the result of the request in the good structure
		formatedComments, err := utils.ParseCommentData(comments)
		if err != nil {
			nw.Error(err.Error())
			log.Printf("[%s] [GetPost] %s", r.RemoteAddr, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get comments successfuly",
			"Posts":   formatedComments,
		})
		if err != nil {
			log.Printf("[%s] [GetComment] %s", r.RemoteAddr, err.Error())
		}
	}
}

func SwitchCaseRequest(table string, db *sql.DB, comment model.Comment) ([]map[string]any, error) {
	if comment.Id != "" {
		return utils.SelectFromDb(table, db, map[string]any{"Id": comment.Id})
	} else {
		return utils.SelectFromDb(table, db, map[string]any{})
	}
}
