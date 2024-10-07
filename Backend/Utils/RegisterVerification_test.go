package utils

// import (
// 	model "social-network/Model"
// 	"testing"

// 	"golang.org/x/crypto/bcrypt"
// )

// func TestRegisterVerification(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		data       model.Register
// 		shouldFail bool
// 	}{
// 		{
// 			name: "Valid registration",
// 			data: model.Register{
// 				Auth: model.Auth{
// 					Email:           "unemail@gmail.com",
// 					Password:        "zXYVhVxp9zxP8qa$",
// 					ConfirmPassword: "zXYVhVxp9zxP8qa$",
// 				},
// 				FirstName: "Jean",
// 				LastName:  "Dujardin",
// 				BirthDate: "1998-01-03",
// 			},
// 			shouldFail: false,
// 		},
// 		{
// 			name: "Password and confirm password do not match",
// 			data: model.Register{
// 				Auth: model.Auth{
// 					Email:           "unemail@gmail.com",
// 					Password:        "zXYVhVxp9zxP8qa$",
// 					ConfirmPassword: "differentPassword",
// 				},
// 				FirstName: "Jean",
// 				LastName:  "Dujardin",
// 				BirthDate: "1998-01-03",
// 			},
// 			shouldFail: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := RegisterVerification(tt.data)
// 			if (err != nil) != tt.shouldFail {
// 				t.Fatalf("Test '%s' échoué : attendu erreur: %v, obtenu: %v", tt.name, tt.shouldFail, err != nil)
// 				return
// 			}
// 		})
// 	}
// }

// func TestIsValidPassword(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		data       string
// 		shouldFail bool
// 	}{
// 		{
// 			name:       "Short Password",
// 			data:       "Ey$21",
// 			shouldFail: true,
// 		},

// 		{
// 			name:       "Contains Uppercase, No Special Char",
// 			data:       "IFBSOSNHFBJ",
// 			shouldFail: true,
// 		},
// 		{
// 			name:       "Contains Number, No Special Char ",
// 			data:       "IDBF2847492",
// 			shouldFail: true,
// 		},
// 		{
// 			name:       "Password Valide",
// 			data:       "zXYVhVxp9@P8qa",
// 			shouldFail: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			valid := IsValidPassword(tt.data)
// 			if valid == tt.shouldFail {
// 				t.Fatalf("Test '%s' échoué : attendu erreur: %v, obtenu: %v", tt.name, tt.shouldFail, !valid)
// 				return
// 			}
// 		})
// 	}
// }

// func TestCreateUuidAndCrypt(t *testing.T) {
// 	// Créer un modèle Register de test
// 	register := &model.Register{
// 		Auth: model.Auth{
// 			Email:    "unemail@gmail.com",
// 			Password: "MonMotDePasse123!",
// 		},
// 		FirstName: "Jean",
// 		LastName:  "Dujardin",
// 		BirthDate: "1990-01-01",
// 	}

// 	// Appeler la fonction CreateUuidAndCrypt
// 	err := CreateUuidAndCrypt(register)
// 	if err != nil {
// 		t.Fatalf("Erreur lors de l'exécution de CreateUuidAndCrypt: %v", err)
// 		return
// 	}

// 	// Vérifier que le mot de passe a bien été crypté
// 	if err := bcrypt.CompareHashAndPassword([]byte(register.Auth.Password), []byte("MonMotDePasse123!")); err != nil {
// 		t.Errorf("Le mot de passe crypté ne correspond pas au mot de passe original")
// 		return
// 	}

// 	// Vérifier que l'UUID a bien été généré
// 	if register.Auth.Id == "" {
// 		t.Errorf("L'UUID n'a pas été généré correctement")
// 		return
// 	}
// }
