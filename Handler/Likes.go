package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"
)

/*
This function takes 1 argument:
  - a pointer to an SQL database object

The purpose of this function is to handle user requests for liking or disliking posts or comments.

The function returns an http.HandlerFunc that can be used as a handler for HTTP requests.
*/
func HandleLike(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to handle errors and responses.
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Struct to hold the data from the request body.
		var request struct {
			// ID of the post or comment to be liked/disliked.
			PostID string `json:"PostID"`
			// ID of the user performing the like/dislike action.
			UserID string `json:"UserID"`
			// Table name indicating whether it's for likes or dislikes.
			Table string `json:"Table"`
		}

		// Decode the JSON request body into the request struct.
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [Like] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		userId, err := utils.DecryptJWT(request.UserID, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [Like] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		request.UserID = userId

		// Validate the provided table name.
		if request.Table != "LikePost" && request.Table != "DislikePost" && request.Table != "LikeComment" && request.Table != "DislikeComment" {
			// Return error if the table is not valid.
			nw.Error("Invalid table")
			log.Printf("[%s] [Like] Invalid table", r.RemoteAddr)
			return
		}

		// Process the like or dislike logic.
		if err := handleLikeLogic(db, request.Table, request.PostID, request.UserID); err != nil {
			// Return error if there is an issue during the like logic processing.
			nw.Error("Error during the like logic")
			log.Printf("[%s] [Like] Error during the like logic : %v", r.RemoteAddr, err)
			return
		}

		// Set the response header to indicate JSON content and respond with success message.
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			// Indicates the operation was successful.
			"Success": true,
			// Success message for the like action.
			"Message": "Like handled successfully",
		})
		if err != nil {
			// Log any error that occurs while encoding the response.
			log.Printf("[%s] [Like] %s", r.RemoteAddr, err.Error())
		}
	}
}

/*
This function takes 4 arguments:
  - a pointer to an SQL database object
  - a string representing the table name for likes/dislikes
  - a string containing the post ID
  - a string containing the user ID

The purpose of this function is to manage the like or dislike logic for a given post or comment.

The function returns an error if there is a problem during the like handling process; otherwise, it returns nil.
*/
func handleLikeLogic(db *sql.DB, table, postID, userID string) error {
	// Check if the user has already liked the post or comment.
	liked, err := hasUserLikedPost(db, table, postID, userID)
	if err != nil {
		// Return error if there is an issue checking the like status.
		return fmt.Errorf("error checking like status: %v", err)
	}

	if liked {
		// If the user has already liked, remove the like.
		err = removeLike(db, table, postID, userID)
	} else {
		// If the user has not liked, add a new like.
		err = addLike(db, table, postID, userID)
	}

	if err != nil {
		// Return error if there is a problem handling the like operation.
		return fmt.Errorf("error handling like: %v", err)
	}

	// Return nil if the operation was successful.
	return nil
}

/*
This function takes 4 arguments:
  - a pointer to an SQL database object
  - a string representing the table name for likes/dislikes
  - a string containing the post ID
  - a string containing the user ID

The purpose of this function is to check if a specific user has liked a given post.

The function returns a boolean indicating whether the user has liked the post and an error if there is an issue during the query.
*/
func hasUserLikedPost(db *sql.DB, table, postID, userID string) (bool, error) {
	// Prepare a SQL statement to count how many times the user has liked the specified post.
	stmt, err := db.Prepare(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE PostID = ? AND UserID = ?", table))
	if err != nil {
		// Return false and the error if statement preparation fails.
		return false, err
	}

	// Declare a variable to hold the count result.
	var count int
	// Execute the query and scan the result into the count variable.
	err = stmt.QueryRow(postID, userID).Scan(&count)
	if err != nil {
		// Return false and the error if querying fails.
		return false, err
	}
	// Return true if count is not zero, indicating the user has liked the post.
	return count != 0, nil
}

/*
This function takes 4 arguments:
  - a pointer to an SQL database object
  - a string representing the table name for likes/dislikes
  - a string containing the post ID
  - a string containing the user ID

The purpose of this function is to add a like for a specific post by a user.

The function returns an error if there is an issue during the insertion or updating process.
*/
func addLike(db *sql.DB, table, postID, userID string) error {
	// Prepare a SQL statement to insert a new like record into the specified table.
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (PostID, UserID) VALUES (?, ?)", table))
	if err != nil {
		// Return the error if statement preparation fails.
		return err
	}

	// Execute the insert statement with the post ID and user ID.
	_, err = stmt.Exec(postID, userID)
	if err != nil {
		// Return the error if execution fails.
		return err
	}

	// Initialize a variable to hold the parent table name.
	parentTable := ""
	// Determine the parent table based on the type of like/dislike action.
	if table == "LikePost" || table == "DislikePost" {
		parentTable = "Post"
	} else if table == "likeComment" || table == "DislikeComment" {
		parentTable = "Comment"
	}

	// Update the like count for the corresponding parent table and post ID.
	return updateLikeCount(db, parentTable, postID, 1)
}

/*
This function takes 4 arguments:
  - a pointer to an SQL database object
  - a string representing the table name for likes/dislikes
  - a string containing the post ID
  - a string containing the user ID

The purpose of this function is to remove a like for a specific post by a user.

The function returns an error if there is an issue during the deletion or updating process.
*/
func removeLike(db *sql.DB, table, postID, userID string) error {
	// Prepare a SQL statement to delete a like record from the specified table.
	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE PostID = ? AND UserID = ?", table))
	if err != nil {
		// Return the error if statement preparation fails.
		return err
	}

	// Execute the delete statement with the post ID and user ID.
	result, err := stmt.Exec(postID, userID)
	if err != nil {
		// Return the error if execution fails.
		return err
	}

	// Get the number of rows affected by the delete operation.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return the error if getting affected rows fails.
		return err
	}

	// Return an error if no rows were affected, indicating the like was not found.
	if rowsAffected == 0 {
		return fmt.Errorf("like not found")
	}

	// Initialize a variable to hold the parent table name.
	parentTable := ""

	// Determine the parent table based on the type of like/dislike action.
	if table == "LikePost" || table == "DislikePost" {
		parentTable = "Post"
	} else if table == "likeComment" || table == "DislikeComment" {
		parentTable = "Comment"
	}

	// Update the like count for the corresponding parent table and post ID.
	return updateLikeCount(db, parentTable, postID, -1)
}

