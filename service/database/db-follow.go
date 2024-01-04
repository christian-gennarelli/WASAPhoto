package database

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

func (db appdbimpl) GetFollowersList(followedUsername string, startDatetime string) (*components.UserList, error) {

	stmt, err := db.c.Prepare("SELECT U.Username, COALESCE('', U.Birthdate), U.ProfilePicPath,  COALESCE('', U.Name) FROM Follow F JOIN User U ON F.Follower = U.Username WHERE F.Followed = ? AND F.CreationDatetime <= ? ORDER BY F.CreationDatetime DESC LIMIT 16")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(followedUsername, startDatetime)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { // We don't care if no user follows the given one, we'll write the StatusNoContent header if it happens to be the case
		return nil, err
	}
	defer rows.Close()

	var userList components.UserList
	for rows.Next() {
		var user components.User
		err = rows.Scan(&user.Username.Value, &user.Birthdate, &user.ProfilePic, &user.Name)
		if err != nil {
			return nil, err
		}

		// Open the image
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

		userList.Users = append(userList.Users, user)
	}

	return &userList, nil

}

func (db appdbimpl) GetFollowingList(followerUsername string, startDatetime string) (*components.UserList, error) {

	stmt, err := db.c.Prepare("SELECT U.Username, COALESCE(U.Birthdate, ''), U.ProfilePicPath, COALESCE(U.Name, '') FROM Follow F JOIN User U ON F.Followed = U.Username WHERE F.Follower = ? AND F.CreationDatetime <= ? ORDER BY F.CreationDatetime DESC LIMIT 16")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(followerUsername, startDatetime)
	if err != nil && err != sql.ErrNoRows { // We don't care if no user follows the given one, we'll write the StatusNoContent header if it happens to be the case
		return nil, err
	}
	defer rows.Close()

	var userList components.UserList
	for rows.Next() {
		var user components.User
		err = rows.Scan(&user.Username.Value, &user.Birthdate, &user.ProfilePic, &user.Name)
		if err != nil {
			return nil, err
		}

		// Open the image
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

		userList.Users = append(userList.Users, user)
	}

	return &userList, nil

}

func (db appdbimpl) FollowUser(followerUsername string, followedUsername string) error {

	stmt, err := db.c.Prepare("INSERT INTO Follow (Follower, Followed, CreationDatetime) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	t := time.Now()
	startDatetime := strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + " " + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
	if _, err = stmt.Exec(followerUsername, followedUsername, startDatetime); err != nil {
		return err
	}

	if _, err = stmt.Exec(followerUsername, followedUsername); err != nil {
		return err
	}

	return nil

}

func (db appdbimpl) UnfollowUser(followerUsername string, followedUsername string) error {

	stmt, err := db.c.Prepare("DELETE FROM Follow WHERE Follower = ? AND Followed = ?")
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(followerUsername, followedUsername); err != nil {
		return err
	}

	return nil

}
