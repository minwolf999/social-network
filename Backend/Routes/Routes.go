package routes

import (
	"fmt"
	"net/http"
	handler "social-network/Handler"
	middleware "social-network/Middleware"
	utils "social-network/Utils"
)

func Routes(mux *http.ServeMux) {

	db, err := utils.OpenDb("sqlite3", "./Database/Database.sqlite")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	mux.HandleFunc("/", handler.Redirect)

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) { handler.Login(w, r) })
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) { middleware.RegisterMiddleware(handler.Register, db) })

	mux.HandleFunc("/home", handler.Home)
}

// Mock du handler Login pour les tests
func mockLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Simuler une réponse OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// Mock du handler Register pour les tests
func mockRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Simuler une réponse OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Register successful"))
}

// Fonction Routes pour les tests avec des mocks
func RoutesForTest(mux *http.ServeMux) {
	// On remplace les handlers réels par des mocks
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Redirected"))
	})

	mux.HandleFunc("/login", mockLoginHandler)
	mux.HandleFunc("/register", mockRegisterHandler)
	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome Home"))
	})
}
