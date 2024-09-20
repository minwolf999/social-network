package utils

import (
	"encoding/base64"
	"strings"
	"testing"

	model "social-network/Model"

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
	value := "test"
	JWT := GenerateJWT(value)

	decrypt, err := DecryptJWT(JWT)
	if err != nil {
		t.Fatalf("Error during the decrypt : %s", err)
		return
	}

	if decrypt != value {
		t.Fatalf("Invalid value")
		return
	}
}
