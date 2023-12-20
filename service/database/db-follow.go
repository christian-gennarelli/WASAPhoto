package database

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

func (db appdbimpl) GetFollowersList(followedUsername string) (*components.UserList, error) {

	stmt, err := db.c.Prepare("SELECT U.Username, U.Birthdate, U.ProfilePicPath, U.Name FROM Follow F JOIN User U ON F.Follower = U.Username WHERE Followed = ?")
	if err != nil {
		return nil, err //fmt.Errorf("error while preparing the SQL statement to obtain the followers of the given username")
	}

	rows, err := stmt.Query(followedUsername)
	if err != nil && err != sql.ErrNoRows { // We don't care if no user follows the given one, we'll write the StatusNoContent header if it happens to be the case
		return nil, err //fmt.Errorf("error while executing the SQL statement to obtain the followers of the given username")
	}

	var userList components.UserList
	for rows.Next() {
		var user components.User
		err = rows.Scan(&user.Username, &user.Birthdate, &user.ProfilePic, &user.Name)
		if err != nil {
			return nil, err //fmt.Errorf("error while scanning the result of the query")
		}

		img, _ := os.Open(user.ProfilePic)
		reader := bufio.NewReader(img)
		content, _ := io.ReadAll(reader)
		user.ProfilePic = base64.StdEncoding.EncodeToString(content)

		userList.Users = append(userList.Users, user)
	}

	return &userList, nil

}

func (db appdbimpl) GetFollowingList(followerUsername string) (*components.UserList, error) {

	stmt, err := db.c.Prepare("SELECT U.Username, COALESCE(U.Birthdate, ''), U.ProfilePicPath, COALESCE(U.Name, '') FROM Follow F JOIN User U ON F.Followed = U.Username WHERE Follower = ?")
	if err != nil {
		return nil, err //fmt.Errorf("error while preparing the SQL statement to obtain the followers of the given username")
	}

	rows, err := stmt.Query(followerUsername)
	if err != nil && err != sql.ErrNoRows { // We don't care if no user follows the given one, we'll write the StatusNoContent header if it happens to be the case
		return nil, err //fmt.Errorf("error while executing the SQL statement to obtain the followers of the given username")
	}

	var userList components.UserList
	for rows.Next() {
		var user components.User
		err = rows.Scan(&user.Username, &user.Birthdate, &user.ProfilePic, &user.Name)
		if err != nil {
			return nil, err //fmt.Errorf("error while scanning the result of the query")
		}

		img, _ := os.Open(user.ProfilePic)
		reader := bufio.NewReader(img)
		content, _ := io.ReadAll(reader)
		user.ProfilePic = base64.StdEncoding.EncodeToString(content)

		userList.Users = append(userList.Users, user)
	}

	return &userList, nil

}

func (db appdbimpl) FollowUser(followerUsername string, followedUsername string) error {

	stmt, err := db.c.Prepare("INSERT INTO Follow (Follower, Followed) VALUES (?, ?)")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to add followerUsername to the list of followers of followedUsername")
	}

	if _, err = stmt.Exec(followerUsername, followedUsername); err != nil {
		return err //fmt.Errorf("error while executing the SQL statement to add followerUsername to the list of followers of followedUsername")
	}

	return nil

}

func (db appdbimpl) UnfollowUser(followerUsername string, followedUsername string) error {

	stmt, err := db.c.Prepare("DELETE FROM Follow WHERE Follower = ? AND Followed = ?")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to remove followerUsername from the list of followers of followedUsername")
	}

	if _, err = stmt.Exec(followerUsername, followedUsername); err != nil {
		return err //fmt.Errorf("error while executing the SQL statement to remove followerUsername from the list of followers of followedUsername")
	}

	return nil

}
