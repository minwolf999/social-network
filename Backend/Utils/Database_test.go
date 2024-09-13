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
	}
	defer db.Close()

	// Check that the connection is not null
	if db == nil {
		t.Fatalf("La connexion à la base de données est nulle")
	}

	// Checks that a simple query can be executed (sanity check)
	err = db.Ping()
	if err != nil {
		t.Fatalf("Impossible de ping la base de données : %v", err)
	}
}

func TestInsertIntoDb(t *testing.T) {
	// Opens an in-memory SQLite database
	db, err := OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
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
	}

	// Calling the InsertIntoDb function to insert data
	err = InsertIntoDb("TestTable", db, "29323HDY73", "John Doe", "JAimeCoder1234")
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
	}

	// Checking that the data has been inserted correctly
	var id string
	var email string
	var password string
	err = db.QueryRow("SELECT Id, Email, Password FROM TestTable WHERE Email = ?", "John Doe").Scan(&id, &email, &password)
	if err != nil {
		t.Fatalf("Erreur lors de la récupération des données : %v", err)
	}

	// Checks
	if email != "John Doe" {
		t.Errorf("Nom attendu 'John Doe', obtenu: %s", email)
	}
	if password != "JAimeCoder1234" {
		t.Errorf("Password attendu JAimeCoder1234, obtenu: %s", password)
	}
}

func TestPrepareStmt(t *testing.T) {
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
	_, err = db.Exec(`INSERT INTO TestTable (Id, Email, Password) VALUES ("019169b0-1302-71ec-a8d5-2615142a12b9","superemail@gmail.com", "JAimeCoder1235"), ("019169b0-1302-71ec-a8d5-2615142a12b9","superemail@gmail.com", "JAimeCoder1234")`)
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
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
	}
	defer rows.Close()

	// Check that the columns are correct
	expectedColumns := []string{"Id", "Email", "Password"}
	for i, col := range expectedColumns {
		if columns[i] != col {
			t.Errorf("Colonne attendue %s, obtenu %s", col, columns[i])
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
		}

		if email != "superemail@gmail.com" {
			t.Errorf("Email attendu 'superemail@gmail.com', obtenu: %s", email)
		}
		if password != "JAimeCoder1234" {
			t.Errorf("Password attendu JAimeCoder1234, obtenu: %s", password)
		}
	} else {
		t.Fatalf("Aucun résultat trouvé pour la requête")
	}

}

func TestSelectFromDb(t *testing.T) {
	// Opens a database in memory
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
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
	}

	// Insert test data
	_, err = db.Exec(`INSERT INTO TestTable (Id, Email, Password) VALUES 
		("1", "superemail@gmail.com", "JAimeCoder1235"), 
		("2", "superemail@gmail.com", "JAimeCoder1234")`)
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
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
	}

	// Check that we got only one line
	if len(result) != 1 {
		t.Fatalf("Nombre de lignes attendu : 1, obtenu : %d", len(result))
	}

	// Checks column values
	row := result[0]
	id := *(row["Id"].(*string))
	email := *(row["Email"].(*string))
	password := *(row["Password"].(*string))

	// Check that the data is correct
	if id != "2" {
		t.Errorf("Id attendu : '2', obtenu : '%s'", id)
	}
	if email != "superemail@gmail.com" {
		t.Errorf("Email attendu : 'superemail@gmail.com', obtenu : '%s'", email)
	}
	if password != "JAimeCoder1234" {
		t.Errorf("Mot de passe attendu : 'JAimeCoder1234', obtenu : '%s'", password)
	}
}
