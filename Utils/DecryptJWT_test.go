package utils

import (
	"encoding/base64"
	model "social-network/Model"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerateJWT(t *testing.T) {
	value := "Test"
	JWT := GenerateJWT(value)

	splitJWT := strings.Split(JWT, ".")
	if len(splitJWT) != 3 {
		t.Errorf("The 3 part of the JWT are not here")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(splitJWT[2]), []byte(model.SecretKey)); err != nil {
		t.Errorf("Invalid secret key: %v", err)
		return
	}

	decrypt, err := base64.StdEncoding.DecodeString(splitJWT[1])
	if err != nil {
		t.Errorf("Invalid original base format")
		return
	}

	if value != string(decrypt) {
		t.Errorf("Invalid value in JWT")
		return
	}
}

func TestDecryptJWT(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	db.Exec(`
		CREATE TABLE IF NOT EXISTS Auth (
			Id VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY,
			Email VARCHAR(100) NOT NULL UNIQUE,
			Password VARCHAR(50) NOT NULL
		);
	`)

	var value struct {
		Id       string
		Email    string
		Password string
	}

	value.Id = "Id"
	value.Email = "Email"
	value.Password = "Password"

	_, err = db.Exec("INSERT INTO Auth VALUES (?,?,?)", value.Id, value.Email, value.Password)
	if err != nil {
		t.Fatalf("Error the insert in the db : %s", err)
		return
	}

	JWT := GenerateJWT(value.Id)

	decrypt, err := DecryptJWT(JWT, db)
	if err != nil {
		t.Fatalf("Error during the decrypt : %s", err)
		return
	}

	if decrypt != value.Id {
		t.Fatalf("Invalid value")
		return
	}
}

func TestIfExistsInDB(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	db.Exec(`
		CREATE TABLE IF NOT EXISTS Auth (
			Id VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY,
			Email VARCHAR(100) NOT NULL UNIQUE,
			Password VARCHAR(50) NOT NULL
		);
	`)

	var value struct {
		Id       string
		Email    string
		Password string
	}

	value.Id = "Id"
	value.Email = "Email"
	value.Password = "Password"

	_, err = db.Exec("INSERT INTO Auth VALUES (?,?,?)", value.Id, value.Email, value.Password)
	if err != nil {
		t.Fatalf("Error the insert in the db : %s", err)
		return
	}

	if err = IfExistsInDB("Auth", db, map[string]any{"Id": value.Id}); err != nil {
		t.Fatalf("Error during the function : %s", err)
		return
	}
}

func TestIfNotExistsInDB(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := model.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
		return
	}
	defer db.Close()

	db.Exec(`
		CREATE TABLE IF NOT EXISTS Auth (
			Id VARCHAR(36) NOT NULL UNIQUE PRIMARY KEY,
			Email VARCHAR(100) NOT NULL UNIQUE,
			Password VARCHAR(50) NOT NULL
		);
	`)

	var value struct {
		Id       string
		Email    string
		Password string
	}

	value.Id = "Id"

	if err = IfNotExistsInDB("Auth", db, map[string]any{"Id": value.Id}); err != nil {
		t.Fatalf("Error during the function : %s", err)
		return
	}
}
