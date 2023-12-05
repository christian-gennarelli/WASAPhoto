package database

import (
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/dchest/uniuri"
)

func (db appdbimpl) IsValid(Username string, ID string) (Valid *bool, err error) {

	// Prepare the SQL statement to return the row containing both the provided username and id, if there is any
	stmt, err := db.c.Prepare("SELECT COUNT(*) FROM Users WHERE Username = ? AND ID = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the list of users with the provided string as substring")
	}

	// Bind the parameters and execute the statement
	rows, err := stmt.Query(Username, ID)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the list of users with the provided string as substring")
	} else {
		defer rows.Close()
	}

	// Check if the returned value is exactly 1: if yes, then the user is valid
	valid := false
	if rows.Next() {
		var numRows string
		err := rows.Scan(&numRows)
		if err != nil {
			return nil, fmt.Errorf("error while parsing the number of rows where username and id coincides with the ones provided")
		} else {
			valid = numRows == "1"
		}
	}

	return &valid, nil

}

func (db appdbimpl) CheckIfUsernameExists(Username string) (*bool, error) {

	stmt, err := db.c.Prepare("SELECT 1 FROM Users WHERE Username = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the check if the provided username exists")
	}

	rows, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the list of users with the provided string as substring")
	} else {
		defer rows.Close()
	}

	valid := false
	if rows.Next() {
		var result string
		err := rows.Scan(&result)
		if err != nil {
			return nil, fmt.Errorf("error while parsing the number of rows where username and id coincides with the ones provided")
		} else {
			valid = result == "1"
		}
	}

	return &valid, nil

}

// If the user does not exist, it will be created, and an identifier is returned. If the user exists, the user identifier is returned.
func (db appdbimpl) PostUserID(Username string) (ID *components.ID, err error) {

	// Prepare the SQL statement
	stmt, err := db.c.Prepare("SELECT ID from User WHERE Username = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the id for the given user (it it exists)")
	}

	// Bind the parameters and execute the statement
	rows, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the id for the given user (it it exists)")
	} else {
		defer rows.Close()
	}

	// Check if the username already existed
	var id components.ID

	// If yes, just return the associated id
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error while extracting the ID from the query")
		}
	} else { // If not, create a new user (and consequently a new ID for it)

		var id components.ID
		id.RandID = uniuri.NewLen(64)

		stmt, err = db.c.Prepare("INSERT INTO User (Username, ID) VALUES (?, ?)")
		if err != nil {
			return nil, fmt.Errorf("error while preparing the SQL statement to create the new user")
		}

		_, err = stmt.Query(Username, id)
		if err != nil {
			return nil, fmt.Errorf("error while performing the query to create the new user")
		}

	}

	return &id, nil

}

func (db appdbimpl) SearchUser(Username string) (*components.UserList, error) {

	// Prepare the SQL statement for finding all the users with "uname" as substring
	stmt, err := db.c.Prepare("SELECT Username FROM Users WHERE Username LIKE '%?%'")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the list of users with the provided string as substring")
	}

	// Bind the parameters and execute the statement
	rows, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the list of users with the provided string as substring")
	} else {
		defer rows.Close()
	}

	// Instantiate the data structure that will hold the list of usernames
	var ulist components.UserList

	// Loop over the rows, and store each user id in the previously instantiated data structure
	for rows.Next() {

		// Retrieve the next username
		var user components.Username
		err = rows.Scan(&user.Uname)
		if err != nil {
			return nil, fmt.Errorf("error while extracting the username from the query")
		}

		// Insert into the returned list of usernames
		ulist.Users = append(ulist.Users, user)

	}

	// Return the list of users
	return &ulist, nil

}

// Retrieve the profile of the user with the provided username
func (db appdbimpl) GetUserProfile(Username string) (*components.Profile, error) {

	// Retrieve the photos posted by this user
	stmt, err := db.c.Prepare("SELECT PostID FROM Post WHERE Author = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the list of posts posted by the user")
	}

	postIDs, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the list of posts posted by the user")
	} else {
		defer postIDs.Close()
	}

	var Posts []components.ID
	for postIDs.Next() {
		var postID components.ID
		err = postIDs.Scan(&postID)
		if err != nil {
			return nil, fmt.Errorf("error while extracting the username from the query")
		}
		Posts = append(Posts, postID)
	}

	// Retrieve the informations about the user with the provided username
	stmt, err = db.c.Prepare("SELECT * FROM Users WHERE Username = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the info about the user with the provided username")
	}

	users, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the info about the user with the provided username")
	} else {
		defer postIDs.Close()
	}

	var User components.User
	for users.Next() {
		err = postIDs.Scan(&User)
		if err != nil {
			return nil, fmt.Errorf("error while extracting the username from the query")
		}
	}

	profile := components.Profile{
		User:  User,
		Posts: Posts,
	}

	return &profile, nil

}
