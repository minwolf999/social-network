package utils

import (
	"errors"
	
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	
	model "social-network/Model"
)

/*
This function takes 1 argument:
  - an Register variable who contain the datas send in the body of the request

The purpose of this function is to Verificate the content of the request make to the Register function.

The function return an error
*/
func RegisterVerification(register model.Register) error {
	// We check if the password match to the confimation of the password
	if register.Auth.Password != register.Auth.ConfirmPassword {
		return errors.New("password and password confirmation do not match")
	}

	// We check if the password is secure enough
	if !IsValidPassword(register.Auth.Password) {
		return errors.New("incorrect password ! the password must contain 8 characters, 1 uppercase letter, 1 special character, 1 number")
	}

	// We check if all the needed information are here
	if register.Auth.Email == "" || register.Auth.Password == "" || register.FirstName == "" || register.LastName == "" || register.BirthDate == "" {
		return errors.New("there is an empty field")
	}

	return nil
}

/*
This function takes 1 argument:
  - a string who contain the password

# The purpose of this function is to look if the password is secure enough

The function return an boolean
*/
func IsValidPassword(password string) bool {
	// We start by initializing the check variable
	isLongEnought := false
	containUpper := false
	containSpeChar := false
	containNumber := false

	// We look if the password contain at least 8 characters
	if len(password) >= 8 {
		isLongEnought = true
	}

	// We look if there is at least 1 lowercase, uppercase, number, special character
	for _, r := range password {
		if r >= 'A' && r <= 'Z' {
			containUpper = true
		} else if r >= '0' && r <= '9' {
			containNumber = true
		} else if r < 'a' || r > 'z' {
			containSpeChar = true
		}
	}

	// If all goes well we return true otherwise false
	if isLongEnought && containNumber && containSpeChar && containUpper {
		return true
	}
	return false
}

/*
This function takes 1 argument:
  - a *Register variable who contain the datas send in the body of the request

# The purpose of this function is to create a new UUID and crypt the password who is in the structure Register

The function return an error
*/
func CreateUuidAndCrypt(register *model.Register) error {
	// We crypt the password and replace the previous password by the crypted version
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Auth.Password), 12)
	if err != nil {
		return errors.New("there is a probleme with bcrypt")
	}
	register.Auth.Password = string(cryptedPassword)

	// We generate a new UUID and store it into the structure
	uuid, err := uuid.NewV7()
	if err != nil {
		return errors.New("there is a probleme with the generation of the uuid")
	}
	register.Auth.Id = uuid.String()

	return nil
}
