package registermiddlewaresubfunction

import (
	"errors"
	"fmt"
	"net/http"
	model "social-network/Model"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RegisterVerification(register model.Register, w http.ResponseWriter) error {
	nw := model.ResponseWriter{
		ResponseWriter: w,
	}

	if register.Auth.Password != register.Auth.ConfirmPassword {
		nw.Error("Password and password confirmation do not match")
		return errors.New("password and password confirmation do not match")
	}

	if !IsValidPassword(register.Auth.Password) {
		nw.Error("Incorrect password ! The password must contain 8 characters, 1 uppercase letter, 1 special character, 1 number")
		return errors.New("incorrect password ! the password must contain 8 characters, 1 uppercase letter, 1 special character, 1 number")
	}

	if register.Auth.Email == "" || register.Auth.Password == "" || register.FirstName == "" || register.LastName == "" || register.BirthDate == "" {
		nw.Error("There is an empty field")
		return errors.New("there is an empty field")
	}

	return nil
}

func IsValidPassword(password string) bool {
	var isLongEnought bool = false
	var containUpper bool = false
	var containSpeChar bool = false
	var containNumber bool = false
	if len(password) >= 8 {
		isLongEnought = true
	}
	for _, r := range password {
		if r >= 'A' && r <= 'Z' {
			containUpper = true
		} else if r >= '0' && r <= '9' {
			containNumber = true
		} else if r < 'a' || r > 'z' {
			containSpeChar = true
		}
	}
	if isLongEnought && containNumber && containSpeChar && containUpper {
		return true
	}
	return false
}

func CreateUuidAndCrypt(register *model.Register, w http.ResponseWriter) error {
	nw := model.ResponseWriter{
		ResponseWriter: w,
	}

	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Auth.Password), 12)
	if err != nil {
		fmt.Println(err)

		nw.Error("Internal Error: There is a probleme with bcrypt")
		return errors.New("there is a probleme with bcrypt")
	}
	register.Auth.Password = string(cryptedPassword)

	uuid, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
		nw.Error("Internal Error: There is a probleme with the generation of the uuid")
		return errors.New("there is a probleme with the generation of the uuid")
	}
	register.Auth.Id = uuid.String()

	return nil
}
