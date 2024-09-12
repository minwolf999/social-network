package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	// Crée un ServeMux et configure les routes sans la base de données
	mux := http.NewServeMux()
	RoutesForTest(mux)

	// Tableau de test avec des cas de routes à tester
	tests := []struct {
		method string
		route  string
		status int // Statut HTTP attendu
	}{
		{"GET", "/", http.StatusOK},
		{"POST", "/login", http.StatusOK},
		{"POST", "/register", http.StatusOK},
		{"GET", "/home", http.StatusOK},
	}

	for _, tt := range tests {
		// Crée une nouvelle requête pour chaque route à tester
		req, err := http.NewRequest(tt.method, tt.route, nil)
		if err != nil {
			t.Fatalf("Erreur lors de la création de la requête : %v", err)
		}

		// Simule une requête avec httptest
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// Vérifie que le code de statut est correct
		if rr.Code != tt.status {
			t.Errorf("Pour la route %s %s, code de statut attendu: %v, obtenu: %v", tt.method, tt.route, tt.status, rr.Code)
		}
	}
}

