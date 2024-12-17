package utils

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	model "social-network/Model"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
This function takes 2 arguments:
  - a string JWT, which represents the JSON Web Token to be decrypted.
  - a pointer to an sql.DB object, representing the database connection.

The purpose of this function is to decrypt the provided JWT and verify its integrity.

The function returns 2 values:
  - a string containing the decrypted ID if the operation is successful.
  - an error if the JWT is invalid, decryption fails, or if the ID does not exist in the database.
*/
func DecryptJWT(JWT string, db *sql.DB) (string, error) {
	// We split the JWT into its three parts: header, payload, and signature
	splitSessionId := strings.Split(JWT, ".")
	if len(splitSessionId) != 3 {
		// Return an error if the JWT does not have exactly three parts
		return "", errors.New("invalide JWT")
	}

	// We compare the hashed password with the provided secret key for integrity verification
	if err := bcrypt.CompareHashAndPassword([]byte(splitSessionId[2]), []byte(model.SecretKey)); err != nil {
		// Return any error encountered during the hash comparison
		return "", err
	}

	// We decode the base64-encoded payload part of the JWT
	decryptId, err := base64.StdEncoding.DecodeString(splitSessionId[1])
	if err != nil {
		// Return any error encountered during the base64 decoding
		return "", err
	}

	// We check if the decrypted ID exists in the "Auth" table in the database
	err = IfExistsInDB("Auth", db, map[string]any{"Id": string(decryptId)})
	// Return the decrypted ID or any error encountered
	return string(decryptId), err
}

/*
This function takes 3 arguments:
  - a string table, which specifies the name of the database table to check.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the arguments for the query, representing the conditions for selecting records.

The purpose of this function is to check if exactly one record exists in the specified table that matches the provided conditions.

The function returns 1 value:
  - an error if the selection fails or if there is not exactly one match.
*/
func IfExistsInDB(table string, db *sql.DB, args map[string]any) error {
	// We call SelectFromDb to retrieve data from the specified table based on the provided conditions
	authData, err := model.SelectFromDb(table, db, args)
	if err != nil {
		// Return any error encountered during data retrieval
		return err
	}

	// We check if the length of the retrieved data is not equal to 1
	if len(authData) == 0 {
		// Return an error if there is no match or multiple matches found
		return errors.New("there is no match")
	}

	// Return nil if exactly one match is found
	return nil
}

/*
This function takes 3 arguments:
  - a string table, which specifies the name of the database table to check.
  - a pointer to an sql.DB object, representing the database connection.
  - a map[string]any containing the arguments for the query, representing the conditions for selecting records.

The purpose of this function is to check if no records exist in the specified table that match the provided conditions.

The function returns 1 value:
  - an error if the selection fails or if there is at least one match.
*/
func IfNotExistsInDB(table string, db *sql.DB, args map[string]any) error {
	// We call SelectFromDb to retrieve data from the specified table based on the provided conditions
	authData, err := model.SelectFromDb(table, db, args)
	if err != nil {
		// Return any error encountered during data retrieval
		return err
	}

	// We check if the length of the retrieved data is not equal to 0
	if len(authData) != 0 {
		// Return an error if there is at least one match found
		return errors.New("there is a match")
	}

	// Return nil if no matches are found
	return nil
}
