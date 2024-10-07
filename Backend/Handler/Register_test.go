package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	model "social-network/Model"
)

func CreateTables(db *sql.DB) {
	db.Exec(`
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS Auth (
			Id VARCHAR(36) NOT NULL,
			Email VARCHAR(100) NOT NULL UNIQUE,
			Password VARCHAR(50) NOT NULL,
			ConnectionAttempt INTEGER,
		
			PRIMARY KEY (Id)
		);
		
		CREATE TABLE IF NOT EXISTS UserInfo (
			Id VARCHAR(36) NOT NULL UNIQUE,
			Email VARCHAR(100) NOT NULL UNIQUE,
			FirstName VARCHAR(50) NOT NULL, 
			LastName VARCHAR(50) NOT NULL,
			BirthDate VARCHAR(20) NOT NULL,
			ProfilePicture VARCHAR(400000),
			Username VARCHAR(50),
			AboutMe VARCHAR(280),
		
			CONSTRAINT fk_id FOREIGN KEY (Id) REFERENCES "Auth"("Id") ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS Post (
		    Id VARCHAR(36) NOT NULL,
		    AuthorId VARCHAR(36) NOT NULL,
		    Text VARCHAR(1000) NOT NULL,
		    Image TEXT,
		    CreationDate VARCHAR(20) NOT NULL,
		    Status TEXT NOT NULL,
		    IsGroup VARCHAR(36),
		    LikeCount INTEGER,
		    DislikeCount INTEGER,
		
			PRIMARY KEY (Id),
		
			CONSTRAINT fk_authorid FOREIGN KEY (AuthorId) REFERENCES "UserInfo"("Id"),
			CONSTRAINT fk_isgroup FOREIGN KEY (IsGroup) REFERENCES "Groups"("Id") ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS LikePost (
			PostId VARCHAR(36) NOT NULL,
			UserId VARCHAR(36) NOT NULL,
		
			CONSTRAINT fk_postid FOREIGN KEY (PostId) REFERENCES "Post"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS DislikePost (
			PostId VARCHAR(36) NOT NULL,
			UserId VARCHAR(36) NOT NULL,
		
			CONSTRAINT fk_postid FOREIGN KEY (PostId) REFERENCES "Post"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS Comment (
			Id VARCHAR(36) NOT NULL,
			AuthorId VARCHAR(36) NOT NULL,
			Text VARCHAR(1000) NOT NULL,
			CreationDate VARCHAR(20) NOT NULL,
			PostId VARCHAR(36),
			LikeCount INTEGER,
			DislikeCount INTEGER,
		
			PRIMARY KEY (Id),
		
			CONSTRAINT fk_authorid FOREIGN KEY (AuthorId) REFERENCES "UserInfo"("Id"),
			CONSTRAINT fk_postid FOREIGN KEY (PostId) REFERENCES "Post"("Id") ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS LikeComment (
			PostId VARCHAR(36) NOT NULL,
			UserId VARCHAR(36) NOT NULL,
		
			CONSTRAINT fk_postid FOREIGN KEY (PostId) REFERENCES "Post"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS DislikeComment (
			PostId VARCHAR(36) NOT NULL,
			UserId VARCHAR(36) NOT NULL,
		
			CONSTRAINT fk_postid FOREIGN KEY (PostId) REFERENCES "Post"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS Follower (
			Id VARCHAR(36) NOT NULL,
			UserId VARCHAR(36) NOT NULL,
			FollowerId VARCHAR(36) NOT NULL,
		
			PRIMARY KEY (Id),
		
			CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_followerid FOREIGN KEY (FollowerId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS Groups (
			Id VARCHAR(36) NOT NULL,
			LeaderId VARCHAR(36) NOT NULL,
			MemberIds TEXT NOT NULL,
			GroupName VARCHAR(200) NOT NULL,
			CreationDate VARCHAR(20) NOT NULL,
		
			PRIMARY KEY (Id),
		
			CONSTRAINT fk_leaderid FOREIGN KEY (LeaderId) REFERENCES "UserInfo"("Id")	
		);
		
		CREATE VIEW PostDetail AS
		  SELECT 
		    p.Id,
			p.Text,
			p.Image,
			p.CreationDate,
			p.IsGroup,
			p.AuthorId,
			p.LikeCount,
			p.DislikeCount,
			u.FirstName,
			u.LastName,
			u.ProfilePicture,
			u.Username
		FROM Post AS p
		INNER JOIN UserInfo AS u ON p.AuthorId = u.Id;
		
		CREATE VIEW CommentDetail AS
		  SELECT 
		    c.Id,
			c.Text,
			c.CreationDate,
			c.AuthorId,
			c.LikeCount,
			c.DislikeCount,
			c.PostId,
			u.FirstName,
			u.LastName,
			u.ProfilePicture,
			u.Username
		FROM Comment AS c
		INNER JOIN UserInfo AS u ON c.AuthorId = u.Id;
	`)
}

func TestRegister(t *testing.T) {
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
			Email:           "unemail@gmail.com",
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
}

func TryRegister(t *testing.T, db *sql.DB, register model.Register) (*httptest.ResponseRecorder, error) {
	// Create a table for testing
	body, err := json.Marshal(register)
	if err != nil {
		return nil, err
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Register(db))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		return nil, fmt.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)

	}

	return rr, nil
}
