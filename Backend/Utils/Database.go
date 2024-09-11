package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDb() (*sql.DB, error) {
	return sql.Open("sqlite3", "./Database/Database.sqlite")
}

func CreateDb() {
	db, err := OpenDb()
	if err != nil {
		fmt.Println("Error:", err)
	}

	defer db.Close()

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
	_, err = db.Exec(r)
	if err != nil {
		fmt.Println("Create Error", err)
	}

}

func InsertIntoDb(tabelName string, Args ...any) error {
	db, err := OpenDb()
	if err != nil {
		return err
	}
	defer db.Close()

	var stringMAP string
	for i, j := range Args {
		if i < len(Args)-1 {
			stringMAP += fmt.Sprintf("\"%s\", ", j)
		} else {
			stringMAP += fmt.Sprintf("\"%s\"", j)
		}
	}

	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s VALUES(%s)", tabelName, stringMAP))
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func SelectFromDb(tabelName string, Args map[string]any) ([][]interface{}, error) {
	column, rows, err := prepareStmt(tabelName, Args)
	if err != nil {
		return nil, err
	}

	var result [][]interface{}
	for rows.Next() {
		row := make([]interface{}, len(column))
		for i := 0; i < len(column); i++ {
			row[i] = new(string)
		}

		if err := rows.Scan(row...); err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}

func prepareStmt(tabelName string, Args map[string]any) ([]string, *sql.Rows, error) {
	db, err := OpenDb()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	var stringMAP string
	for i, j := range Args {
		stringMAP += fmt.Sprintf("%s=%s ", i, j)
	}

	stmt, err := db.Prepare(fmt.Sprintf("SELECT * from %s where %s ", tabelName, stringMAP))
	if err != nil {
		return nil, nil, err
	}

	rows, err := stmt.Query(stmt)
	if err != nil {
		return nil, nil, err
	}

	column, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	return column, rows, nil
}
