package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"

	"github.com/gofrs/uuid"
)

// CreateComment handles the creation of a new comment.
// It expects the request body to contain the comment details in JSON format.
func CreateComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to handle error responses
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a Comment object to hold the decoded request body
		var comment model.Comment
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [CreateComment] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Check if the comment text is empty
		if comment.Text == "" {
			nw.Error("There is no text for the comment")
			log.Printf("[%s] [CreateComment] There is no text for the comment", r.RemoteAddr)
			return
		}

		// Decrypt the AuthorId from the JWT to get the actual author ID
		decryptAuthorId, err := utils.DecryptJWT(comment.AuthorId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [CreateComment] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		comment.AuthorId = decryptAuthorId // Set the decrypted author ID

		// Create a Post object to check if the post exists
		var post = model.Post{
			Id: comment.PostId,
		}

		// Check if the post exists in the database
		if err := post.SelectFromDb(db, map[string]any{"Id": post.Id}); err != nil {
			nw.Error("Internal error: Problem during database query: " + err.Error())
			log.Printf("[%s] [CreateComment] %s", r.RemoteAddr, err.Error())
			return
		}

		// Check if the post has valid content
		if post.Text == "" {
			nw.Error("There is no post with the Id : " + comment.PostId)
			log.Printf("[%s] [CreateComment] There is no post with the Id : %s", r.RemoteAddr, comment.PostId)
			return
		}

		// Generate a new UUID for the comment
		uid, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
			log.Printf("[%s] [CreateComment] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}
		comment.Id = uid.String() // Set the generated UUID as the comment ID

		// Insert the comment into the database
		if err = comment.InsertIntoDb(db); err != nil {
			nw.Error("Internal Error: There is a problem during the push in the DB: " + err.Error())
			log.Printf("[%s] [CreateComment] %s", r.RemoteAddr, err.Error())
			return
		}

		notifId, err := uuid.NewV7()
		if err != nil {
			nw.Error("There is a problem with the generation of the uuid") // Handle UUID generation error
			log.Printf("[%s] [CreateComment] There is a problem with the generation of the uuid : %s", r.RemoteAddr, err)
			return
		}

		var userData model.Register
		if err = userData.SelectFromDb(db, map[string]any{"Id": comment.AuthorId}); err != nil {
			nw.Error("There is a problem during the fetching of the user") // Handle UUID generation error
			log.Printf("[%s] [CreateComment] There is a problem during the fetching of the user : %s", r.RemoteAddr, err)
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
			UserId:      post.AuthorId,
			Status:      "Comment",
			Description: fmt.Sprintf("A comment as been posted by %s for your post \"%s\"", userDataName, post.Text),
			GroupId:     post.IsGroup,
			OtherUserId: "",
		}

		if err = notification.InsertIntoDb(db); err != nil {
			nw.Error("There is a probleme during the sending of a notification")
			log.Printf("[%s] [CreateComment] There is a probleme during the sending of a notification : %s", r.RemoteAddr, err)
			return
		}

		model.ConnectedWebSocket.Mu.Lock()
		for i := range model.ConnectedWebSocket.Conn {
			_, isOk := model.ConnectedWebSocket.Conn[i]
			if isOk {
				var WebsocketMessage struct {
					Type        string
					AuthorId    string
					PostId      string
					Description string
					Value       model.Comment
				}

				var commentdetail model.Comment
				if err = commentdetail.SelectFromDb(db, map[string]any{"Id": comment.Id}); err != nil {
					nw.Error("Error during the fetch of the db")
					log.Printf("[%s] [CreateComment] Error during the fetch of the db : %s", r.RemoteAddr, err)
					return
				}

				WebsocketMessage.Type = "Comment"
				WebsocketMessage.AuthorId = post.AuthorId
				WebsocketMessage.PostId = post.Id
				WebsocketMessage.Description = fmt.Sprintf("A comment as been posted by %s for your post \"%s\"", userDataName, post.Text)
				WebsocketMessage.Value = commentdetail


				if err = model.ConnectedWebSocket.Conn[i].WriteJSON(WebsocketMessage); err != nil {

					nw.Error("Error during the communication with the websocket")
					log.Printf("[%s] [CreateComment] Error during the communication with the websocket : %s", r.RemoteAddr, err)
					return
				}
			}
		}
		model.ConnectedWebSocket.Mu.Unlock()

		// Send a success response in JSON format
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

// GetComment handles the retrieval of comments.
// It expects the request body to contain the comment ID and post ID.
func GetComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom ResponseWriter to handle error responses
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Initialize a Comment object to hold the decoded request body
		var comment model.Comment
		// Decode the JSON request body into the comment object
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			nw.Error("Invalid request body") // Send error if decoding fails
			log.Printf("[%s] [GetComment] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// Decrypt the AuthorId from the JWT to ensure the request is valid
		_, err := utils.DecryptJWT(comment.AuthorId, db)
		if err != nil {
			nw.Error("Invalid JWT") // Handle invalid JWT error
			log.Printf("[%s] [GetComment] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		var comments model.Comments // Initialize a Comments object to hold the retrieved comments

		// Check if a specific post ID is provided
		if comment.PostId != "" {
			err = comment.SelectFromDb(db, map[string]any{"Id": comment.Id}) // Retrieve a specific comment
			comments = append(comments, comment)                             // Add the comment to the comments slice
		} else {
			// If no post ID is provided, retrieve all comments
			err = comments.SelectFromDb(db, map[string]any{})
		}
		if err != nil {
			nw.Error("Error during the select in the db") // Handle database selection error
			log.Printf("[%s] [GetComment] Error during the select in the db : %v", r.RemoteAddr, err)
			return
		}

		// Send a success response with the retrieved comments
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"Success": true,
			"Message": "Get comments successfully",
			"Posts":   comments, // Include the comments in the response
		})
		if err != nil {
			log.Printf("[%s] [GetComment] %s", r.RemoteAddr, err.Error())
		}
	}
}
