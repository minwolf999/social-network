package utils

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestOpenDb(t *testing.T) {
	// Ouvre une base de données SQLite en mémoire
	db, err := OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer db.Close()

	// Vérifie que la connexion n'est pas nulle
	if db == nil {
		t.Fatalf("La connexion à la base de données est nulle")
	}

	// Vérifie qu'on peut exécuter une simple requête (sanity check)
	err = db.Ping()
	if err != nil {
		t.Fatalf("Impossible de ping la base de données : %v", err)
	}
}

func TestCreateDb(t *testing.T) {
	// Ouvre une base de données SQLite en mémoire
	db, err := OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer db.Close()

	// Exécute la création des tables
	CreateDb(db)

	// Test si la table "Auth" a été créée avec succès
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='Auth'").Scan(&tableName)
	if err != nil {
		t.Fatalf("La table Auth n'a pas été créée : %v", err)
	}

	if tableName != "Auth" {
		t.Errorf("Table 'Auth' non trouvée, trouvée: %s", tableName)
	}

	// Test si la table "UserInfo" a été créée avec succès
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='UserInfo'").Scan(&tableName)
	if err != nil {
		t.Fatalf("La table UserInfo n'a pas été créée : %v", err)
	}

	if tableName != "UserInfo" {
		t.Errorf("Table 'UserInfo' non trouvée, trouvée: %s", tableName)
	}
}

func TestInsertIntoDb(t *testing.T) {
	// Ouvre une base de données SQLite en mémoire
	db, err := OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer db.Close()

	// Crée une table pour les tests
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

	// Appel de la fonction InsertIntoDb pour insérer des données
	err = InsertIntoDb("TestTable", db, "29323HDY73", "John Doe", "JAimeCoder1234")
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
	}

	// Vérification que les données ont bien été insérées
	var id string
	var email string
	var password string
	err = db.QueryRow("SELECT Id, Email, Password FROM TestTable WHERE Email = ?", "John Doe").Scan(&id, &email, &password)
	if err != nil {
		t.Fatalf("Erreur lors de la récupération des données : %v", err)
	}

	// Vérifications
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

	// Crée une table de test
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

	// Insère des données de test
	_, err = db.Exec(`INSERT INTO TestTable (Id, Email, Password) VALUES ("019169b0-1302-71ec-a8d5-2615142a12b9","superemail@gmail.com", "JAimeCoder1235"), ("019169b0-1302-71ec-a8d5-2615142a12b9","superemail@gmail.com", "JAimeCoder1234")`)
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
	}

	// Appel de la fonction PrepareStmt avec des arguments de test
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

	// Vérifie que les colonnes sont correctes
	expectedColumns := []string{"Id", "Email", "Password"}
	for i, col := range expectedColumns {
		if columns[i] != col {
			t.Errorf("Colonne attendue %s, obtenu %s", col, columns[i])
		}
	}

	// Vérifie que les résultats sont corrects
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
	// Ouvre une base de données en mémoire
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer db.Close()

	// Crée une table de test
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

	// Insère des données de test
	_, err = db.Exec(`INSERT INTO TestTable (Id, Email, Password) VALUES 
		("1", "superemail@gmail.com", "JAimeCoder1235"), 
		("2", "superemail@gmail.com", "JAimeCoder1234")`)
	if err != nil {
		t.Fatalf("Erreur lors de l'insertion des données : %v", err)
	}

	// Arguments pour la sélection (exemple avec Email et Password)
	args := map[string]any{
		"Email":    "superemail@gmail.com",
		"Password": "JAimeCoder1234",
	}

	// Appel de la fonction SelectFromDb
	result, err := SelectFromDb("TestTable", db, args)
	if err != nil {
		t.Fatalf("Erreur lors de l'exécution de SelectFromDb : %v", err)
	}

	// Vérifie qu'on a obtenu une seule ligne
	if len(result) != 1 {
		t.Fatalf("Nombre de lignes attendu : 1, obtenu : %d", len(result))
	}

	// Vérifie les valeurs des colonnes
	row := result[0]
	id := *(row[0].(*string)) // Convertir les valeurs récupérées en *string
	email := *(row[1].(*string))
	password := *(row[2].(*string))

	// Vérifie que les données sont correctes
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
