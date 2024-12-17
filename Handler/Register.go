package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	model "social-network/Model"
	utils "social-network/Utils"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nw := model.ResponseWriter{
			ResponseWriter: w,
		}

		// We read the request body and unmarshal it into a structure
		var register model.Register
		if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
			nw.Error("Invalid request body")
			log.Printf("[%s] [Register] Invalid request body: %v", r.RemoteAddr, err)
			return
		}

		// We look if all is good in the datas send in the body of the request
		if err := RegisterVerification(register); err != nil {
			nw.Error(err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
			return
		}

		// We generate an UUID and crypt the password
		if err := CreateUuidAndCrypt(&register); err != nil {
			nw.Error(err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
			return
		}

		register.Id = register.Auth.Id
		register.Email = register.Auth.Email

		if len(register.ProfilePicture) > 400000 {
			nw.Error("To big image")
			log.Printf("[%s] [Register] To big image", r.RemoteAddr)
			return
		}

		// We get the row in the db where the email is equal to the email send
		if err := utils.IfExistsInDB("Auth", db, map[string]any{"Email": register.Auth.Email}); err != nil && err.Error() != "there is no match" {
			nw.Error("Email is already used : " + err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, "Email is already used")
			return
		}

		// We insert in the table Auth of the db the id, email and password of the people trying to register
		if err := register.Auth.InsertIntoDb(db); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
			return
		}

		// We insert in the table UserInfo of the db the rest of the values
		if err := register.InsertIntoDb(db); err != nil {
			nw.Error("Internal Error: There is a probleme during the push in the DB: " + err.Error())
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
			return
		}

		// We send a success response to the request
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"Success":   true,
			"Message":   "Login successfully",
			"sessionId": utils.GenerateJWT(register.Auth.Id),
		})
		if err != nil {
			log.Printf("[%s] [Register] %s", r.RemoteAddr, err.Error())
		}
	}
}

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