/*
This function takes 4 arguments:
  - a pointer to an SQL database object
  - a string representing the table name where the like count is stored
  - a string containing the post ID for which the like count needs to be updated
  - an integer (delta) that indicates how much to change the like count (positive for adding likes, negative for removing likes)

The purpose of this function is to update the like count for a specific post in the database.

The function returns an error if there is an issue during the update process.
*/
func updateLikeCount(db *sql.DB, table, postID string, delta int) error {
	// Prepare a SQL statement to update the LikeCount in the specified table.
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE %s SET LikeCount = LikeCount + ? WHERE ID = ?", table))
	if err != nil {
		// Return the error if statement preparation fails.
		return err
	}

	// Execute the update statement with the delta value and post ID.
	_, err = stmt.Exec(delta, postID)

	// Return any error that occurs during execution.
	return err
}

func GetLike(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// Struct to hold the data from the request body.
		var request struct {
			UserId string `json:"UserId"`
			PostId string `json:"PostId"`
			Table  string `json:"Table"`
		}

		// Decode the JSON request body into the request struct.
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			// Return error if the request body is invalid.
			nw.Error("Invalid request body")
			log.Printf("[%s] [GetLike] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		userId, err := utils.DecryptJWT(request.UserId, db)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetLike] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		request.UserId = userId

		if request.Table != "LikePost" && request.Table != "LikeComment" {
			nw.Error("Invalid table")
			log.Printf("[%s] [GetLike] Invalid table", r.RemoteAddr)
			return
		}

		userData, err := model.SelectFromDb(request.Table, db, map[string]any{"GroupId": request.PostId})
		if err != nil {
			nw.Error("Error during the fetch of the DB")
			log.Printf("[%s] [GetLike] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		serializedData, err := json.Marshal(userData)
		if err != nil {
			nw.Error("Error during the fetch of the DB")
			log.Printf("[%s] [GetLike] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		if request.Table == "LikePost" {
			var PostLike []struct {
				UserId string `json:"UserId"`
				PostId string `json:"PostId"`
			}

			if err = json.Unmarshal(serializedData, &PostLike); err != nil {
				nw.Error("Error during the fetch of the DB")
				log.Printf("[%s] [GetLike] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
				return
			}

			// Set the response header to indicate JSON content and respond with success message.
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(map[string]any{
				// Indicates the operation was successful.
				"Success": true,
				// Success message for the like action.
				"Message": "Like getted successfully",

				"Value": PostLike,
			})
			if err != nil {
				// Log any error that occurs while encoding the response.
				log.Printf("[%s] [Like] %s", r.RemoteAddr, err.Error())
			}
		} else if request.Table == "LikeComment" {
			var CommentLike []struct {
				UserId    string `json:"UserId"`
				CommentId string `json:"CommentId"`
			}

			if err = json.Unmarshal(serializedData, &CommentLike); err != nil {
				nw.Error("Error during the fetch of the DB")
				log.Printf("[%s] [GetLike] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
				return
			}

			// Set the response header to indicate JSON content and respond with success message.
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(map[string]any{
				// Indicates the operation was successful.
				"Success": true,
				// Success message for the like action.
				"Message": "Like getted successfully",

				"Value": CommentLike,
			})
			if err != nil {
				// Log any error that occurs while encoding the response.
				log.Printf("[%s] [Like] %s", r.RemoteAddr, err.Error())
			}
		}
	}
}
