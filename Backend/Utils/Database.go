package utils

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDb(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

func CreateDb(db *sql.DB) {
	var err error
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

func InsertIntoDb(tabelName string, db *sql.DB, Args ...any) error {
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

func SelectFromDb(tabelName string, db *sql.DB, Args map[string]any) ([][]interface{}, error) {
	column, rows, err := PrepareStmt(tabelName, db, Args)
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

func PrepareStmt(tabelName string, db *sql.DB, Args map[string]any) ([]string, *sql.Rows, error) {
	var whereClauses []string
	var params []any

	// Construction de la clause WHERE avec des paramètres
	for column, value := range Args {
		// Utilise "?" pour les paramètres
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", column))
		// Ajoute les valeurs correspondantes
		params = append(params, value)
	}

	// Joint les clauses WHERE avec "AND" pour former la condition
	whereString := ""
	if len(whereClauses) > 0 {
		whereString = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Construit la requête SQL avec les clauses WHERE
	query := fmt.Sprintf("SELECT * FROM %s %s", tabelName, whereString)

	// Prépare la requête SQL
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, nil, err
	}
	defer stmt.Close()

	// Exécute la requête en passant les paramètres
	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, nil, err
	}

	// Récupère les colonnes du résultat
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	return columns, rows, nil
}

