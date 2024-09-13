package utils

import (
	"database/sql"
	"fmt"
	"strings"

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
This function takes a minimum of 2 arguments:
  - a string who is the name of the table
  - the values you want to add in this table (all types are accepted)

The objective of this function is to insert values in a table of the db.

The function gonna return:
  - an error
*/
func InsertIntoDb(tabelName string, db *sql.DB, Args ...any) error {
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
  - an map where each key corresponds to a column in the db
  - an error
*/
func SelectFromDb(tabelName string, db *sql.DB, Args map[string]any) ([]map[string]any, error) {
	// We prepare and execute the select request
	column, rows, err := PrepareStmt(tabelName, db, Args)
	if err != nil {
		return nil, err
	}

	// We loop on the result to stock them into the []map[string]any
	var result []map[string]any
	for rows.Next() {
		// We initialize the variable who gonna contain the current result row
		row := make(map[string]any)

		values := make([]interface{}, len(column))
		for i := 0; i < len(column); i++ {
			values[i] = new(string)
		}

		// We fill the variable with the values of the row
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		// We add the values row by row in the current map
		for i, v := range column {
			row[v] = values[i]
		}

		// We add the values of the current row into the []map[string]any
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
func PrepareStmt(tabelName string, db *sql.DB, Args map[string]any) ([]string, *sql.Rows, error) {

	var whereClauses []string
	var params []any

	// Building the WHERE clause with parameters
	for column, value := range Args {
		// Use "?" for parameters
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", column))
		// Add the corresponding values
		params = append(params, value)
	}

	// Joins the WHERE clauses with "AND" to form the condition
	whereString := ""
	if len(whereClauses) > 0 {
		whereString = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Build the SQL query with WHERE clauses
	query := fmt.Sprintf("SELECT * FROM %s %s", tabelName, whereString)

	// Prepare the SQL query
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, nil, err
	}
	defer stmt.Close()

	// Executes the query passing the parameters
	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, nil, err
	}

	// Retrieves the columns of the result
	column, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	// We return the column and the result of the SQL request
	return column, rows, nil
}
