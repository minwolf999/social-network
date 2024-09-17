package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	model "social-network/Model"
	utils "social-network/Utils"
	"strings"
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

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Login(db))

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

	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGenerateJWT(t *testing.T) {
	value := "Test"
	JWT := GenerateJWT(value)

	splitJWT := strings.Split(JWT, ".")
	if len(splitJWT) != 3 {
		t.Errorf("The 3 part of the JWT are not here")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(splitJWT[2]), []byte(model.SecretKey)); err != nil {
		t.Errorf("Invalid secret key: %v", err)
		return
	}

	decrypt, err := base64.StdEncoding.DecodeString(splitJWT[1])
	if err != nil {
		t.Errorf("Invalid original base format")
		return
	}

	if value != string(decrypt) {
		t.Errorf("Invalid value in JWT")
	}
}

func TestParseUserData(t *testing.T) {
	testMap := map[string]any{
		"Email":    "unemail@gmail.com",
		"Password": "MonMotDePasse123!",
	}

	userData, err := ParseUserData(testMap)
	if err != nil {
		t.Errorf("Error during the parse: %v", err)
		return
	}

	if userData.Email != testMap["Email"] {
		t.Errorf("Email before and after the parse are not the same")
		return
	}

	if userData.Password != testMap["Password"] {
		t.Errorf("password before and after the parse are not the same")
		return
	}
}
