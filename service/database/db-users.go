package database

import (
	"fmt"

	"github.com/dchest/uniuri"
)

// If the user does not exist, it will be created, and an identifier is returned. If the user exists, the user identifier is returned.
func (db appdbimpl) PostUserID(Username string) (ID string, err error) {

	// Prepare the SQL statement
	stmt, err := db.c.Prepare("SELECT ID from User WHERE Username = ?")
	if err != nil {
		return "", fmt.Errorf("error while preparing the SQL statement to obtain the id for the given user (it it exists)")
	}

	// Bind the parameters and execute the statement
	rows, err := stmt.Query(Username)
	if err != nil {
		return "", fmt.Errorf("error while performing the query to obtain the id for the given user (it it exists)")
	}

	defer rows.Close()

	// Check if the username already existed
	var id string

	// If yes, just return the associated id
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return "", fmt.Errorf("error while extracting the ID from the query")
		}
	}

	// If not, create a new user (and consequently a new ID for it)
	id = uniuri.NewLen(64)

	stmt, err = db.c.Prepare("INSERT INTO User (Username, ID) VALUES (?, ?)")
	if err != nil {
		return "", fmt.Errorf("error while preparing the SQL statement to create the new user")
	}

	_, err = stmt.Query(Username, id)
	if err != nil {
		return "", fmt.Errorf("error while performing the query to create the new user")
	}

	return id, nil

}
