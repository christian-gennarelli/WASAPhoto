package database

import (
	"database/sql"
	"fmt"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

// Retrieve the profile of the user with the provided username
func (db appdbimpl) GetUserProfile(Username string) (*components.Profile, error) {

	// Retrieve the informations about the user with the provided username
	stmt, err := db.c.Prepare("SELECT Username, ID FROM User WHERE Username = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the info about the user with the provided username")
	}

	var User components.User
	if err = stmt.QueryRow(Username).Scan(&User.Username, &User.ID); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("error while executing the SQL statement to obtain the info about the user with the provided username")
	}

	// Retrieve the photos posted by this user
	stmt, err = db.c.Prepare("SELECT PostID FROM Post WHERE Author = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the list of posts posted by the user")
	}

	postIDs, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the list of posts posted by the user")
	}
	defer postIDs.Close()

	var Posts []components.ID
	for postIDs.Next() {
		var postID components.ID
		if err = postIDs.Scan(&postID); err != nil {
			return nil, fmt.Errorf("error while extracting the posts from the query")
		}
		Posts = append(Posts, postID)
	}

	profile := components.Profile{
		User:  User,
		Posts: Posts,
	}

	return &profile, nil

}
