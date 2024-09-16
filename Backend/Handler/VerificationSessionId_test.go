package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	model "social-network/Model"
	utils "social-network/Utils"

	"golang.org/x/crypto/bcrypt"
)

func TestVerificationSessionId(t *testing.T) {
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
		Id:       "test",
		Email:    "unemail@gmail.com",
		Password: "MonMotDePasse123!",
	}

	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), 12)
	if err != nil {
		t.Fatalf("Erreur lors du cryptage du mot de passe : %v", err)
	}

	if err = utils.InsertIntoDb("Auth", db, login.Id, login.Email, string(cryptedPassword)); err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
	}

	sessionId := base64.StdEncoding.EncodeToString([]byte(login.Id))

	body, err := json.Marshal(sessionId)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/VerificationSessionId", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(VerificationSessionId(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "Valid cookie"
	bodyValue := make(map[string]any)

	if err = json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
	}

	if bodyValue["Message"] != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}


func TestCheckDatasForCookie(t *testing.T) {
	login := []map[string]any{
		{
			"Id":       "test",
			"Email":    "unemail@gmail.com",
			"Password": "MonMotDePasse123!",
		},
	}

	if err := CheckDatasForCookie(login); err != nil {
		t.Fatalf("Erreur lors de l'execution de la fonction: %s", err)
	}
}