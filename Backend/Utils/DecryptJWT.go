package utils

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	model "social-network/Model"
)

/*
This function takes 1 argument:
  - a string who contain the value to set in the JWT

The purpose of this function is to create a JWT with a value crypted inside him.

The function return 1 value:
  - a string who is the JWT
*/
func GenerateJWT(str string) string {
	// We convert the header object in base 64 for the first part
	header := base64.StdEncoding.EncodeToString([]byte(`{
		"typ": "JWT"
	}`))

	// We convert the value in base 64 for the second part
	content := base64.StdEncoding.EncodeToString([]byte(str))

retry:
	// We hash the key for the last part od the JWT
	key, err := bcrypt.GenerateFromPassword([]byte(model.SecretKey), 12)
	if err != nil {
		fmt.Println(err)
	}
	if strings.Contains(string(key), ".") {
		goto retry
	}

	// We assemble the 3 part with a . between each part
	return header + "." + content + "." + string(key)
}

/*
This function takes 1 argument:
  - a string who contain the JWT

The purpose of this function is to decrypt the value defined in the JWT and verify that it is correctly formatted.

The function return 2 values:
  - a string who contain the value set in the JWT
  - an error
*/
func DecryptJWT(JWT string, db *sql.DB) (string, error) {
	// We split the 3 part of the JWT
	splitSessionId := strings.Split(JWT, ".")
	if len(splitSessionId) != 3 {
		return "", errors.New("invalide JWT")
	}

	// We check if the secret key is good
	if err := bcrypt.CompareHashAndPassword([]byte(splitSessionId[2]), []byte(model.SecretKey)); err != nil {
		return "", err
	}

	// We decode the sessionId
	decryptId, err := base64.StdEncoding.DecodeString(splitSessionId[1])
	if err != nil {

	}

	IfExistsInDB("Auth", db, map[string]any{"Id": string(decryptId)})
	return string(decryptId), err
}

func IfExistsInDB(table string, db *sql.DB, args map[string]any) error {
	authData, err := SelectFromDb(table, db, args)
	if err != nil {
		return err
	}

	if len(authData) != 1 {
		return errors.New("there is no match")
	}

	return nil
}
