package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	model "social-network/Model"

	"testing"
)

func TestAddFollower(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	rr, err := TryRegister(t, db, model.Register{
		Auth: model.Auth{
			Email:           "test1@gmail.com",
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

	expected := "Register successfully"
	// Check the response body is what we expect.
	bodyValue := make(map[string]any)

	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	JWT := bodyValue["sessionId"]

	rr, err = TryRegister(t, db, model.Register{
		Auth: model.Auth{
			Email:           "unemail7@gmail.com",
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
	bodyValue = make(map[string]any)
	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	var id string
	db.QueryRow("SELECT Id FROM Auth WHERE Email = ?", "unemail7@gmail.com").Scan(&id)

	follow := map[string]any{
		"UserId":     JWT,
		"FollowerId": id,
	}

	body, err := json.Marshal(follow)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/addFollower", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(AddFollower(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected = "Add follower successfully"
	// Check the response body is what we expect.
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

func TestRemoveFollower(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	rr, err := TryRegister(t, db, model.Register{
		Auth: model.Auth{
			Email:           "test1@gmail.com",
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

	expected := "Register successfully"
	// Check the response body is what we expect.
	bodyValue := make(map[string]any)

	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	JWT := bodyValue["sessionId"]

	rr, err = TryRegister(t, db, model.Register{
		Auth: model.Auth{
			Email:           "unemail7@gmail.com",
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
	bodyValue = make(map[string]any)
	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	var id string
	db.QueryRow("SELECT Id FROM Auth WHERE Email = ?", "unemail7@gmail.com").Scan(&id)

	follow := map[string]any{
		"UserId":     JWT,
		"FollowerId": id,
	}

	body, err := json.Marshal(follow)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/addFollower", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(AddFollower(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected = "Add follower successfully"
	// Check the response body is what we expect.
	bodyValue = make(map[string]any)
	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err = http.NewRequest("POST", "/removeFollower", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(RemoveFollower(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected = "Remove follower successfully"
	// Check the response body is what we expect.
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

func TestGetFollowed(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	rr, err := TryRegister(t, db, model.Register{
		Auth: model.Auth{
			Email:           "test1@gmail.com",
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

	expected := "Register successfully"
	// Check the response body is what we expect.
	bodyValue := make(map[string]any)

	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	JWT := bodyValue["sessionId"]

	rr, err = TryRegister(t, db, model.Register{
		Auth: model.Auth{
			Email:           "unemail7@gmail.com",
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
	bodyValue = make(map[string]any)
	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	var id string
	db.QueryRow("SELECT Id FROM Auth WHERE Email = ?", "unemail7@gmail.com").Scan(&id)

	follow := map[string]any{
		"UserId":     JWT,
		"FollowerId": id,
	}

	body, err := json.Marshal(follow)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/addFollower", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(AddFollower(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected = "Add follower successfully"
	// Check the response body is what we expect.
	bodyValue = make(map[string]any)
	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	follow = map[string]any{
		"UserId": JWT,
	}

	body, err = json.Marshal(follow)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err = http.NewRequest("POST", "/getFollowed", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetFollowed(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected = "Get followed successfuly"
	// Check the response body is what we expect.
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

func TestGetFollower(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	rr, err := TryRegister(t, db, model.Register{
		Auth: model.Auth{
			Email:           "test1@gmail.com",
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

	expected := "Register successfully"
	// Check the response body is what we expect.
	bodyValue := make(map[string]any)

	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	JWT1 := bodyValue["sessionId"]

	rr, err = TryRegister(t, db, model.Register{
		Auth: model.Auth{
			Email:           "unemail7@gmail.com",
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
	bodyValue = make(map[string]any)
	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	JWT2 := bodyValue["sessionId"]

	var id string
	db.QueryRow("SELECT Id FROM Auth WHERE Email = ?", "unemail7@gmail.com").Scan(&id)

	follow := map[string]any{
		"UserId":     JWT1,
		"FollowerId": id,
	}

	body, err := json.Marshal(follow)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/addFollower", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(AddFollower(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected = "Add follower successfully"
	// Check the response body is what we expect.
	bodyValue = make(map[string]any)
	if err := json.Unmarshal(rr.Body.Bytes(), &bodyValue); err != nil {
		t.Fatalf("Erreur lors de la réception de la réponse de la requête : %v", err)
		return
	}
	if bodyValue["Success"] != true {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		return
	}

	follow = map[string]any{
		"UserId": JWT2,
	}

	body, err = json.Marshal(follow)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err = http.NewRequest("POST", "/getFollower", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetFollower(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected = "Get followed successfuly"
	// Check the response body is what we expect.
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
