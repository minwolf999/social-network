package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	model "social-network/Model"
	utils "social-network/Utils"

	"github.com/gofrs/uuid"
)

/*
This function handles the creation of a new post.

It takes a pointer to an SQL database as an argument and returns an http.HandlerFunc.

The function verifies the request, processes the data, and interacts with the database to create a new post.
*/
func CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize a new ResponseWriter for structured error handling.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Create a variable to hold the post data.
		var post model.Post

		// Decode the incoming JSON request body into the post structure.
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [CreatePost] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the Author ID from the JWT.
		decryptAuthorId, err := utils.DecryptJWT(post.AuthorId, db)
		if err != nil {
			// Return error if JWT decryption fails.
			nw.Error("Invalid JWT")
			log.Printf("[%s] [CreatePost] Error during the decrypt of the JWT: %v", r.RemoteAddr, err)
			return
		}
		// Update the Author ID with the decrypted value.
		post.AuthorId = decryptAuthorId

		// Validate the post fields: text, creation date, and status.
		if post.Text == "" || post.CreationDate == "" ||
			(post.Status != "public" && post.Status != "private" && strings.Split(post.Status, " | ")[0] != "almost private") {
			// Return error if any required fields are empty or invalid.
			nw.Error("There is an empty field")
			log.Printf("[%s] [CreatePost] There is an empty field", r.RemoteAddr)
			return
		}

		// Generate a new UUID for the post.
		uid, err := uuid.NewV7()
		if err != nil {
			// Return error if UUID generation fails.
			nw.Error("There is a problem with the generation of the uuid")
			log.Printf("[%s] [CreatePost] There is a problem with the generation of the uuid: %s", r.RemoteAddr, err)
			return
		}
		// Set the post ID to the newly generated UUID.
		post.Id = uid.String()

		// Attempt to insert the post into the database.
		if err = post.InsertIntoDb(db); err != nil {
			// Return error if there is a problem during the database insertion.
			nw.Error("Internal Error: There is a problem during the push in the DB: " + err.Error())
			log.Printf("[%s] [CreatePost] %s", r.RemoteAddr, err.Error())
			return
		}

		if post.IsGroup != "" {
			var group model.Group
			if err = group.SelectFromDb(db, map[string]any{"Id": post.IsGroup}); err != nil {
				nw.Error("There is a problem during the fetching of the group") // Handle UUID generation error
				log.Printf("[%s] [CreatePost] There is a problem during the fetching of the group : %s", r.RemoteAddr, err)
				return
			}

			var userData model.Register
			if err = userData.SelectFromDb(db, map[string]any{"Id": post.AuthorId}); err != nil {
				nw.Error("There is a problem during the fetching of the user") // Handle UUID generation error
				log.Printf("[%s] [CreatePost] There is a problem during the fetching of the user : %s", r.RemoteAddr, err)
				return
			}

			var userDataName string
			if userData.Username == "" {
				userDataName = userData.FirstName + " " + userData.LastName
			} else {
				userDataName = userData.Username
			}

			group.SplitMembers()
			for i := range group.SplitMemberIds {

				notifId, err := uuid.NewV7()
				if err != nil {
					nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
					log.Printf("[%s] [CreatePost] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
					return
				}

				notification := model.Notification{
					Id:          notifId.String(),
					UserId:      group.SplitMemberIds[i],
					Status:      "Group",
					Description: fmt.Sprintf("A new post as been send by %s for the group %s", userDataName, group.GroupName),
					GroupId:     group.Id,
					OtherUserId: "",
				}

				if err = notification.InsertIntoDb(db); err != nil {
					nw.Error("There is a probleme during the sending of a notification")
					log.Printf("[%s] [CreatePost] There is a probleme during the sending of a notification : %s", r.RemoteAddr, err)
					return
				}

				model.ConnectedWebSocket.Mu.Lock()
				_, isOk := model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]]
				if isOk {
					var WebsocketMessage struct {
						Type        string
						GroupId     string
						Description string
						Value       model.Post
					}

					WebsocketMessage.Type = "GroupPost"
					WebsocketMessage.GroupId = group.Id
					WebsocketMessage.Description = "A post has been send to the group"
					WebsocketMessage.Value = post

					if err = model.ConnectedWebSocket.Conn[group.SplitMemberIds[i]].WriteJSON(WebsocketMessage); err != nil {

						nw.Error("Error during the communication with the websocket")
						log.Printf("[%s] [CreatePost] Error during the communication with the websocket : %s", r.RemoteAddr, err)
						return
					}
				}
				model.ConnectedWebSocket.Mu.Unlock()
			}
		}

		// Set response headers for JSON content.
		w.Header().Set("Content-Type", "application/json")

		// Encode the response JSON for successful post creation.
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Post created successfully",
			// Return the newly created post ID.
			"IdPost": post.Id,
		})
		if err != nil {
			log.Printf("[%s] [CreatePost] %s", r.RemoteAddr, err.Error())
		}
	}
}

/*
GetPost handles the retrieval of posts from the database.

It checks for a specific post by its ID or retrieves all posts if no ID is provided.
It also handles JWT decryption for the author's ID.
*/
func GetPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize a new ResponseWriter for structured error handling.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Create a variable to hold the post data from the request.
		var post model.Post

		// Decode the incoming JSON request body into the post structure.
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetPost] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the Author ID from the JWT to ensure it's valid.
		JWT, err := utils.DecryptJWT(post.AuthorId, db)
		if err != nil {
			// Return error if JWT decryption fails.
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetPost] Error during the decrypt of the JWT: %v", r.RemoteAddr, err)
			return
		}

		// Create a variable to hold the retrieved posts.
		var posts model.Posts

		// If a specific post ID is provided, retrieve that post; otherwise, retrieve all posts.
		if post.Id != "" {
			// Retrieve post by ID
			err = post.SelectFromDb(db, map[string]any{"Id": post.Id})
			// Add the retrieved post to the posts slice.
			posts[0] = post
		} else {
			// Retrieve all posts.
			err = posts.SelectFromDb(db, map[string]any{})
		}
		if err != nil {
			// Return error if there is a problem during the database query.
			nw.Error("Error during the select in the db")
			log.Printf("[%s] [GetPost] Error during the select in the db: %v", r.RemoteAddr, err)
			return
		}

		for i := 0; i < len(posts); i++ {
			if (JWT == posts[i].Status && posts[i].Status == "private" && !IsFollowedBy(JWT, posts[i].AuthorId, db)) || JWT == posts[i].Status && (strings.Split(posts[i].Status, " | ")[0] == "almost private" && !slices.Contains(strings.Split(posts[i].Status, " | ")[1:], post.AuthorId)) || JWT == posts[i].Status && posts[i].IsGroup != "" {
				if i < len(posts)-1 {
					// Remove the post from the middle of the slice.
					posts = append(posts[:i], posts[i+1:]...)
				} else {
					// Remove the last post.
					posts = posts[:i]
				}

				i--
			}
		}

		// Set response headers for JSON content.
		w.Header().Set("Content-Type", "application/json")

		// Encode the response JSON for successfully retrieved posts.
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			// Note the correction from "successfuly" to "successfully"
			"Message": "Get posts successfully",
			// Return the list of retrieved posts.
			"Posts": posts,
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [GetPost] %s", r.RemoteAddr, err.Error())
		}
	}
}
