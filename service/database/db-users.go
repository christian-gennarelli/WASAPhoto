package database

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/dchest/uniuri"
)

// If the user does not exist, it will be created, and an identifier is returned. If the user exists, the user identifier is returned.
func (db appdbimpl) PostUserID(Username string) (*string, error) {

	// Prepare the SQL statement
	stmt, err := db.c.Prepare("SELECT ID from User WHERE Username = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the id for the given user (if it exists)")
	}
	defer stmt.Close()

	var ID string
	// Bind the parameters and execute the statement
	row := stmt.QueryRow(Username)

	if err = row.Err(); err != nil {
		return nil, err
	}

	if err = row.Scan(&ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ID = uniuri.NewLen(64)

			// HOW CAN I STORE A []BYTE IN A SQL RECORD??
			// uuid = uuid.NewV4()
			// if err != nil {
			// 	return nil, err
			// }

			stmt, err = db.c.Prepare("INSERT INTO User (Username, ID) VALUES (?, ?)")
			if err != nil {
				return nil, err
			}

			if _, err = stmt.Exec(Username, ID); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &ID, nil
}

func (db appdbimpl) SearchUser(Username string) (*components.UserList, error) {

	// Prepare the SQL statement for finding all the users with "Value" as substring
	stmt, err := db.c.Prepare("SELECT Username, COALESCE(Birthdate, ''), ProfilePicPath, COALESCE(Name, '') FROM User WHERE Username LIKE '%'||?||'%' ")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the list of users with the provided string as substring")
	}
	defer stmt.Close()

	// Bind the parameters and execute the statement
	rows, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the list of users with the provided string as substring")
	}
	defer rows.Close()

	// Instantiate the data structure that will hold the list of usernames
	var ulist components.UserList

	// Loop over the rows, and store each user id in the previously instantiated data structure
	for rows.Next() {

		// Retrieve the next username
		var user components.User
		if err = rows.Scan(&user.Username, &user.Birthdate, &user.ProfilePic, &user.Name); err != nil {
			return nil, fmt.Errorf("error while extracting the username from the query")
		}

		//Open the image
		img, err := os.Open(user.ProfilePic)
		if err != nil {
			return nil, err
		}
		reader := bufio.NewReader(img)
		// Read it
		content, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		// Convert it in base64
		user.ProfilePic = base64.StdEncoding.EncodeToString(content)

		// Insert into the returned list of usernames
		ulist.Users = append(ulist.Users, user)

		if err = rows.Err(); err != nil {
			return nil, err
		}

	}

	// Return the list of users
	return &ulist, nil

}

func (db appdbimpl) UpdateUsername(NewUsername string, OldUsername string) error {

	stmt, err := db.c.Prepare("UPDATE User SET Username = ? WHERE Username = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(NewUsername, OldUsername); err != nil {
		return err
	}

	return nil

}

func (db appdbimpl) GetUsernameByToken(Id string) (*string, error) {

	stmt, err := db.c.Prepare("SELECT Username FROM User WHERE ID = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var username string
	err = stmt.QueryRow(Id).Scan(&username)
	if err != nil {
		return nil, err
	}

	return &username, nil

}

func (db appdbimpl) GetOwnerUsernameOfComment(Id string) (*string, error) {

	stmt, err := db.c.Prepare("SELECT Author FROM Comment WHERE CommentID = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var username string
	if err = stmt.QueryRow(Id).Scan(&username); err != nil {
		return nil, err
	}

	return &username, nil

}
