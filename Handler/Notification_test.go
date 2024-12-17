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

func TestGetAllNotifications(t *testing.T) {
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
			Id:              "userid",
			Email:           "unemail@gmail.com",
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

	body, err := json.Marshal(utils.GenerateJWT(userData.Id))
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getAllNotifications", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllNotifications(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Group obtained successfully"
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

func TestGetGroupNotification(t *testing.T) {
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
			Id:              "userid",
			Email:           "unemail@gmail.com",
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

	group := model.Group{
		Id:           "group",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "group test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	data := map[string]any{
		"GroupId": group.Id,
		"UserId":  utils.GenerateJWT(userData.Id),
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getGroupNotification", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetGroupNotification(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Group obtained successfully"
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

func TestGetUserNotification(t *testing.T) {
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
			Id:              "userid",
			Email:           "unemail@gmail.com",
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

	userData2 := model.Register{
		Auth: model.Auth{
			Id:              "userid2",
			Email:           "unemail2@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean2",
		LastName:  "Dujardin2",
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

	group := model.Group{
		Id:           "group",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "group test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	data := map[string]any{
		"UserId":      utils.GenerateJWT(userData.Id),
		"OtherUserId": userData2.Id,
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getUserNotification", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUserNotification(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Group obtained successfully"
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

func TestDeleteAllNotifications(t *testing.T) {
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
			Id:              "userid",
			Email:           "unemail@gmail.com",
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

	notification := model.Notification{
		Id:          "notificationId",
		UserId:      userData.Id,
		Status:      "test",
		Description: "test",
	}

	if err = notification.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	body, err := json.Marshal(utils.GenerateJWT(userData.Id))
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/deleteAllNotifications", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteAllNotifications(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Group obtained successfully"
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

func TestDeleteAllGroupNotifications(t *testing.T) {
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
			Id:              "userid",
			Email:           "unemail@gmail.com",
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

	group := model.Group{
		Id:           "groupId",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}
	if err = group.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	notification := model.Notification{
		Id:          "notificationId",
		UserId:      userData.Id,
		Status:      "test",
		Description: "test",
		GroupId:     group.Id,
	}

	if err = notification.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	data := map[string]any{
		"UserId":  utils.GenerateJWT(userData.Id),
		"GroupId": group.Id,
	}
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/deleteAllGroupNotifications", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteAllGroupNotifications(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Group obtained successfully"
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

func TestDeleteAllUserNotifications(t *testing.T) {
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
			Id:              "userid",
			Email:           "unemail@gmail.com",
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

	userData2 := model.Register{
		Auth: model.Auth{
			Id:              "userid2",
			Email:           "unemail2@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean2",
		LastName:  "Dujardin2",
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

	notification := model.Notification{
		Id:          "notificationId",
		UserId:      userData.Id,
		Status:      "test",
		Description: "test",
		OtherUserId: userData2.Id,
	}

	if err = notification.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	data := map[string]any{
		"UserId":      utils.GenerateJWT(userData.Id),
		"OtherUserId": userData2.Id,
	}
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/deleteAllUserNotifications", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteAllUserNotifications(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Group obtained successfully"
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
