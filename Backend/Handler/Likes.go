package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// HandleLike manages the like system for posts
func HandleLike(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			PostID string `json:"PostID"`
			UserID string `json:"UserID"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := handleLikeLogic(db, request.PostID, request.UserID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error handling like: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Like handled successfully"})
	}
}

func handleLikeLogic(db *sql.DB, postID, userID string) error {
	// Check if the user has already liked the post
	liked, err := hasUserLikedPost(db, postID, userID)
	if err != nil {
		return fmt.Errorf("error checking like status: %v", err)
	}

	if liked {
		// If already liked, remove the like
		err = removeLike(db, postID, userID)
	} else {
		// If not liked, add the like
		err = addLike(db, postID, userID)
	}

	if err != nil {
		return fmt.Errorf("error handling like: %v", err)
	}

	return nil
}

func hasUserLikedPost(db *sql.DB, postID, userID string) (bool, error) {
	query := "SELECT COUNT(*) FROM LikePost WHERE PostID = ? AND UserID = ?"
	var count int
	err := db.QueryRow(query, postID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func addLike(db *sql.DB, postID, userID string) error {
	_, err := db.Exec("INSERT INTO LikePost (PostID, UserID) VALUES (?, ?)", postID, userID)
	if err != nil {
		return err
	}
	return updateLikeCount(db, postID, 1)
}

func removeLike(db *sql.DB, postID, userID string) error {
	result, err := db.Exec("DELETE FROM LikePost WHERE PostID = ? AND UserID = ?", postID, userID)
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
	return updateLikeCount(db, postID, -1)
}

func updateLikeCount(db *sql.DB, postID string, delta int) error {
	_, err := db.Exec("UPDATE Post SET LikeCount = LikeCount + ? WHERE ID = ?", delta, postID)
	return err
}
