package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os/exec"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/dchest/uniuri"
)

// If the user does not exist, it will be created, and an identifier is returned. If the user exists, the user identifier is returned.
func (db appdbimpl) PostUserID(Username string) (*components.User, error) {

	// Prepare the SQL statement
	stmt, err := db.c.Prepare("SELECT ID, Username, ProfilePicPath, COALESCE(Birthdate, ''), COALESCE(Name, '') from User WHERE Username = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the id for the given user (if it exists)")
	}
	defer stmt.Close()

	var user components.User
	// Bind the parameters and execute the statement
	row := stmt.QueryRow(Username)

	if err = row.Err(); err != nil {
		return nil, err
	}

	if err = row.Scan(&user.ID, &user.Username, &user.ProfilePic, &user.Birthdate, &user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			user.ID = uniuri.NewLen(64)
			user.Username = Username
			user.ProfilePic = "profile_pics/" + user.Username + ".png"

			stmt, err = db.c.Prepare("INSERT INTO User (Username, ID, ProfilePicPath) VALUES (?, ?, ?)")
			if err != nil {
				return nil, err
			}
			defer stmt.Close()

			if _, err = stmt.Exec(user.Username, user.ID, user.ProfilePic); err != nil {
				return nil, err
			}

			// Make a copy of the default profile pic and rename it
			srcFolder := "photos/profile_pics/default.png"
			destFolder := "photos/" + user.ProfilePic
			cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
			if err := cpCmd.Run(); err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (db appdbimpl) UpdateUsername(NewUsername string, OldUsername string) error {

	stmt, err := db.c.Prepare("UPDATE User SET Username = ?, ProfilePicPath = ? WHERE Username = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(NewUsername, "profile_pics/"+NewUsername+".png", OldUsername); err != nil {
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
