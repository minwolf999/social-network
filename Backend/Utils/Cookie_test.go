package utils

import (
	"net/http/httptest"
	"testing"
)

// Test function to simulate the http.ResponseWriter and check cookie setting
func TestSetCookie(t *testing.T) {
	// Use httptest to simulate http.ResponseWriter
	rr := httptest.NewRecorder()

	// Call the setCookie function
	SetCookie(rr, "testValue")

	// Check if the cookie is set correctly in the response header
	cookies := rr.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("Expected a cookie to be set")
	}

	// Validate cookie name and value
	cookie := cookies[0]
	if cookie.Name != "sessionId" {
		t.Errorf("Expected cookie name 'testCookie', got %s", cookie.Name)
	}
	if cookie.Value != "testValue" {
		t.Errorf("Expected cookie value 'testValue', got %s", cookie.Value)
	}
}
