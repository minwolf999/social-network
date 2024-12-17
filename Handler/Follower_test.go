package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	model "social-network/Model"
	utils "social-network/Utils"

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

	var userData1 = model.Register{
		Auth: model.Auth{
			Id:              "userid1",
			Email:           "test1@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData1.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData1.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	JWT := utils.GenerateJWT(userData1.Id)

	var userData2 = model.Register{
		Auth: model.Auth{
			Id:              "userid2",
			Email:           "unemail7@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData2.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData2.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	follow := map[string]any{
		"FollowerId": JWT,
		"FollowedId": userData2.Id,
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
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddFollower(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected := "Add follower successfully"
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

	var userData1 = model.Register{
		Auth: model.Auth{
			Id:              "userid1",
			Email:           "test1@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData1.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData1.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	JWT := utils.GenerateJWT(userData1.Id)

	userData2 := model.Register{
		Auth: model.Auth{
			Id:              "userid2",
			Email:           "unemail7@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData2.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData2.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	follow := model.Follower{
		Id:         "followid",
		FollowerId: userData1.Id,
		FollowedId: userData2.Id,
	}

	if err = follow.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	follow.FollowerId = JWT
	body, err := json.Marshal(follow)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/removeFollower", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RemoveFollower(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected := "Remove follower successfully"
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

	userData1 := model.Register{
		Auth: model.Auth{
			Id:              "userid1",
			Email:           "test1@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData1.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData1.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	JWT := utils.GenerateJWT(userData1.Id)

	userData2 := model.Register{
		Auth: model.Auth{
			Id:              "userid2",
			Email:           "unemail7@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData2.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData2.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	follow := model.Follower{
		Id:         "followid",
		FollowerId: userData1.Id,
		FollowedId: userData2.Id,
	}

	if err = follow.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	var follower struct {
		UserId      string `json:"UserId"`
		OtherUserId string `json:"OtherUserId"`
	}
	follower.UserId = JWT
	follower.OtherUserId = ""

	body, err := json.Marshal(follower)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getFollowed", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFollowed(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected := "Get followed successfuly"
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

	userData1 := model.Register{
		Auth: model.Auth{
			Id:              "userid1",
			Email:           "test1@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData1.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData1.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	// JWT1 := bodyValue["sessionId"]

	userData2 := model.Register{
		Auth: model.Auth{
			Id:              "userid2",
			Email:           "unemail7@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	if err = userData2.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData2.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	// JWT2 := bodyValue["sessionId"]

	follow := model.Follower{
		Id:         "followerid",
		FollowerId: userData1.Id,
		FollowedId: userData2.Id,
	}

	if err = follow.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	var follower struct {
		UserId      string `json:"UserId"`
		OtherUserId string `json:"OtherUserId"`
	}

	follower.UserId = utils.GenerateJWT(userData2.Id)

	body, err := json.Marshal(follower)
	if err != nil {
		t.Fatal(err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getFollower", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFollower(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Fatal(err)
		return
	}

	expected := "Get followed successfuly"
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
}
