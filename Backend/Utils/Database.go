package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	model "social-network/Model"
	"strings"

	"github.com/gofrs/uuid"
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
This function takes 1 argument:
  - a connection with th db

# The purpose of this function is to fill the db with 100 value in each table

The function return an error
*/
func LoadData(db *sql.DB) error {
	// We set a list of first and last name
	firstNameList := []string{"Alice", "Bob", "Charlie", "David", "Eva", "Frank", "Grace", "Hannah", "Ivy", "Jack", "Karen", "Leo", "Mia", "Noah", "Olivia", "Paul", "Quinn", "Ruby", "Sam", "Tina", "Uma", "Victor", "Wendy", "Xander", "Yara", "Zane", "Adrian", "Bella", "Carl", "Diana", "Ethan", "Fiona", "George", "Helen", "Isaac", "Julia", "Kevin", "Lara", "Michael", "Nina", "Oscar", "Penny", "Quentin", "Rachel", "Steve", "Tara", "Uriel", "Violet", "Walter", "Xenia", "Yves", "Zelda", "Arthur", "Bianca", "Colin", "Derek", "Emma", "Felix", "Gina", "Harry", "Iris", "James", "Kara", "Louis", "Maria", "Nathan", "Owen", "Pam", "Ron", "Sophie", "Tom", "Ursula", "Vincent", "Will", "Ximena", "Yvonne", "Zach", "Angela", "Bruno", "Claire", "Damien", "Elise", "Freddy", "Gloria", "Henry", "Isabelle", "Julien", "Kurt", "Liam", "Nadine", "Olga", "Peter", "Quincy", "Rosie", "Simon", "Tracy", "Ulrich", "Victoria", "Wayne", "Xia", "Yasmine", "Zeke"}
	lastNameList := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker", "Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores", "Green", "Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell", "Carter", "Roberts", "Gomez", "Phillips", "Evans", "Turner", "Diaz", "Parker", "Cruz", "Edwards", "Collins", "Reyes", "Stewart", "Morris", "Morales", "Murphy", "Cook", "Rogers", "Gutierrez", "Ortiz", "Morgan", "Cooper", "Peterson", "Bailey", "Reed", "Kelly", "Howard", "Ramos", "Kim", "Cox", "Ward", "Richardson", "Watson", "Brooks", "Chavez", "Wood", "James", "Bennett", "Gray", "Mendoza", "Ruiz", "Hughes", "Price", "Alvarez", "Castillo", "Sanders", "Patel", "Myers", "Long", "Ross", "Foster", "Jimenez", "Powell", "Jenkins", "Perry", "Russell", "Sullivan", "Bell", "Coleman", "Butler", "Henderson", "Barnes"}

	subjects := []string{"Le développeur", "L'algorithme", "Go", "Le projet open-source", "La communauté", "L'interface utilisateur", "Le serveur", "La technologie", "Le programmeur", "Le code"}
	verbs := []string{"améliore", "optimise", "développe", "test", "documente", "déploie", "corrige", "met à jour", "construit", "gère"}
	objects := []string{"la performance", "l'API", "le backend", "le frontend", "l'application web", "le système de fichiers", "la sécurité", "la base de données", "les requêtes HTTP", "la fonctionnalité asynchrone"}
	adverbs := []string{"rapidement", "efficacement", "avec soin", "de manière innovante", "avec passion", "de façon concurrente", "sans erreur", "avec succès", "sans latence", "de manière fluide"}

	// We loop at least 100 times
	for i := 0; i < 100; i++ {
		// We create a variable of type Register and set the value inside (we get a random first and last name in the 2 lists)
		var user model.Register
		user.FirstName = firstNameList[rand.Intn(len(firstNameList))]
		user.LastName = lastNameList[rand.Intn(len(lastNameList))]
		user.Auth.Password = "Azerty&1234"
		user.Auth.Email = fmt.Sprintf("%s.%s@gmail.com", user.FirstName, user.LastName)

		// We generate 3 random number for the birthDate
		day := rand.Intn(31-0) + 0
		mounth := rand.Intn(12-0) + 0
		year := rand.Intn(2024-1980) + 1980
		user.BirthDate = fmt.Sprintf("%d-%d-%d", year, mounth, day)

		// We create an UUID and hash the password
		if err := CreateUuidAndCrypt(&user); err != nil {
			return err
		}

		// We insert the values in the tables
		if err := InsertIntoDb("Auth", db, user.Auth.Id, user.Auth.Email, user.Auth.Password); err != nil {
			i--
			continue
		}
		if err := InsertIntoDb("UserInfo", db, user.Auth.Id, user.Auth.Email, user.FirstName, user.LastName, user.BirthDate, user.ProfilePicture, user.Username, user.AboutMe); err != nil {
			i--
			continue
		}

		var post model.Post
		post.AuthorId = user.Auth.Id
		uuid, err := uuid.NewV7()
		if err != nil {
			return errors.New("there is a probleme with the generation of the uuid")
		}
		post.Id = uuid.String()

		day = rand.Intn(31-0) + 0
		mounth = rand.Intn(12-0) + 0
		year = rand.Intn(2024-1980) + 1980
		post.CreationDate = fmt.Sprintf("%d-%d-%d", year, mounth, day)

		post.Text = fmt.Sprintf("%s %s %s %s.", subjects[rand.Intn(len(subjects))], verbs[rand.Intn(len(verbs))], objects[rand.Intn(len(objects))], adverbs[rand.Intn(len(adverbs))])

		if err := InsertIntoDb("Post", db, post.Id, post.AuthorId, post.Text, post.Image, post.CreationDate, post.IsGroup); err != nil {
			i--
			continue
		}
	}

	return nil
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
			stringMAP += fmt.Sprintf("\"%v\", ", j)
		} else {
			stringMAP += fmt.Sprintf("\"%v\"", j)
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

	// Fetch the column types for correct handling of the result set
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// We loop on the result to stock them into the []map[string]any
	var result []map[string]any
	for rows.Next() {
		// We initialize the variable who gonna contain the current result row
		row := make(map[string]any)

		values := make([]interface{}, len(columnTypes))
		for i, ct := range columnTypes {
			switch ct.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "CHAR": // handle string types
				values[i] = new(string)
			case "INT", "INTEGER", "BIGINT": // handle integer types
				values[i] = new(int64)
			case "FLOAT", "DOUBLE", "REAL": // handle float types
				values[i] = new(float64)
			case "BOOL", "BOOLEAN": // handle boolean types
				values[i] = new(bool)
			default:
				values[i] = new(interface{}) // fallback for unknown types
			}
		}
			switch ct.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "CHAR": // handle string types
				values[i] = new(string)
			case "INT", "INTEGER", "BIGINT": // handle integer types
				values[i] = new(int64)
			case "FLOAT", "DOUBLE", "REAL": // handle float types
				values[i] = new(float64)
			case "BOOL", "BOOLEAN": // handle boolean types
				values[i] = new(bool)
			default:
				values[i] = new(interface{}) // fallback for unknown types
			}
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
  - an array of string who contains the name of the rows
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

