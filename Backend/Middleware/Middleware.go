package middleware

import (
	"database/sql"
	"net/http"
	
	model "social-network/Model"
)


/*
This function takes 1 argument:
  - an http.HandlerFunc who is handleFunction who gonna be called after this function

The purpose of this function is to set the good Cors method.

The function return an http.HandlerFunc (it's a function)
*/
func SetHeaderAccessControll(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		next(w, r)
	}
}


/*
This function takes 1 argument:
  - an http.HandlerFunc who is handleFunction who gonna be called after this function

The purpose of this function is to check if the request method is a good request method.

The function return an http.HandlerFunc (it's a function)
*/
func LookMethod(next func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		if r.Method != http.MethodPost {
			nw.Error("Invalid method !")
			return
		}

		next(w, r, db)
	}
}
