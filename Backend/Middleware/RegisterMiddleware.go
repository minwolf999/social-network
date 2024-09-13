package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	registermiddlewaresubfunction "social-network/Middleware/RegisterMiddlewareSubFunction"
	model "social-network/Model"
)

/*
This function takes 1 argument:
  - an http.HandlerFunc who is handleFunction who gonna be called after this function

The purpose of this function is to Verificate the content of the request make to the Register function.

The function return an http.HandlerFunc (it's a function)
*/
func RegisterMiddleware(next func(w http.ResponseWriter, r *http.Request, db *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		fmt.Println("new request")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Accept-Encoding", "gzip, deflate, br")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Accept", "/*")

		// We read the request body and unmarshal it into a structure
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var register model.Register
		json.Unmarshal(body, &register)
		json.Unmarshal(body, &register.Auth)

		// We look if all is good in the datas send in the body of the request
		if err := registermiddlewaresubfunction.RegisterVerification(register); err != nil {
			nw.Error(err.Error())
			return
		}

		// We generate an UUID and crypt the password
		if err := registermiddlewaresubfunction.CreateUuidAndCrypt(&register); err != nil {
			nw.Error(err.Error())
			return
		}

		// We marshal the datas and set it in the context
		json, _ := json.Marshal(register)
		ctx := context.WithValue(r.Context(), model.RegisterCtx, json)
		r = r.WithContext(ctx)

		next(w, r, db)
	}
}
