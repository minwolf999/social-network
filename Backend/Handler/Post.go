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

func CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var post model.Post
		json.Unmarshal(body, &post)

		// We decrypt the post author Id
		decryptAuthorId, err := utils.DecryptJWT(post.AuthorId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [CreatePost] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		post.AuthorId = decryptAuthorId

		if post.Text == "" {
			nw.Error("There is no text for the post")
			log.Printf("[%s] [CreatePost] There is no text for the post", r.RemoteAddr)
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
		if err = utils.InsertIntoDb("Post", db, post.Id, post.AuthorId, post.Text, post.Image, post.IsGroup); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [Createpost] %s", r.RemoteAddr, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Post created successfully",
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
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var post model.Post
		json.Unmarshal(body, &post)

		// We decrypt the post author ID
		_, err := utils.DecryptJWT(post.AuthorId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetPost] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		// We check if there is a precise Post to get and make the request
		var posts []map[string]any
		if post.Id != "" {
			posts, err = utils.SelectFromDb("Post", db, map[string]any{"Id": post.Id})
		} else {
			posts, err = utils.SelectFromDb("Post", db, map[string]any{})
		}
		if err != nil {
			nw.Error("Error during the select in the db")
			log.Printf("[%s] [GetPost] Error during the select in the db : %v", r.RemoteAddr, err)
			return
		}

		// We parse the result of the request in the good structure
		formatedPosts, err := utils.ParsePostData(posts)
		if err != nil {
			nw.Error(err.Error())
			log.Printf("[%s] [GetPost] %s", r.RemoteAddr, err.Error())
			return
		}

		//---------------------------------------------------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------------------------------------
		//							Filtrer les posts pour n'obtenir que ceux des amis
		//---------------------------------------------------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------------------------------------
		//---------------------------------------------------------------------------------------------------------------------------------------------------

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get posts successfuly",
			"Posts":   formatedPosts,
		})
		if err != nil {
			log.Printf("[%s] [GetPost] %s", r.RemoteAddr, err.Error())
		}
	}
}
