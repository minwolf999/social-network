package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	model "social-network/Model"
)

/*
This function takes 1 argument:
  - an http.HandlerFunc who is handleFunction who gonna be called after this function

The purpose of this function is to Verificate the content of the request make to the Login function.

The function return an http.HandlerFunc (it's a function)
*/
func LoginMiddleware(next func(w http.ResponseWriter, r *http.Request, db *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// We read the request body and unmarshal it into a structure
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var login model.Auth
		json.Unmarshal(body, &login)

		// We look if all is good in the datas send in the body of the request
		if login.Email == "" || login.Password == "" {
			nw.Error("There is an empty field")
			return
		}

		// We marshal the datas and set it in the context
		json, _ := json.Marshal(login)
		ctx := context.WithValue(r.Context(), model.LoginCtx, json)
		r = r.WithContext(ctx)

		next(w, r, db)
	}
}
