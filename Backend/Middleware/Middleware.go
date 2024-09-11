package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	registermiddlewaresubfunction "social-network/Middleware/RegisterMiddlewareSubFunction"
	model "social-network/Model"
)

func RegisterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var register model.Register
		json.Unmarshal(body, &register)
		json.Unmarshal(body, &register.Auth)

		if err := registermiddlewaresubfunction.RegisterVerification(register, w); err != nil {
			return
		}

		if err := registermiddlewaresubfunction.CreateUuidAndCrypt(&register, w); err != nil {
			return
		}

		json, _ := json.Marshal(register)
		ctx := context.WithValue(r.Context(), model.RegisterCtx, json)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
