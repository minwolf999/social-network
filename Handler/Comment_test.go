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

func TestCreateComment(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	var userData = model.Register{
		Auth: model.Auth{
			Id:              "userid",
			Email:           "unemail3@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	JWT := utils.GenerateJWT(userData.Id)

	post := model.Post{
		Id:           "postid",
		AuthorId:     userData.Id,
		Text:         "Test",
		CreationDate: "1970-01-01",
		Status:       "public",
	}

	if err = post.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	rr, err := TryCreateComment(t, db, JWT, post.Id)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Check the response body is what we expect.
	expected := "Comment created successfully"
	bodyValue := make(map[string]any)

	if err = json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}

	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}
}

func TryCreateComment(t *testing.T, db *sql.DB, authorId, postId string) (*httptest.ResponseRecorder, error) {
	post := model.Comment{
		PostId:       postId,
		AuthorId:     authorId,
		Text:         "Test",
		CreationDate: "1970-01-01",
	}

	body, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/createComment/Test", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateComment(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		return nil, fmt.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	return rr, nil
}

func TestGetComment(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	userData := model.Register{
		Auth: model.Auth{
			Id: "userid",
			Email:           "unemail4@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	post := model.Post{
		Id:           "postid",
		AuthorId:     userData.Id,
		Text:         "Test",
		CreationDate: "1970-01-01",
		Status:       "public",
	}

	if err = post.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	comment := model.Comment{
		Id:           "commentid",
		PostId:       post.Id,
		AuthorId:     userData.Id,
		Text:         "Test",
		CreationDate: "1970-01-01",
	}

	if err = comment.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	comment.AuthorId = utils.GenerateJWT(comment.AuthorId)

	body, err := json.Marshal(comment)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getComment/Test", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetComment(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	expected := "Get comments successfuly"
	bodyValue := make(map[string]any)

	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}

	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}
}
