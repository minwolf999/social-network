package utils

import (
	"encoding/base64"
	"net/http"

	model "social-network/Model"
)

/*
This function takes 2 arguments:
  - an http.ResponseWriter
  - an Register variable who contain the datas send in the body of the request

The purpose of this function is to create a session cookie.

The function return no value
*/
func SetCookie(w http.ResponseWriter, register model.Register) {
	cookieEmail := http.Cookie{
		Name:     "sessionId",
		Value:    base64.StdEncoding.EncodeToString([]byte(register.Auth.Id)),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(w, &cookieEmail)
}
