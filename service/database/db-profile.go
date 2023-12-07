package database

import (
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

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
	stmt, err = db.c.Prepare("SELECT * FROM User WHERE Username = ?")
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
