package utils

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestOpenDb(t *testing.T) {
	// Opens an in-memory SQLite database
	db, err := OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
		return
	}
	defer db.Close()

	// Check that the connection is not null
	if db == nil {
		t.Fatalf("La connexion à la base de données est nulle")
		return
	}

	// Checks that a simple query can be executed (sanity check)
	err = db.Ping()
	if err != nil {
		t.Fatalf("Impossible de ping la base de données : %v", err)
		return
	}
}

/* func TestLoadData(t *testing.T) {
	// Opens an in-memory SQLite database
	db, err := OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
		return
	}
	defer db.Close()

	// Create a table for testing
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS Auth (
			Id VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY,
			Email VARCHAR(100) NOT NULL UNIQUE,
			Password VARCHAR(50) NOT NULL,
			ConnectionAttempt INTEGER
		);
			
		CREATE TABLE IF NOT EXISTS UserInfo (
			Id VARCHAR(36) NOT NULL UNIQUE REFERENCES "Auth"("Id"),
			Email VARCHAR(100) NOT NULL UNIQUE REFERENCES "Auth"("Email"),
			FirstName VARCHAR(50) NOT NULL, 
			LastName VARCHAR(50) NOT NULL,
			BirthDate VARCHAR(20) NOT NULL,
			ProfilePicture VARCHAR(400000),
			Username VARCHAR(50),
			AboutMe VARCHAR(280)
		);
			
		CREATE TABLE IF NOT EXISTS Post (
			Id VARCHAR(36) NOT NULL UNIQUE,
			AuthorId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id"),
			Text VARCHAR(1000) NOT NULL,
			Image VARCHAR(100),
			CreationDate VARCHAR(20) NOT NULL,
			IsGroup VARCHAR(36) REFERENCES "Groups"("Id"),
			LikeCount INTEGER,
			DislikeCount INTEGER
		);
			
		CREATE TABLE IF NOT EXISTS Comment (
			Id VARCHAR(36) NOT NULL UNIQUE,
			AuthorId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id"),
			Text VARCHAR(1000) NOT NULL,
			CreationDate VARCHAR(20) NOT NULL,

			PostId VARCHAR(36) REFERENCES "Post"("Id"),

			LikeCount INTEGER,
			DislikeCount INTEGER
		);
			
		CREATE TABLE IF NOT EXISTS Groups (
			Id VARCHAR(36) NOT NULL UNIQUE
		);
		`)
	if err != nil {
		t.Fatalf("Erreur lors de la création de la table : %v", err)
		return
	}

	err = LoadData(db)
	if err != nil {
		t.Fatalf("Erreur pendant la résolution de la fonction : %v", err)
		return
	}
}*/

