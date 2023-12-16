package database

import (
	"database/sql"
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/dchest/uniuri"
)

// If the user does not exist, it will be created, and an identifier is returned. If the user exists, the user identifier is returned.
func (db appdbimpl) PostUserID(Username string) (*components.ID, error) {

	// Prepare the SQL statement
	stmt, err := db.c.Prepare("SELECT ID from User WHERE Username = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the id for the given user (if it exists)")
	}

	var id string
	var ID components.ID
	// Bind the parameters and execute the statement
	err = stmt.QueryRow(Username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			ID.RandID = uniuri.NewLen(64)

			stmt, err = db.c.Prepare("INSERT INTO User (Username, ID, Birthdate) VALUES (?, ?, ?)")
			if err != nil {
				return nil, fmt.Errorf("error while preparing the SQL statement to create the new user")
			}

			if _, err = stmt.Exec(Username, ID.RandID, "2023-12-15"); err != nil {
				return nil, fmt.Errorf("error while performing the query to create the new user")
			}
		} else {
			return nil, fmt.Errorf("error while performing the query to obtain the id for the given user (if it exists)")
		}
	} else {
		ID.RandID = id
	}

	return &ID, nil
}

func (db appdbimpl) SearchUser(Username string) (*components.UserList, error) {

	// Prepare the SQL statement for finding all the users with "uname" as substring
	stmt, err := db.c.Prepare("SELECT Username FROM User WHERE Username LIKE '%?%'")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the list of users with the provided string as substring")
	}

	// Bind the parameters and execute the statement
	users, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the list of users with the provided string as substring")
	}
	defer users.Close()

	// Instantiate the data structure that will hold the list of usernames
	var ulist components.UserList

	// Loop over the rows, and store each user id in the previously instantiated data structure
	for users.Next() {

		// Retrieve the next username
		var user components.Username
		err = users.Scan(&user.Uname)
		if err != nil {
			return nil, fmt.Errorf("error while extracting the username from the query")
		}

		// Insert into the returned list of usernames
		ulist.Users = append(ulist.Users, user)

	}

	// Return the list of users
	return &ulist, nil

}

func (db appdbimpl) UpdateUsername(NewUsername string, OldUsername string) error {

	stmt, err := db.c.Prepare("UPDATE User SET Username = ? WHERE Username = ?")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to updating the username")
	}

	_, err = stmt.Exec(NewUsername, OldUsername)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return err
		// }
		return err //fmt.Errorf("error while performing the query to obtain the info about the user with the provided username")
	}

	return nil

}

func (db appdbimpl) GetUsernameByToken(Id string) (*components.Username, error) {

	stmt, err := db.c.Prepare("SELECT Username FROM User WHERE ID = ?")
	if err != nil {
		return nil, err //fmt.Errorf("error encountered while preparing the query to retrieve the username associated with the given token")
	}

	var username components.Username
	err = stmt.QueryRow(Id).Scan(&username.Uname)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, err
		// }
		return nil, err //fmt.Errorf("error while executing the query to retrieve the username associated with the given token")
	}

	return &username, nil

}

func (db appdbimpl) GetOwnerUsernameOfComment(Id string) (*components.Username, error) {

	stmt, err := db.c.Prepare("SELECT Author FROM Comment WHERE CommentID = ?")
	if err != nil {
		return nil, err //fmt.Errorf("error while preparing the SQL statement to retrieve the author of the provided comment")
	}

	var username components.Username
	if err = stmt.QueryRow(Id).Scan(&username.Uname); err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, err
		// }
		return nil, err //fmt.Errorf("error while executing the SQL statement to retrieve the author of the provided comment")
	}

	return &username, nil

}
