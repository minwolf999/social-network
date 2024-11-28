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

	"golang.org/x/crypto/bcrypt"
)

func CreateTables(db *sql.DB) {
	db.Exec(`
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
			ProfilePicture TEXT,
			Banner TEXT,
			Username VARCHAR(50),
			AboutMe VARCHAR(280),
			Status VARCHAR(20),
			GroupsJoined TEXT,
		
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
			Image Text,
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
			FollowerId VARCHAR(36) NOT NULL,
			FollowedId VARCHAR(36) NOT NULL,

			PRIMARY KEY (Id),

			CONSTRAINT fk_followerid FOREIGN KEY (FollowerId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_followedid FOREIGN KEY (FollowedId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS Groups (
			Id VARCHAR(36) NOT NULL,
			LeaderId VARCHAR(36) NOT NULL,
			MemberIds TEXT NOT NULL,
			GroupName VARCHAR(200) NOT NULL,
			GroupDescription VARCHAR(500),
			CreationDate VARCHAR(20) NOT NULL,
			Banner TEXT,
			GroupPicture TEXT,

			PRIMARY KEY (Id),

			CONSTRAINT fk_leaderid FOREIGN KEY (LeaderId) REFERENCES "UserInfo"("Id")	
		);

		CREATE TABLE IF NOT EXISTS Event (
			Id VARCHAR(36),
			GroupId VARCHAR(36),
			OrganisatorId VARCHAR(36),
			Title VARCHAR(200),
			Description VARCHAR(1000),
			DateOfTheEvent VARCHAR(20),
			Image TEXT,

			PRIMARY KEY (Id),

			CONSTRAINT fk_groupid FOREIGN KEY (GroupId) REFERENCES "Groups"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_organisatorid FOREIGN KEY (OrganisatorId) REFERENCES "UserInfo"("Id")
		);

		CREATE TABLE IF NOT EXISTS JoinEvent (
			EventId VARCHAR(36),
			UserId VARCHAR(36),

			CONSTRAINT fk_eventid FOREIGN KEY (EventId) REFERENCES "Event"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS DeclineEvent (
			EventId VARCHAR(36),
			UserId VARCHAR(36),
		
			CONSTRAINT fk_eventid FOREIGN KEY (EventId) REFERENCES "Event"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
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

		CREATE VIEW IF NOT EXISTS EventDetail AS
		  SELECT 
		    e.Id,
		    g.GroupName,

		    CASE 
		      WHEN u1.Username = '' THEN CONCAT(u1.FirstName, ' ', u1.LastName)
		      ELSE u1.Username 
		    END AS Organisator,

		    e.Title,
		    e.Description,
		    e.DateOfTheEvent,

		    GROUP_CONCAT(DISTINCT CASE 
		        WHEN u2.Username = '' THEN CONCAT(u2.FirstName, ' ', u2.LastName)
		        ELSE u2.Username 
		    END) AS JoinUsers,

		    GROUP_CONCAT(DISTINCT CASE
		        WHEN u3.Username = '' THEN CONCAT(u3.FirstName, ' ', u3.LastName)
		        ELSE u3.Username
		    END) AS DeclineUsers


		FROM Event AS e
		INNER JOIN Groups AS g ON g.Id = e.GroupId
		INNER JOIN UserInfo AS u1 ON u1.Id = e.OrganisatorId

		LEFT JOIN JoinEvent AS j ON j.EventId = e.Id
		LEFT JOIN UserInfo AS u2 ON u2.Id = j.UserId

		LEFT JOIN DeclineEvent AS d ON d.EventId = e.Id
		LEFT JOIN UserInfo AS u3 ON u3.Id = d.UserId;


		CREATE TABLE IF NOT EXISTS FollowingRequest (
			FollowerId VARCHAR(36) NOT NULL,
			FollowedId VARCHAR(36) NOT NULL,

			CONSTRAINT fk_followerid FOREIGN KEY (FollowerId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_followedid FOREIGN KEY (FollowedId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS JoinGroupRequest (
			UserId VARCHAR(36) NOT NULL,
			GroupId VARCHAR(36) NOT NULL,

			CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_followerid FOREIGN KEY (GroupId) REFERENCES "Groups"("Id") ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS InviteGroupRequest (
			SenderId VARCHAR(36) NOT NULL,
			GroupId VARCHAR(36) NOT NULL,
			ReceiverId VARCHAR(36) NOT NULL,

			CONSTRAINT fk_senderid FOREIGN KEY (SenderId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_followerid FOREIGN KEY (GroupId) REFERENCES "Groups"("Id") ON DELETE CASCADE,
			CONSTRAINT fk_receiverid FOREIGN KEY (ReceiverId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);

		CREATE VIEW IF NOT EXISTS FollowDetail AS
		SELECT 
			f.Id,
			f.FollowerId AS FollowerId,
			User.Username AS Follower_Username,
			
			CASE
			WHEN CONCAT(User.FirstName, ' ', User.LastName) = ' ' THEN ''
			ELSE CONCAT(User.FirstName, ' ', User.LastName)
			END AS Follower_Name,

			f.FollowedId AS FollowedId,
			Follower.Username AS Followed_Username,

			CASE
			WHEN CONCAT(Follower.FirstName, ' ', Follower.LastName) = ' ' THEN ''
			ELSE CONCAT(Follower.FirstName, ' ', Follower.LastName)
			END AS Followed_Name

		FROM Follower AS f
		INNER JOIN UserInfo AS User ON User.Id = f.FollowerId
		INNER JOIN UserInfo AS Follower ON Follower.Id = f.FollowedId;

		CREATE VIEW IF NOT EXISTS GroupDetail AS
		SELECT 
			g.Id,
			g.LeaderId,
			
			CASE 
			WHEN u.Username = '' THEN CONCAT(u.FirstName, ' ', u.LastName)
			ELSE u.Username 
			END AS Leader,

			g.MemberIds,
			g.groupName,
			g.CreationDate

		FROM Groups AS g
		INNER JOIN UserInfo AS u ON u.Id = g.LeaderId;

		CREATE VIEW IF NOT EXISTS FollowRequestDetail AS
		SELECT
			f.FollowerId AS FollowerId,
			CASE 
				WHEN User.Username = '' THEN CONCAT(User.FirstName, ' ', User.LastName)
				ELSE User.Username 
			END AS Follower_Name,

			f.FollowedId AS FollowedId,
			CASE 
				WHEN Follower.Username = '' THEN CONCAT(Follower.FirstName, ' ', Follower.LastName)
				ELSE Follower.Username 
			END AS Followed_Name
			

		FROM FollowingRequest AS f
		INNER JOIN UserInfo AS User ON User.Id = f.FollowerId
		INNER JOIN UserInfo AS Follower ON Follower.Id = f.FollowedId;

		CREATE View IF NOT EXISTS JoinGroupRequestDetail AS
			SELECT
				j.UserId,
				
				CASE 
					WHEN u.Username = '' THEN CONCAT(u.FirstName, ' ', u.LastName)
					ELSE u.Username 
				END AS User_Name,

				j.GroupId,
				g.GroupName

			FROM JoinGroupRequest AS j
			INNER JOIN UserInfo AS u ON u.Id = j.UserId
			INNER JOIN Groups AS g ON g.Id = j.GroupId;

		CREATE View IF NOT EXISTS InviteGroupRequestDetail AS
		SELECT
			i.SenderId,
			
			CASE 
				WHEN Sender.Username = '' THEN CONCAT(Sender.FirstName, ' ', Sender.LastName)
				ELSE Sender.Username 
			END AS Sender_Name,

			i.GroupId,
			g.GroupName,

			i.ReceiverId,
			
			CASE
				WHEN Receiver.Username = '' THEN CONCAT(Receiver.FirstName, ' ', Receiver.LastName)
				ELSE Receiver.Username 
			END AS Receiver_Name
			
		FROM InviteGroupRequest AS i
		INNER JOIN UserInfo AS Sender ON Sender.Id = i.SenderId
		INNER JOIN Groups AS g ON g.Id = i.GroupId
		INNER JOIN UserInfo AS Receiver ON Receiver.Id = i.ReceiverId;

		CREATE TABLE IF NOT EXISTS Notification (
			Id VARCHAR(36) NOT NULL,
			UserId VARCHAR(36) NOT NULL,
			Status VARCHAR(100) NOT NULL,
			Description VARCHAR(256) NOT NULL,
			
			GroupId VARCHAR(36),
			OtherUserId VARCHAR(36),

			CONSTRAINT fk_userid FOREIGN KEY (UserId) REFERENCES "UserInfo"("Id") ON DELETE CASCADE
		);
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

func TestRegisterVerification(t *testing.T) {
	tests := []struct {
		name       string
		data       model.Register
		shouldFail bool
	}{
		{
			name: "Valid registration",
			data: model.Register{
				Auth: model.Auth{
					Email:           "unemail@gmail.com",
					Password:        "zXYVhVxp9zxP8qa$",
					ConfirmPassword: "zXYVhVxp9zxP8qa$",
				},
				FirstName: "Jean",
				LastName:  "Dujardin",
				BirthDate: "1998-01-03",
			},
			shouldFail: false,
		},
		{
			name: "Password and confirm password do not match",
			data: model.Register{
				Auth: model.Auth{
					Email:           "unemail@gmail.com",
					Password:        "zXYVhVxp9zxP8qa$",
					ConfirmPassword: "differentPassword",
				},
				FirstName: "Jean",
				LastName:  "Dujardin",
				BirthDate: "1998-01-03",
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterVerification(tt.data)
			if (err != nil) != tt.shouldFail {
				t.Fatalf("Test '%s' échoué : attendu erreur: %v, obtenu: %v", tt.name, tt.shouldFail, err != nil)
				return
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		shouldFail bool
	}{
		{
			name:       "Short Password",
			data:       "Ey$21",
			shouldFail: true,
		},

		{
			name:       "Contains Uppercase, No Special Char",
			data:       "IFBSOSNHFBJ",
			shouldFail: true,
		},
		{
			name:       "Contains Number, No Special Char ",
			data:       "IDBF2847492",
			shouldFail: true,
		},
		{
			name:       "Password Valide",
			data:       "zXYVhVxp9@P8qa",
			shouldFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := IsValidPassword(tt.data)
			if valid == tt.shouldFail {
				t.Fatalf("Test '%s' échoué : attendu erreur: %v, obtenu: %v", tt.name, tt.shouldFail, !valid)
				return
			}
		})
	}
}

func TestCreateUuidAndCrypt(t *testing.T) {
	// Créer un modèle Register de test
	register := &model.Register{
		Auth: model.Auth{
			Email:    "unemail@gmail.com",
			Password: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	// Appeler la fonction CreateUuidAndCrypt
	err := CreateUuidAndCrypt(register)
	if err != nil {
		t.Fatalf("Erreur lors de l'exécution de CreateUuidAndCrypt: %v", err)
		return
	}

	// Vérifier que le mot de passe a bien été crypté
	if err := bcrypt.CompareHashAndPassword([]byte(register.Auth.Password), []byte("MonMotDePasse123!")); err != nil {
		t.Errorf("Le mot de passe crypté ne correspond pas au mot de passe original")
		return
	}

	// Vérifier que l'UUID a bien été généré
	if register.Auth.Id == "" {
		t.Errorf("L'UUID n'a pas été généré correctement")
		return
	}
}
