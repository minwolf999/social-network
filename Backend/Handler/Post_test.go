package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	model "social-network/Model"
	utils "social-network/Utils"

	"testing"
)

func TestCreatePost(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := utils.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	rr, err := TryRegister(t, db, model.Register{
			Auth: model.Auth{
				Email:           "unemail1@gmail.com",
				Password:        "MonMotDePasse123!",
				ConfirmPassword: "MonMotDePasse123!",
			},
			FirstName: "Jean",
			LastName:  "Dujardin",
			BirthDate: "1990-01-01",
		})
	if err != nil {
		t.Fatal(err)
		return
	}

	// Check the response body is what we expect.
	expected := "Register successfully"
	bodyValue := make(map[string]any)

	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}

	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	rr, err = TryCreatePost(t, db, fmt.Sprint(bodyValue["sessionId"]))
	if err != nil {
		t.Fatal(err)
		return
	}

	// Check the response body is what we expect.
	expected = "Post created successfully"
	bodyValue = make(map[string]any)

	if err = json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}

	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}
}

func TryCreatePost(t *testing.T, db *sql.DB, authorId string) (*httptest.ResponseRecorder, error) {
	// Create a table for testing
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Post (
			Id VARCHAR(36) NOT NULL UNIQUE,
			AuthorId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id"),
			Text VARCHAR(1000) NOT NULL,
			Image VARCHAR(100),
			IsGroup INTEGER NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}

	post := model.Post{
		Id:       "Test",
		AuthorId: authorId,
		Text:     "Test",
	}

	body, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/CreatePost", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePost(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		return nil, fmt.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	return rr, nil
}

func TestGetPost(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := utils.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	rr, err := TryRegister(t, db, model.Register{
		Auth: model.Auth{
			Email:           "unemail2@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	// Check the response body is what we expect.
	expected := "Register successfully"
	bodyValue := make(map[string]any)

	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}

	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	JWT := fmt.Sprint(bodyValue["sessionId"])
	rr, err = TryCreatePost(t, db, JWT)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Check the response body is what we expect.
	expected = "Post created successfully"
	bodyValue = make(map[string]any)

	if err = json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}

	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	post := model.Post{
		Id:       "Test",
		AuthorId: JWT,
		Text:     "Test",
	}

	body, err := json.Marshal(post)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getPost", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePost(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	expected = "Get posts successfuly"
	bodyValue = make(map[string]any)

	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}

	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}
}

func TestParsePostData(t *testing.T) {
	testMap := []map[string]any{
		{
			"AuthorId": "id",
			"Text":     "Hello wold!",
			"IsGroup":  "0",
		},
	}

	userData, err := ParsePostData(testMap)
	if err != nil {
		t.Errorf("Error during the parse: %v", err)
		return
	}

	if userData[0].Text != testMap[0]["Text"] {
		t.Errorf("Text before and after the parse are not the same")
		return
	}

	if userData[0].IsGroup != testMap[0]["IsGroup"] {
		t.Errorf("IsGroup before and after the parse are not the same")
		return
	}

	if userData[0].AuthorId != testMap[0]["AuthorId"] {
		t.Errorf("AuthorId before and after the parse are not the same")
		return
	}
}
