package database

import "fmt"

func (db appdbimpl) FollowUser(followerUsername string, followingUsername string) error {

	stmt, err := db.c.Prepare("INSERT INTO Follow (Follower, Following) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to add followerUsername to the list of followers of followingUsername")
	}

	_, err = stmt.Query(followerUsername, followingUsername)
	if err != nil {
		return fmt.Errorf("error while executing the SQL statement to add followerUsername to the list of followers of followingUsername")
	}

	return nil

}

func (db appdbimpl) UnfollowUser(followerUsername string, followingUsername string) error {

	stmt, err := db.c.Prepare("DELETE FROM Follow WHERE Follower = ? AND Following = ?)")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to remove followerUsername from the list of followers of followingUsername")
	}

	_, err = stmt.Query(followerUsername, followingUsername)
	if err != nil {
		return fmt.Errorf("error while executing the SQL statement to remove followerUsername from the list of followers of followingUsername")
	}

	return nil

}
