package database

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

func (db appdbimpl) GetFollowersList(followedUsername string) (*[]components.User, error) {

	stmt, err := db.c.Prepare("SELECT U.Username, COALESCE('', U.Birthdate), U.ProfilePicPath,  COALESCE('', U.Name) FROM Follow F JOIN User U ON F.Follower = U.Username WHERE F.Followed = ? ORDER BY F.CreationDatetime DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(followedUsername)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { // We don't care if no user follows the given one, we'll write the StatusNoContent header if it happens to be the case
		return nil, err
	}
	defer rows.Close()

	var userList []components.User
	for rows.Next() {
		var user components.User
		err = rows.Scan(&user.Username, &user.Birthdate, &user.ProfilePic, &user.Name)
		if err != nil {
			return nil, err
		}

		userList = append(userList, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &userList, nil

}

func (db appdbimpl) GetFollowingList(followerUsername string) (*[]components.User, error) {

	stmt, err := db.c.Prepare("SELECT U.Username, COALESCE(U.Birthdate, ''), U.ProfilePicPath, COALESCE(U.Name, '') FROM Follow F JOIN User U ON F.Followed = U.Username WHERE F.Follower = ? ORDER BY F.CreationDatetime DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(followerUsername)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { // We don't care if no user follows the given one, we'll write the StatusNoContent header if it happens to be the case
		return nil, err
	}
	defer rows.Close()

	var userList []components.User
	for rows.Next() {
		var user components.User
		err = rows.Scan(&user.Username, &user.Birthdate, &user.ProfilePic, &user.Name)
		if err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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

	return nil

}

func (db appdbimpl) UnfollowUser(followerUsername string, followedUsername string) error {

	stmt, err := db.c.Prepare("DELETE FROM Follow WHERE Follower = ? AND Followed = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(followerUsername, followedUsername); err != nil {
		return err
	}

	return nil

}