func TestInsertIntoDb(t *testing.T) {
	// Opens an in-memory SQLite database
	db, err := OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
		return
	}
	defer db.Close()

	// Create a table for testing
	_, err = db.Exec(`
		CREATE TABLE TestTable (
			Id TEXT,
			Email TEXT,
			Password INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("Erreur lors de la création de la table : %v", err)
		return
	}

	// Calling the InsertIntoDb function to insert data
	err = InsertIntoDb("TestTable", db, "29323HDY73", "John Doe", "JAimeCoder1234")
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
		return
	}

	// Checking that the data has been inserted correctly
	var id string
	var email string
	var password string
	err = db.QueryRow("SELECT Id, Email, Password FROM TestTable WHERE Email = ?", "John Doe").Scan(&id, &email, &password)
	if err != nil {
		t.Fatalf("Erreur lors de la récupération des données : %v", err)
		return
	}

	// Checks
	if email != "John Doe" {
		t.Errorf("Nom attendu 'John Doe', obtenu: %s", email)
		return
	}
	if password != "JAimeCoder1234" {
		t.Errorf("Password attendu JAimeCoder1234, obtenu: %s", password)
		return
	}
}

func TestPrepareStmt(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
		return
	}
	defer db.Close()

	// Create a test table
	_, err = db.Exec(`
		CREATE TABLE TestTable (
			Id TEXT,
			Email VARCHAR(50),
			Password TEXT
		)
	`)
	if err != nil {
		t.Fatalf("Erreur lors de la création de la table : %v", err)
		return
	}

	// Insert test data
	_, err = db.Exec(`INSERT INTO TestTable (Id, Email, Password) VALUES ("019169b0-1302-71ec-a8d5-2615142a12b9","superemail@gmail.com", "JAimeCoder1235"), ("019169b0-1302-71ec-a8d5-2615142a12b9","superemail@gmail.com", "JAimeCoder1234")`)
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
		return
	}

	// Calling the PrepareStmt function with test arguments
	args := map[string]any{
		"Id":       "019169b0-1302-71ec-a8d5-2615142a12b9",
		"Email":    "superemail@gmail.com",
		"Password": "JAimeCoder1234",
	}

	columns, rows, err := PrepareStmt("TestTable", db, args)
	if err != nil {
		t.Fatalf("Erreur lors de l'exécution de PrepareStmt : %v", err)
		return
	}
	defer rows.Close()

	// Check that the columns are correct
	expectedColumns := []string{"Id", "Email", "Password"}
	for i, col := range expectedColumns {
		if columns[i] != col {
			t.Errorf("Colonne attendue %s, obtenu %s", col, columns[i])
			return
		}
	}

	// Check that the results are correct
	var id string
	var email string
	var password string
	if rows.Next() {
		err = rows.Scan(&id, &email, &password)
		if err != nil {
			t.Fatalf("Erreur lors de la lecture des résultats : %v", err)
			return
		}

		if email != "superemail@gmail.com" {
			t.Errorf("Email attendu 'superemail@gmail.com', obtenu: %s", email)
			return
		}
		if password != "JAimeCoder1234" {
			t.Errorf("Password attendu JAimeCoder1234, obtenu: %s", password)
			return
		}
	} else {
		t.Fatalf("Aucun résultat trouvé pour la requête")
		return
	}

}

func TestSelectFromDb(t *testing.T) {
	// Opens a database in memory
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
		return
	}
	defer db.Close()

	// Create a test table
	_, err = db.Exec(`
		CREATE TABLE TestTable (
			Id TEXT,
			Email TEXT,
			Password TEXT
		)
	`)
	if err != nil {
		t.Fatalf("Erreur lors de la création de la table : %v", err)
		return
	}

	// Insert test data
	_, err = db.Exec(`INSERT INTO TestTable (Id, Email, Password) VALUES 
		("1", "superemail@gmail.com", "JAimeCoder1235"), 
		("2", "superemail@gmail.com", "JAimeCoder1234")`)
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
		return
	}

	// Arguments for selection (example with Email and Password)
	args := map[string]any{
		"Email":    "superemail@gmail.com",
		"Password": "JAimeCoder1234",
	}

	// Calling the SelectFromDb function
	result, err := SelectFromDb("TestTable", db, args)
	if err != nil {
		t.Fatalf("Erreur lors de l'exécution de SelectFromDb : %v", err)
		return
	}

	// Check that we got only one line
	if len(result) != 1 {
		t.Fatalf("Nombre de lignes attendu : 1, obtenu : %d", len(result))
		return
	}

	// Checks column values
	row := result[0]
	id := *(row["Id"].(*string))
	email := *(row["Email"].(*string))
	password := *(row["Password"].(*string))

	// Check that the data is correct
	if id != "2" {
		t.Errorf("Id attendu : '2', obtenu : '%s'", id)
		return
	}
	if email != "superemail@gmail.com" {
		t.Errorf("Email attendu : 'superemail@gmail.com', obtenu : '%s'", email)
		return
	}
	if password != "JAimeCoder1234" {
		t.Errorf("Mot de passe attendu : 'JAimeCoder1234', obtenu : '%s'", password)
		return
	}
}

func TestPrepareUpdateStmt(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create test table
	_, err = db.Exec(`
		CREATE TABLE TestTable (
			Id TEXT PRIMARY KEY,
			Email VARCHAR(50),
			Password TEXT
		)
	`)
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

	// Insert test data
	_, err = db.Exec(`INSERT INTO TestTable (Id, Email, Password) VALUES 
		("019169b0-1302-71ec-a8d5-2615142a12b9", "superemail@gmail.com", "JAimeCoder1234"),
		("119169b0-1302-71ec-a8d5-2615142a12b9", "anotheremail@gmail.com", "Password5678")`)
	if err != nil {
		t.Fatalf("Error inserting data: %v", err)
	}

	// Prepare update statement
	args := map[string]any{
		"Id":       "019169b0-1302-71ec-a8d5-2615142a12b9",
		"Email":    "updatedemail@gmail.com",
		"Password": "NewPassword9876",
	}
	colsToUpdate := []string{"Email", "Password"}

	err = PrepareUpdateStmt("TestTable", db, args, colsToUpdate)
	if err != nil {
		t.Fatalf("Error executing PrepareUpdateStmt: %v", err)
	}

	// Check that the update was successful
	var email, password string
	err = db.QueryRow("SELECT Email, Password FROM TestTable WHERE Id = ?", args["Id"]).Scan(&email, &password)
	if err != nil {
		t.Fatalf("Error verifying results: %v", err)
	}

	if email != "updatedemail@gmail.com" {
		t.Errorf("Expected email 'updatedemail@gmail.com', got: %s", email)
	}
	if password != "NewPassword9876" {
		t.Errorf("Expected password 'NewPassword9876', got: %s", password)
	}
}

func TestUpdateDb(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer db.Close()

	// Create a test table
	_, err = db.Exec(`
		CREATE TABLE TestTable (
			Id TEXT,
			Email VARCHAR(50),
			Password TEXT
		)
	`)
	if err != nil {
		t.Fatalf("Erreur lors de la création de la table : %v", err)
	}

	// Insert test data
	_, err = db.Exec(`INSERT INTO TestTable (Id, Email, Password) VALUES 
		("019169b0-1302-71ec-a8d5-2615142a12b9","superemail@gmail.com", "JAimeCoder1234"),
		("119169b0-1302-71ec-a8d5-2615142a12b9","anotheremail@gmail.com", "Password5678")`)
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
	}

	// Calling the UpdateInDb function with test arguments
	updateArgs := map[string]any{
		"Email":    "newemail@gmail.com",
		"Password": "UpdatedPass1234",
	}
	whereArgs := map[string]any{
		"Id": "019169b0-1302-71ec-a8d5-2615142a12b9",
	}

	err = UpdateDb("TestTable", db, updateArgs, whereArgs)
	if err != nil {
		t.Fatalf("Erreur lors de l'exécution de UpdateDb : %v", err)
	}

	// Check that the update was successful
	var email, password string
	err = db.QueryRow("SELECT Email, Password FROM TestTable WHERE Id = ?", whereArgs["Id"]).Scan(&email, &password)
	if err != nil {
		t.Fatalf("Erreur lors de la vérification des résultats : %v", err)
	}

	if email != "newemail@gmail.com" {
		t.Errorf("Email attendu 'newemail@gmail.com', obtenu: %s", email)
	}
	if password != "UpdatedPass1234" {
		t.Errorf("Password attendu UpdatedPass1234, obtenu: %s", password)
	}
}

func TestRemoveFromDb(t *testing.T) {
	// Opens an in-memory SQLite database
	db, err := OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
		return
	}
	defer db.Close()

	// Create a table for testing
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS TestTable (
			Id TEXT,
			Email TEXT,
			Password INTEGER
		);

		CREATE TABLE IF NOT EXISTS Follower (
			Id VARCHAR(36) NOT NULL UNIQUE,
			UserId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id"),
			FollowerId VARCHAR(36) NOT NULL REFERENCES "UserInfo"("Id")
		);
	`)
	if err != nil {
		t.Fatalf("Erreur lors de la création de la table : %v", err)
		return
	}

	// Calling the InsertIntoDb function to insert data
	if err = InsertIntoDb("TestTable", db, "test1", "John Doe1", "JAimeCoder1234"); err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
		return
	}

	// Calling the InsertIntoDb function to insert data
	if err = InsertIntoDb("TestTable", db, "test2", "John Doe2", "JAimeCoder1234"); err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
		return
	}

	// Calling the InsertIntoDb function to insert data
	if err = InsertIntoDb("Follower", db, "id", "Test1", "Test2"); err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
		return
	}

	if err = RemoveFromDB("Follower", db, map[string]any{"UserId": "Test1", "FollowerId": "Test2"}); err != nil {
		t.Fatalf("Erreur lors de la suppression des données : %v", err)
		return
	}
}