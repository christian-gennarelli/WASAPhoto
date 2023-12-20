package database

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

func (db appdbimpl) BanUser(bannerUsername string, bannedUsername string) error {

	stmt, err := db.c.Prepare("INSERT INTO Ban (Banner, Banned) VALUES (?, ?)")
	if err != nil {
		return err //fmt.Errorf("error while preparing the statement to ban the user")
	}

	if _, err = stmt.Exec(bannerUsername, bannedUsername); err != nil {
		return err //fmt.Errorf("errof while executing the statement to ban the user")
	}

	if err = db.UnfollowUser(bannerUsername, bannedUsername); err != nil {
		return err
	}

	return nil
}

func (db appdbimpl) UnbanUser(bannerUsername string, bannedUsername string) error {

	stmt, err := db.c.Prepare("DELETE FROM Ban WHERE Banner = ? AND Banned = ?")
	if err != nil {
		return err //fmt.Errorf("error while preparing the statement to ban the user")
	}

	if _, err = stmt.Exec(bannerUsername, bannedUsername); err != nil {
		return err //fmt.Errorf("errof while executing the statement to ban the user")
	}

	return nil
}

func (db appdbimpl) GetBanUserList(bannerUsername string) (*components.UserList, error) {

	stmt, err := db.c.Prepare("SELECT U.Username, U.ProfilePicPath, COALESCE(U.Birthdate, ''), COALESCE(U.Name, '') FROM Ban B JOIN User U ON B.Banned = U.Username WHERE Banner = ?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(bannerUsername)
	if err != nil {
		return nil, err
	}

	var bannedUserList components.UserList
	for rows.Next() {
		var bannedUser components.User
		if err = rows.Scan(&bannedUser.Username, &bannedUser.ProfilePic, &bannedUser.Birthdate, &bannedUser.Name); err != nil {
			return nil, err
		}

		// // Open image and turn it into base64
		img, _ := os.Open(bannedUser.ProfilePic)
		reader := bufio.NewReader(img)
		content, _ := io.ReadAll(reader)
		bannedUser.ProfilePic = base64.StdEncoding.EncodeToString(content)

		bannedUserList.Users = append(bannedUserList.Users, bannedUser)
	}

	return &bannedUserList, nil

}

func (db appdbimpl) CheckIfBanned(bannerUsername string, bannedUsername string) (*bool, error) {

	stmt, err := db.c.Prepare("SELECT Banned FROM Ban WHERE Banner = ? AND Banned = ?")
	if err != nil {
		return nil, err
	}

	row, foo, valid := stmt.QueryRow(bannerUsername, bannedUsername), "", true
	if err = row.Scan(&foo); err != nil {
		if err == sql.ErrNoRows {
			return &valid, nil
		}
		return nil, err
	}

	valid = false
	return &valid, nil
}
