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

func TestHandleLike(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	var register = model.Register{
		Auth: model.Auth{
			Id:       "userId",
			Email:    "email",
			Password: "password",
		},

		FirstName: "firstname",
		LastName:  "lastname",
		BirthDate: "monday",
	}

	if err = register.Auth.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	if err = register.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	
	var post = model.Post{
		Id:           "postId",
		AuthorId:     register.Id,
		Text:         "text",
		CreationDate: "now",
		Status:       "public",
	}
	if err = post.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}
	
	register.Id = utils.GenerateJWT(register.Id)

	var like = map[string]any {
		"PostID": post.Id,
		"UserID": register.Id,
		"Table": "LikePost",
	}

	body, err := json.Marshal(like)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
		return
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/like", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
		return
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleLike(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		return
	}

	// Check the response body is what we expect.
	expected := "Like handled successfully"
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

func TestHandleLikeLogic(t *testing.T) {
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	var register = model.Register{
		Auth: model.Auth{
			Id:       "userId",
			Email:    "email",
			Password: "password",
		},

		FirstName: "firstname",
		LastName:  "lastname",
		BirthDate: "monday",
	}

	if err = register.Auth.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	if err = register.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	var post = model.Post{
		Id:           "postId",
		AuthorId:     "userId",
		Text:         "text",
		CreationDate: "now",
		Status:       "public",
	}
	if err = post.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	if err = handleLikeLogic(db, "LikePost", "postId", "userId"); err != nil {
		t.Fatalf("error during the function : %v", err)
		return
	}
}

func TestHasUserLikedPost(t *testing.T) {
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	var register = model.Register{
		Auth: model.Auth{
			Id:       "userId",
			Email:    "email",
			Password: "password",
		},

		FirstName: "firstname",
		LastName:  "lastname",
		BirthDate: "monday",
	}

	if err = register.Auth.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	if err = register.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	var post = model.Post{
		Id:           "postId",
		AuthorId:     "userId",
		Text:         "text",
		CreationDate: "now",
		Status:       "public",
	}
	if err = post.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	if err = addLike(db, "LikePost", "postId", "userId"); err != nil {
		t.Fatalf("error while adding like : %v", err)
		return
	}

	bo, err := hasUserLikedPost(db, "LikePost", "postId", "userId")
	if err != nil {
		t.Fatalf("error during the function : %v", err)
	}

	if !bo {
		t.Fatal("The result is not the good")
	}
}

func TestAddLikeAndRemove(t *testing.T) {
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	var register = model.Register{
		Auth: model.Auth{
			Id:       "userId",
			Email:    "email",
			Password: "password",
		},

		FirstName: "firstname",
		LastName:  "lastname",
		BirthDate: "monday",
	}

	if err = register.Auth.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	if err = register.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	var post = model.Post{
		Id:           "postId",
		AuthorId:     "userId",
		Text:         "text",
		CreationDate: "now",
		Status:       "public",
	}
	if err = post.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	if err = addLike(db, "LikePost", "postId", "userId"); err != nil {
		t.Fatalf("error while adding like : %v", err)
		return
	}

	if err = removeLike(db, "LikePost", "postId", "userId"); err != nil {
		t.Fatalf("error while removing like : %v", err)
		return
	}
}

func TestUpdateLikeCount(t *testing.T) {
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	CreateTables(db)

	var register = model.Register{
		Auth: model.Auth{
			Id:       "userId",
			Email:    "email",
			Password: "password",
		},

		FirstName: "firstname",
		LastName:  "lastname",
		BirthDate: "monday",
	}

	if err = register.Auth.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	if err = register.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	var post = model.Post{
		Id:           "hello",
		AuthorId:     "userId",
		Text:         "text",
		CreationDate: "now",
		Status:       "public",
	}
	if err = post.InsertIntoDb(db); err != nil {
		t.Fatalf("%v", err)
		return
	}

	if err = updateLikeCount(db, "Post", "hello", 1); err != nil {
		t.Fatalf("Error during the function : %v", err)
		return
	}
}
