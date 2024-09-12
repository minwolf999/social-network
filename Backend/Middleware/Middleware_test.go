package middleware

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	model "social-network/Model"
	utils "social-network/Utils"
	"testing"
)

func TestRegisterMiddleware(t *testing.T) {
	// Crée un mock de base de données (ou une vraie connexion en mémoire)
	db, err := utils.OpenDb("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Erreur lors de la création de la base de données en mémoire : %v", err)
	}
	defer db.Close()

	// Crée une structure Register de test
	register := model.Register{
		Auth: model.Auth{
			Email:           "unemail@gmail.com",
			Password:        "MonMotDePasse123!",
			ConfirmPassword: "MonMotDePasse123!",
		},
		FirstName: "Jean",
		LastName:  "Dujardin",
		BirthDate: "1990-01-01",
	}

	// Sérialise la structure Register en JSON
	body, err := json.Marshal(register)
	if err != nil {
		t.Fatalf("Erreur lors de la sérialisation du corps de la requête : %v", err)
	}

	// Crée une requête HTTP POST simulée
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Crée un ResponseRecorder pour capturer la réponse
	rr := httptest.NewRecorder()

	// Crée une fonction next factice pour passer au prochain handler après le middleware
	next := func(w http.ResponseWriter, r *http.Request, db *sql.DB) {
		// Si le middleware est exécuté correctement, cela signifie qu'on arrive à cette fonction
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Middleware passed and next handler called"))
	}

	// Appel du middleware
	handler := RegisterMiddleware(next, db)
	handler.ServeHTTP(rr, req)

	// Vérifie le code de statut de la réponse
	if rr.Code != http.StatusOK {
		t.Errorf("Code de statut incorrect : attendu %v, obtenu %v", http.StatusOK, rr.Code)
	}

	// Vérifie le corps de la réponse
	expectedBody := "Middleware passed and next handler called"
	if rr.Body.String() != expectedBody {
		t.Errorf("Corps de réponse incorrect : attendu %v, obtenu %v", expectedBody, rr.Body.String())
	}
}
