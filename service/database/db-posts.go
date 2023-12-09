package database

import (
	"fmt"

	"github.com/dchest/uniuri"
)

func (db appdbimpl) CheckIfPostExists(PostID string) (*bool, error) {

	stmt, err := db.c.Prepare("SELECT 1 FROM Post WHERE PostID = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the check if the provided post exists")
	}

	rows, err := stmt.Query(PostID)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to check if the provided post exists")
	} else {
		defer rows.Close()
	}

	valid := true
	if !rows.Next() {
		valid = false
	}

	return &valid, nil

}

func (db appdbimpl) CheckIfOwnerPost(Username string, PostID string) (*bool, error) {

	stmt, err := db.c.Prepare("SELECT 1 FROM User U JOIN Post P ON U.Username = P.Author WHERE U.Username = ? AND P.PostID = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to check if the given post is owned by the given user")
	}

	rows, err := stmt.Query(Username, PostID)
	if err != nil {
		return nil, fmt.Errorf("error while executing the SQL query to retrieve the username associated to the given Auth token")
	}

	valid := true
	if !rows.Next() {
		valid = false
	}

	return &valid, nil

}

func (db appdbimpl) AddLikeToPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("INSERT INTO Like (PostID, Liker) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to add the like")
	}

	_, err = stmt.Query(PostID, Username)
	if err != nil {
		return fmt.Errorf("error while executing the query to add the like")
	}

	return nil

}

func (db appdbimpl) RemoveLikeFromPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("DELETE FROM Like WHERE PostID = ? AND Liker = ?  ")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to add the like")
	}

	_, err = stmt.Query(PostID, Username)
	if err != nil {
		return fmt.Errorf("error while executing the query to add the like")
	}

	return nil

}

func (db appdbimpl) AddCommentToPost(PostID string, Body string, CreationDatetime string, Author string) error {

	stmt, err := db.c.Prepare("INSERT INTO Comment (CommentID, PostID, Author, CreationDatetime, Comment) VALUES (?, ?, ?, CONVERT(DATETIME, ?), ?)")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to add the comment")
	}

	// Generate the comment id
	commentID := uniuri.NewLen(64)

	_, err = stmt.Exec(commentID, PostID, Author, CreationDatetime, Body)
	if err != nil {
		return fmt.Errorf("error while executing the query to add the comment")
	}

	return nil

}

func (db appdbimpl) RemoveCommentFromPost(PostID string, CommentID string) error {

	stmt, err := db.c.Prepare("DELETE FROM Comment WHERE PostID = ? AND CommentID = ?")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to remove the comment")
	}

	_, err = stmt.Exec(PostID, CommentID)
	if err != nil {
		return fmt.Errorf("error while executing the query to remove the comment")
	}

	return nil

}