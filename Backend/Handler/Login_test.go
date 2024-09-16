package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	model "social-network/Model"
	utils "social-network/Utils"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := utils.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
	}
	defer db.Close()

	// Create a table for testing
	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS Auth (
			Id VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY,
			Email VARCHAR(100) NOT NULL UNIQUE,
			Password VARCHAR(50) NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("Erreur lors de la création de la table : %v", err)
	}
	
	// Crée une structure Register de test
	login := model.Auth{
		Email:    "unemail@gmail.com",
		Password: "MonMotDePasse123!",
	}

	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), 12)
	if err != nil {
		t.Fatalf("Erreur lors du cryptage du mot de passe : %v", err)
	}

	if err = utils.InsertIntoDb("Auth", db, "0", login.Email, string(cryptedPassword)); err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
	}

	body, err := json.Marshal(login)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	next := func(next func(http.ResponseWriter, *http.Request, *sql.DB), db *sql.DB) http.HandlerFunc {
		// Si le middleware est exécuté correctement, cela signifie qu'on arrive à cette fonction
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r, db)
		}
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(next(Login, db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "Login successfully"
	bodyValue := make(map[string]any)

	if err = json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
	}

	if bodyValue["Message"] != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
