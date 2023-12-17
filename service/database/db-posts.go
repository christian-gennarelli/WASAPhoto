package database

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

func (db appdbimpl) CheckIfPostExists(PostID string) error {

	stmt, err := db.c.Prepare("SELECT PostID FROM Post WHERE PostID = ?")
	if err != nil {
		return err //fmt.Errorf("error encountered while preparing the query to check if the given post exists")
	}

	var id components.ID
	if err = stmt.QueryRow(PostID).Scan(&id.Value); err != nil {
		// if err == sql.ErrNoRows {
		// 	return err
		// }
		return err //fmt.Errorf("error encountered while executing the query to check if the given post exists")
	}

	return nil

}

func (db appdbimpl) CheckIfOwnerPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("SELECT P.Author FROM User U JOIN Post P ON U.Username = P.Author WHERE U.Username = ? AND P.PostID = ?")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to check if the given post is owned by the given user")
	}

	var author components.Username
	if err = stmt.QueryRow(Username, PostID).Scan(&author.Value); err != nil {
		// if err == sql.ErrNoRows {
		// 	return err
		// }
		return err //fmt.Errorf("error while executing the SQL query to retrieve the username associated to the given Auth token")
	}

	return nil

}

func (db appdbimpl) AddLikeToPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("INSERT INTO Like (PostID, Liker) VALUES (?, ?)")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to add the like")
	}

	_, err = stmt.Query(PostID, Username)
	if err != nil {
		return err //fmt.Errorf("error while executing the query to add the like")
	}

	return nil

}

func (db appdbimpl) RemoveLikeFromPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("DELETE FROM Like WHERE PostID = ? AND Liker = ?")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to add the like")
	}

	_, err = stmt.Query(PostID, Username)
	if err != nil {
		return err //fmt.Errorf("error while executing the query to add the like")
	}

	return nil

}

func (db appdbimpl) AddCommentToPost(PostID string, Body string, CreationDatetime string, Author string) error {

	stmt, err := db.c.Prepare("INSERT INTO Comment (PostID, Author, CreationDatetime, Comment) VALUES (?, ?, CONVERT(DATETIME, ?), ?)")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to add the comment")
	}

	_, err = stmt.Query(PostID, Author, CreationDatetime, Body)
	if err != nil {
		return err //fmt.Errorf("error while executing the query to add the comment")
	}

	return nil

}

func (db appdbimpl) RemoveCommentFromPost(PostID string, CommentID string) error {

	stmt, err := db.c.Prepare("DELETE FROM Comment WHERE PostID = ? AND CommentID = ?")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to remove the comment")
	}

	_, err = stmt.Exec(PostID, CommentID)
	if err != nil {
		return err //fmt.Errorf("error while executing the query to remove the comment")
	}

	return nil

}
