/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// Default methods
	GetName() (string, error)
	SetName(name string) error
	Ping() error

	// Custom methods
	GetUsernameByToken(Id string) (*components.Username, error)
	CheckIfUsernameExists(Username string) (*bool, error)
	CheckIfPostExists(PostID string) (*bool, error)
	CheckCombinationIsValid(Username string, ID string) (*bool, error)
	CheckIfOwnerPost(Username string, PostID string) (*bool, error)
	PostUserID(Username string) (*components.ID, error)
	SearchUser(Username string) (*components.UserList, error)
	GetUserProfile(Username string) (*components.Profile, error)
	UpdateUsername(NewUsername string, OldUsername string) error
	AddLikeToPost(Username string, PostID string) error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Read the content of the SQL statement needed to build up the structure of the database (if not already present)
	file, err := os.Open("database_setup.txt")
	if err != nil {
		fmt.Println("Error encoutered while reading the database SQL statement from the corresponding file!")
		return nil, err
	}
	defer file.Close()

	// What bufio.Scanner.Scan() does is to split the content of the file according in so-called tokens, which are unit of text separated by a delimiter (\n is the default one.)
	var db_setup string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() // bufio.Scanner.Text() return the content of the current token.
		db_setup += line + "\n"
	}

	if scanner.Err() != nil {
		fmt.Println("Error from the scanner!")
		return nil, scanner.Err()
	}

	// Run the following SQL statement to build the structure of the database
	stmt, err := db.Prepare(db_setup)
	if err != nil {
		fmt.Println("Error while preparing the SQL statement for preparing the database!")
		return nil, err
	}

	// Execute the SQL statement
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("Error while executing the SQL statement to prepare the database!")
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
