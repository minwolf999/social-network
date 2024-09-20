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

		if comment.Text == "" {
			nw.Error("There is no text for the comment")
			log.Printf("[%s] [CreateComment] There is no text for the comment", r.RemoteAddr)
			return
		}

		// We decrypt the comment author Id
		decryptAuthorId, err := utils.DecryptJWT(comment.AuthorId)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [CreateComment] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}
		comment.AuthorId = decryptAuthorId

		// We check if the id given for the parent post fit with a real post id in the db
		post, err := utils.SelectFromDb("Post", db, map[string]any{"Id": comment.PostId})
		if err != nil {
			nw.Error("Internal error: Problem during database query: " + err.Error())
			log.Printf("[%s] [CreateComment] %s", r.RemoteAddr, err.Error())
			return
		}

		if len(post) != 1 {
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
		if err = utils.InsertIntoDb("Comment", db, comment.Id, comment.AuthorId, comment.Text, comment.CreationDate, comment.PostId); err != nil {
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
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var comment model.Comment
		json.Unmarshal(body, &comment)

		// We take the post id from the url used
		comment.PostId = r.PathValue("postId")

		// We decrypt the post author ID
		_, err := utils.DecryptJWT(comment.AuthorId)
		if err != nil {
			nw.Error("Invalid JWT")
			log.Printf("[%s] [GetPost] Error during the decrypt of the JWT : %v", r.RemoteAddr, err)
			return
		}

		var comments []map[string]any
		// We check if there is a precise Comment to get and make the request
		if comment.Id != "" {
			comments, err = utils.SelectFromDb("Comment", db, map[string]any{"Id": comment.Id})
		} else {
			comments, err = utils.SelectFromDb("Comment", db, map[string]any{})
		}
		if err != nil {
			nw.Error("Error during the select in the db")
			log.Printf("[%s] [GetComment] Error during the select in the db : %v", r.RemoteAddr, err)
			return
		}

		// We parse the result of the request in the good structure
		formatedComments, err := ParseCommentData(comments)
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

func ParseCommentData(userData []map[string]any) ([]model.Comment, error) {
	// We marshal the map to get it in []byte
	serializedData, err := json.Marshal(userData)
	if err != nil {
		return nil, errors.New("internal error: conversion problem")
	}

	// We Unmarshal in the good structure
	var postResult []model.Comment
	err = json.Unmarshal(serializedData, &postResult)
	return postResult, err
}
