package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	// Create a ServeMux and configure routes without the database
	mux := http.NewServeMux()
	RoutesForTest(mux)

	// Test board with route cases to test
	tests := []struct {
		method string
		route  string
		status int
	}{
		{"GET", "/", http.StatusOK},
		{"POST", "/login", http.StatusOK},
		{"POST", "/register", http.StatusOK},
		{"GET", "/home", http.StatusOK},
	}

	for _, tt := range tests {
		// Create a new query for each route to be tested
		req, err := http.NewRequest(tt.method, tt.route, nil)
		if err != nil {
			t.Fatalf("Erreur lors de la création de la requête : %v", err)
			return
		}

		// Simulate a request with httptest
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// Check that the status code is correct
		if rr.Code != tt.status {
			t.Errorf("Pour la route %s %s, code de statut attendu: %v, obtenu: %v", tt.method, tt.route, tt.status, rr.Code)
			return
		}
	}
}

