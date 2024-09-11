package utils

import (
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
