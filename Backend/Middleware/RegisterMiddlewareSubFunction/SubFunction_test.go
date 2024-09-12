package registermiddlewaresubfunction

import (
	model "social-network/Model"
	"testing"
)

func TestRegisterVerification(t *testing.T) {
	tests := []struct {
		name       string
		data       model.Register
		shouldFail bool
	}{
		{
			name: "Valid registration",
			data: model.Register{
				Auth: model.Auth{
					Email:           "unemail@gmail.com",
					Password:        "zXYVhVxp9zxP8qa$",
					ConfirmPassword: "zXYVhVxp9zxP8qa$",
				},
				FirstName: "Jean",
				LastName:  "Dujardin",
				BirthDate: "1998-01-03",
			},
			shouldFail: false,
		},
		{
			name: "Password and confirm password do not match",
			data: model.Register{
				Auth: model.Auth{
					Email:           "unemail@gmail.com",
					Password:        "zXYVhVxp9zxP8qa$",
					ConfirmPassword: "differentPassword",
				},
				FirstName: "Jean",
				LastName:  "Dujardin",
				BirthDate: "1998-01-03",
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterVerification(tt.data)
			if (err != nil) != tt.shouldFail {
				t.Fatalf("Test '%s' échoué : attendu erreur: %v, obtenu: %v", tt.name, tt.shouldFail, err != nil)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {

}
