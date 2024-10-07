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

func CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var post model.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [CreatePost] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We decrypt the post author Id
		decryptAuthorId, err := utils.DecryptJWT(post.AuthorId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [CreatePost] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		post.AuthorId = decryptAuthorId

		if post.Text == "" || post.CreationDate == "" || (post.Status != "public" && post.Status != "private" && strings.Split(post.Status, " | ")[0] != "almost private") {
			nw.Error("There is an empty field")
			log.Printf("[%s] [CreatePost] there is an empty field", r.RemoteAddr)
			return
		}

		// We create a UID for the post
		uuid, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a probleme with the generation of the uuid")
			log.Printf("[%s] [CreatePost] There is a probleme with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}
		post.Id = uuid.String()

		// We insert the post in the db
		if err = post.InsertIntoDb(db); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [Createpost] %s", r.RemoteAddr, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Post created successfully",
			"IdPost":  post.Id,
		})
		if err != nil {
			log.Printf("[%s] [CreatePost] %s", r.RemoteAddr, err.Error())
		}
	}
}

func GetPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var post model.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetPost] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We decrypt the post author ID
		_, err := utils.DecryptJWT(post.AuthorId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetPost] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		// We check if there is a precise Post to get and make the request
		var posts model.Posts
		if post.Id != "" {
			err = post.SelectFromDbById(db)
			posts[0] = post
		} else {
			err = posts.SelectAllFromDb(db)
		}
		if err != nil {
			nw.Error("Error during the select in the db")
			log.Printf("[%s] [GetPost] Error during the select in the db : %v", r.RemoteAddr, err)
			return
		}

		for i, v := range posts {
			if (v.Status == "private" && !IsFollowedBy(post.AuthorId, v.AuthorId, db)) || (strings.Split(v.Status, " | ")[0] == "almost private" && !slices.Contains(strings.Split(v.Status, " | ")[1:], post.AuthorId)) {
				if i < len(posts)-1 {
					posts = append(posts[:i], posts[i+1:]...)
				} else {
					posts = posts[:i]
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get posts successfuly",
			"Posts":   posts,
		})
		if err != nil {
			log.Printf("[%s] [GetPost] %s", r.RemoteAddr, err.Error())
		}
	}
}
