package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	model "social-network/Model"
	utils "social-network/Utils"
	"testing"
)

func TestCreateGroup(t *testing.T) {
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

	JWT := utils.GenerateJWT(userData.Id)

	var user = map[string]any{
		"LeaderId":     JWT,
		"GroupName":    "test",
		"CreationDate": "now",
	}

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/createGroup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateGroup(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Group created successfully"
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

func TestLeaveGroup(t *testing.T) {
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
			Email:           "unemail@gmail.com",
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

	userData2 := model.Register{
		Auth: model.Auth{
			Id:              "userid2",
			Email:           "unemail1@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName:    "Jean",
		LastName:     "Dujardin",
		BirthDate:    "1990-01-01",
		GroupsJoined: "groupid",
	}

	if err = userData2.Auth.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	if err = userData2.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData1.Id,
		MemberIds:    fmt.Sprintf("%s | %s", userData1.Id, userData2.Id),
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	Leave := map[string]any{
		"UserId":  utils.GenerateJWT(userData2.Id),
		"GroupId": group.Id,
	}

	body, err := json.Marshal(Leave)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/LeaveGroup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LeaveGroup(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Group joined successfuly"
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

func TestGetGroup(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	getGroup := map[string]any{
		"UserId":  utils.GenerateJWT(userData.Id),
		"groupId": group.Id,
	}

	body, err := json.Marshal(getGroup)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getGroup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetGroup(db))

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

func TestGetAllGroups(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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
	req, err := http.NewRequest("POST", "/getAllGroups", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllGroups(db))

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

func TestGetGroupsJoined(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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
	req, err := http.NewRequest("POST", "/getGroupsJoined", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetGroupsJoined(db))

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

func TestGetGroupsPosts(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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
	req, err := http.NewRequest("POST", "/getGroupsPosts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetGroupsPosts(db))

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

func TestDeleteGroup(t *testing.T) {
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

	JWT := utils.GenerateJWT(userData.Id)

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	user := map[string]any{
		"UserId":  JWT,
		"groupId": group.Id,
	}

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/deleteGroup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteGroup(db))

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

func TestJoinGroup(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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

	data := map[string]any{
		"UserId":  utils.GenerateJWT(userData2.Id),
		"GroupId": group.Id,
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/joinGroup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(JoinGroup(db))

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

func TestGetJoinRequest(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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

	joinRequest := model.JoinGroupRequest{
		UserId:  userData2.Id,
		GroupId: group.Id,
	}

	if err = joinRequest.InsertIntoDb(db); err != nil {
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
	req, err := http.NewRequest("POST", "/getJoinRequest", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetJoinRequest(db))

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

func TestDeclineJoinRequest(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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

	joinRequest := model.JoinGroupRequest{
		UserId:  userData2.Id,
		GroupId: group.Id,
	}

	if err = joinRequest.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	data := map[string]any{
		"UserId":     utils.GenerateJWT(userData.Id),
		"GroupId":    group.Id,
		"JoinUserId": userData2.Id,
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/declineJoinRequest", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeclineJoinRequest(db))

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

func TestAcceptJoinRequest(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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

	joinRequest := model.JoinGroupRequest{
		UserId:  userData2.Id,
		GroupId: group.Id,
	}

	if err = joinRequest.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	data := map[string]any{
		"UserId":     utils.GenerateJWT(userData.Id),
		"GroupId":    group.Id,
		"JoinUserId": userData2.Id,
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/acceptJoinRequest", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AcceptJoinRequest(db))

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

func TestInviteGroup(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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

	data := map[string]any{
		"SenderId":   utils.GenerateJWT(userData.Id),
		"GroupId":    group.Id,
		"ReceiverId": userData2.Id,
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/inviteGroup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InviteGroup(db))

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

func TestGetInvitationGroup(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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

	joinRequest := model.InviteGroupRequest{
		SenderId:  userData.Id,
		GroupId: group.Id,
		ReceiverId: userData2.Id,
	}

	if err = joinRequest.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	data := utils.GenerateJWT(userData2.Id)

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getInvitationGroup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetInvitationGroup(db))

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

func TestDeclineInvitationGroup(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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

	joinRequest := model.InviteGroupRequest{
		SenderId:  userData.Id,
		GroupId: group.Id,
		ReceiverId: userData2.Id,
	}

	if err = joinRequest.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	data := map[string]any{
		"SenderId":   userData.Id,
		"GroupId":    group.Id,
		"ReceiverId": utils.GenerateJWT(userData2.Id),
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/declineInvitationGroup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeclineInvitationGroup(db))

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

func TestAcceptInvitationGroup(t *testing.T) {
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

	var group = model.Group{
		Id:           "groupid",
		LeaderId:     userData.Id,
		MemberIds:    userData.Id,
		GroupName:    "test",
		CreationDate: "now",
	}

	if err = group.InsertIntoDb(db); err != nil {
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

	joinRequest := model.InviteGroupRequest{
		SenderId:  userData.Id,
		GroupId: group.Id,
		ReceiverId: userData2.Id,
	}

	if err = joinRequest.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	data := map[string]any{
		"SenderId":   userData.Id,
		"GroupId":    group.Id,
		"ReceiverId": utils.GenerateJWT(userData2.Id),
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/acceptInvitationGroup", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AcceptInvitationGroup(db))

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
