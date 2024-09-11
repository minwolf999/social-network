package utils

import (
	"fmt"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

/*
This function takes 2 arguments:
  - the type of the db (mysql, sqlite, ...)
  - the path to find the db from the file main.go

The objective of this function is to open the connection to the db.

The function gonna return:
  - a connection to the db
  - an error
*/
func OpenDb(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

/*
This function takes no argument

The objective of this function is to create the table in the db.

The function returns no value
*/
func CreateDb() {
	// We open the db and close at the end of this fonction
	db, err := OpenDb("sqlite3", "./Database/Database.sqlite")
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer db.Close()

	// We initialize the SQL query by writing it into a string
	r := `
	CREATE TABLE IF NOT EXISTS Auth (
		Id VARCHAR(36) NOT NULL PRIMARY KEY,
		Email VARCHAR(100) NOT NULL,
		Password VARCHAR(50) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS UserInfo (
		Id VARCHAR(36) NOT NULL REFERENCES "Auth"("Id"),
		Email VARCHAR(100) NOT NULL REFERENCES "Auth"("Email"),
		FirstName VARCHAR(50) NOT NULL, 
		LastName VARCHAR(50) NOT NULL,
		Birth VARCHAR(20) NOT NULL,
		Avatar VARCHAR(100),
		Username VARCHAR(50),
		AboutMe VARCHAR(280)  
	);
	`

	// We execute the SQL request
	_, err = db.Exec(r)
	if err != nil {
		fmt.Println("Create Error", err)
	}

}

/*
This function takes a minimum of 2 arguments:
  - a string who is the name of the table
  - the values you want to add in this table (all types are accepted)

The objective of this function is to insert values in a table of the db.

The function gonna return:
  - an error
*/
func InsertIntoDb(tabelName string, Args ...any) error {
	// We open the db and close at the end of this fonction
	db, err := OpenDb("sqlite3", "./Database/Database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// We format the values to write them into a string
	var stringMAP string
	for i, j := range Args {
		if i < len(Args)-1 {
			stringMAP += fmt.Sprintf("\"%s\", ", j)
		} else {
			stringMAP += fmt.Sprintf("\"%s\"", j)
		}
	}

	// We prepare the SQL query to avoid SQL injections
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s VALUES(%s)", tabelName, stringMAP))
	if err != nil {
		return err
	}

	// We execute the SQL request
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

/*
This function takes 2 arguments:
  - a string who is the name of the table
  - a map with containing the wanted values in each row

The objective of this function is to get values in a table of the db.

The function gonna return:
  - an [][]interface where each []interface corresponds to a row in the db
  - an error
*/
func SelectFromDb(tabelName string, Args map[string]any) ([][]interface{}, error) {
	// We prepare and execute the select request
	column, rows, err := prepareStmt(tabelName, Args)
	if err != nil {
		return nil, err
	}

	// We loop on the result to stock them into the [][]interface
	var result [][]interface{}
	for rows.Next() {
		// We initialize the variable who gonna contain the current result row
		row := make([]interface{}, column)
		for i := 0; i < column; i++ {
			row[i] = new(string)
		}

		// We fill the variable with the values of the row
		if err := rows.Scan(row...); err != nil {
			return nil, err
		}

		// We add the values of the current row into the [][]structure
		result = append(result, row)
	}

	return result, nil
}

/*
This function takes 2 arguments:
  - a string who is the name of the table
  - a map with containing the wanted values in each row

The objective of this function is to format, prepare and execute the SQL request.

The function gonna return:
  - an int who corresponds to the row quantity of the table
  - a *sql.Rows who contain the result of the SQL request
*/
func prepareStmt(tabelName string, Args map[string]any) (int, *sql.Rows, error) {
	// We open the db and close at the end of this fonction
	db, err := OpenDb("sqlite3", "./Database/Database.sqlite")
	if err != nil {
		return 0, nil, err
	}
	defer db.Close()

	// We format the values to write them into a string
	var stringMAP string
	for i, j := range Args {
		stringMAP += fmt.Sprintf("%s=%s ", i, j)
	}

	// We prepare the SQL query to avoid SQL injections
	stmt, err := db.Prepare(fmt.Sprintf("SELECT * from %s where %s ", tabelName, stringMAP))
	if err != nil {
		return 0, nil, err
	}

	// We execute the SQL request and stock the result of the request inside a *sql.Rows
	rows, err := stmt.Query(stmt)
	if err != nil {
		return 0, nil, err
	}

	// We stock the []string who contains the name of all the rows
	column, err := rows.Columns()
	if err != nil {
		return 0, nil, err
	}

	// We return the row quantity and the result of the SQL request
	return len(column), rows, nil
}
