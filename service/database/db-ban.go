package database

import (
	"bufio"
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
	defer stmt.Close()

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
	defer stmt.Close()

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
	defer stmt.Close()

	rows, err := stmt.Query(bannerUsername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// Returns nil if banned (either of the two directions of the ban), err for some internal server error, sql.ErrNoRows if no ban has been found
func (db appdbimpl) CheckIfBanned(bannerUsername string, bannedUsername string) error {

	stmt, err := db.c.Prepare("SELECT Banned FROM Ban WHERE Banner = ? AND Banned = ? UNION SELECT Banner FROM Ban WHERE Banner = ? AND Banned = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	row, foo := stmt.QueryRow(bannerUsername, bannedUsername, bannedUsername, bannerUsername), ""
	if err = row.Scan(&foo); err != nil {
		return err
	}

	return nil
}