func UpdateDb(tableName string, db *sql.DB, updateArgs map[string]any, whereArgs map[string]any) error {
	// Prepare the list of columns to update
	var colsToUpdate []string
	for col := range updateArgs {
		colsToUpdate = append(colsToUpdate, col)
	}

	// Combine updateArgs and whereArgs for the PrepareUpdateStmt function
	allArgs := make(map[string]any)
	for k, v := range updateArgs {
		allArgs[k] = v
	}
	for k, v := range whereArgs {
		allArgs[k] = v
	}

	// We prepare and execute the update request
	err := PrepareUpdateStmt(tableName, db, allArgs, colsToUpdate)
	return err
}

func PrepareUpdateStmt(tableName string, db *sql.DB, args map[string]any, colsToUpdate []string) (error) {
	var (
		setClauses   []string
		whereClauses []string
		params       []any
	)

	// Building the SET clause
	for _, col := range colsToUpdate {
		if value, ok := args[col]; ok {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", col))
			params = append(params, value)
			delete(args, col)
		}
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no columns to update")
	}

	// Building the WHERE clause with remaining parameters
	for column, value := range args {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", column))
		params = append(params, value)
	}

	// Construct the query
	query := fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(setClauses, ", "))
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Execute the query
	_, err := db.Exec(query, params...)
	return err
}

/*
This function takes 2 arguments:
  - a string who is the name of the table
  - a map with containing the wanted values in each row

The objective of this function is to Remove a row in the Db.

The function gonna return:
  - an error
*/
func RemoveFromDB(tabelName string, db *sql.DB, Args map[string]any) error {
	var whereCondition string
	var whereValues []any

	// We build the condition of the request
	for k, v := range Args {
		whereCondition += fmt.Sprintf("%s = ?", k)
		whereValues = append(whereValues, v)

		if len(whereValues) != len(Args) {
			whereCondition += " AND "
		}
	}

	if whereCondition != "" {
		whereCondition = "WHERE " + whereCondition
	}
	
	// We prepare the request
	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s %s", tabelName, whereCondition))
	if err != nil {
		return err
	}

	// We execute the SQL request
	_, err = stmt.Exec(whereValues...)
	if err != nil {
		return err
	}

	return nil
}