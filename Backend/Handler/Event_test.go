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

func TestCreateEvent(t *testing.T) {
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
			Id:       "userid",
			Email:    "email",
			Password: "password",
		},
		FirstName: "name",
		LastName:  "name",
		BirthDate: "now",
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

	userData.Id = utils.GenerateJWT(userData.Id)

	var event = model.Event{
		OrganisatorId:  userData.Id,
		GroupId:        group.Id,
		Title:          "titre",
		Description:    "description",
		DateOfTheEvent: "now",
	}

	body, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/createEvent", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateEvent(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Event created successfully"
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

func TestJoinEvent(t *testing.T) {
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
			Id:       "userid",
			Email:    "email",
			Password: "password",
		},
		FirstName: "name",
		LastName:  "name",
		BirthDate: "now",
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

	var event = model.Event{
		Id:             "eventid",
		OrganisatorId:  userData.Id,
		GroupId:        group.Id,
		Title:          "titre",
		Description:    "description",
		DateOfTheEvent: "now",
	}

	if err = event.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	var joinEvent = model.JoinEvent{
		EventId: event.Id,
		UserId:  utils.GenerateJWT(userData.Id),
	}

	body, err := json.Marshal(joinEvent)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/joinEvent", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(JoinEvent(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Event joined successfully"
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

func TestDeclineEvent(t *testing.T) {
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
			Id:       "userid",
			Email:    "email",
			Password: "password",
		},
		FirstName: "name",
		LastName:  "name",
		BirthDate: "now",
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

	var event = model.Event{
		Id:             "eventid",
		OrganisatorId:  userData.Id,
		GroupId:        group.Id,
		Title:          "titre",
		Description:    "description",
		DateOfTheEvent: "now",
	}

	if err = event.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	var joinEvent = model.DeclineEvent{
		EventId: event.Id,
		UserId:  utils.GenerateJWT(userData.Id),
	}

	body, err := json.Marshal(joinEvent)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/declineEvent", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeclineEvent(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Event declined successfully"
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

func TestGetEvent(t *testing.T) {
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
			Id:       "userid",
			Email:    "email",
			Password: "password",
		},
		FirstName: "name",
		LastName:  "name",
		BirthDate: "now",
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

	var event = model.Event{
		Id:             "eventid",
		OrganisatorId:  userData.Id,
		GroupId:        group.Id,
		Title:          "titre",
		Description:    "description",
		DateOfTheEvent: "now",
	}

	if err = event.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	var joinEvent = model.JoinEvent{
		EventId: event.Id,
		UserId:  userData.Id,
	}

	if err = joinEvent.InsertIntoDb(db); err != nil {
		t.Fatal(err)
		return
	}

	var data = map[string]any{
		"EventId": event.Id,
		"UserId":  utils.GenerateJWT(userData.Id),
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/getEvent", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetEvent(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Event detail get successfully"
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