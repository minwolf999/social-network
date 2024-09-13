package utils

import (
	"net/http"
	"time"
)

/*
This function takes 2 arguments:
  - an http.ResponseWriter
  - an Register variable who contain the datas send in the body of the request

The purpose of this function is to create a session cookie.

The function return no value
*/
func SetCookie(w http.ResponseWriter, value string) {
	cookieEmail := http.Cookie{
		Name:     "sessionId",
		Value:    value,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(w, &cookieEmail)
}
