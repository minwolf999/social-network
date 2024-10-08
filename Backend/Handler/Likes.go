package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "social-network/Model"
)

// HandleLike manages the like system for posts
func HandleLike(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		var request struct {
			PostID string `json:"PostID"`
			UserID string `json:"UserID"`
			Table  string `json:"Table"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [Like] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		if request.Table != "LikePost" && request.Table != "DislikePost" && request.Table != "LikeComment" && request.Table != "DislikeComment" {
			nw.Error("Invalid table")
			log.Printf("[%s] [Like] Invalid table", r.RemoteAddr)
			return
		}

		if err := handleLikeLogic(db, request.Table, request.PostID, request.UserID); err != nil {
			nw.Error("Error during the like logic")
			log.Printf("[%s] [Like] Error during the like logic : %v", r.RemoteAddr, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"Success":   true,
			"Message":   "Like handled successfully",
		})
		if err != nil {
			log.Printf("[%s] [Like] %s", r.RemoteAddr, err.Error())
		}
	}
}

func handleLikeLogic(db *sql.DB, table, postID, userID string) error {
	// Check if the user has already liked the post
	liked, err := hasUserLikedPost(db, table, postID, userID)
	if err != nil {
		return fmt.Errorf("error checking like status: %v", err)
	}

	if liked {
		// If already liked, remove the like
		err = removeLike(db, table, postID, userID)
	} else {
		// If not liked, add the like
		err = addLike(db, table, postID, userID)
	}

	if err != nil {
		return fmt.Errorf("error handling like: %v", err)
	}

	return nil
}

func hasUserLikedPost(db *sql.DB, table, postID, userID string) (bool, error) {
	stmt, err := db.Prepare(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE PostID = ? AND UserID = ?", table))
	if err != nil {
		return false, err
	}
	
	var count int
	err = stmt.QueryRow(postID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func addLike(db *sql.DB, table, postID, userID string) error {
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (PostID, UserID) VALUES (?, ?)", table))
	if err != nil {
		return err
	}
	
	_, err = stmt.Exec(postID, userID)
	if err != nil {
		return err
	}

	parentTable := ""
	if table == "LikePost" || table == "DislikePost" {
		parentTable = "Post"
	} else if table == "likeComment" || table == "DislikeComment" {
		parentTable = "Comment"
	}

	return updateLikeCount(db, parentTable, postID, 1)
}

func removeLike(db *sql.DB, table, postID, userID string) error {
	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE PostID = ? AND UserID = ?", table))
	if err != nil {
		return err
	}

	result, err := stmt.Exec(postID, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("like not found")
	}

	parentTable := ""
	if table == "LikePost" || table == "DislikePost" {
		parentTable = "Post"
	} else if table == "likeComment" || table == "DislikeComment" {
		parentTable = "Comment"
	}

	return updateLikeCount(db, parentTable, postID, -1)
}

func updateLikeCount(db *sql.DB, table, postID string, delta int) error {
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE %s SET LikeCount = LikeCount + ? WHERE ID = ?", table))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(delta, postID)
	return err
}
