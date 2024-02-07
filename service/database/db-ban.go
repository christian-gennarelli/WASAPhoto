package database

import (
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

func (db appdbimpl) BanUser(bannerUsername string, bannedUsername string) error {

	stmt, err := db.c.Prepare("INSERT INTO Ban (Banner, Banned, CreationDatetime) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	t := time.Now()
	startDatetime := strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + " " + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
	if _, err = stmt.Exec(bannerUsername, bannedUsername, startDatetime); err != nil {
		return err
	}

	if err = db.UnfollowUser(bannerUsername, bannedUsername); err != nil {
		return err
	}

	stmt, err = db.c.Prepare("DELETE FROM Like WHERE Liker = ? AND PostID IN (SELECT P.PostID FROM Post P WHERE P.Author = ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(bannerUsername, bannedUsername); err != nil {
		return err
	}
	if _, err = stmt.Exec(bannedUsername, bannerUsername); err != nil {
		return err
	}

	stmt, err = db.c.Prepare("DELETE FROM Comment WHERE Author = ? AND PostID IN (SELECT P.PostID FROM Post P WHERE P.Author = ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(bannerUsername, bannedUsername); err != nil {
		return err
	}
	if _, err = stmt.Exec(bannedUsername, bannerUsername); err != nil {
		return err
	}

	stmt, err = db.c.Prepare("DELETE FROM Follow WHERE Follower = ? AND Followed = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(bannerUsername, bannedUsername); err != nil {
		return err
	}
	if _, err = stmt.Exec(bannedUsername, bannerUsername); err != nil {
		return err
	}

	return nil
}

func (db appdbimpl) UnbanUser(bannerUsername string, bannedUsername string) error {

	stmt, err := db.c.Prepare("DELETE FROM Ban WHERE Banner = ? AND Banned = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(bannerUsername, bannedUsername); err != nil {
		return err
	}

	return nil
}

func (db appdbimpl) GetBanUserList(bannerUsername string) (*[]components.User, error) {

	stmt, err := db.c.Prepare("SELECT U.Username, U.ProfilePicPath, COALESCE(U.Birthdate, ''), COALESCE(U.Name, '') FROM Ban B JOIN User U ON B.Banned = U.Username WHERE B.Banner = ? ORDER BY B.CreationDatetime DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(bannerUsername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bannedUserList []components.User
	for rows.Next() {
		var bannedUser components.User
		if err = rows.Scan(&bannedUser.Username, &bannedUser.ProfilePic, &bannedUser.Birthdate, &bannedUser.Name); err != nil {
			return nil, err
		}

		bannedUserList = append(bannedUserList, bannedUser)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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

	if err := row.Err(); err != nil {
		return err
	}

	if err = row.Scan(&foo); err != nil {
		return err
	}

	return nil
}
