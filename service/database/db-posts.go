package database

import (
	"database/sql"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

func (db appdbimpl) CheckIfOwnerPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("SELECT P.Author FROM User U JOIN Post P ON U.Username = P.Author WHERE U.Username = ? AND P.PostID = ?")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to check if the given post is owned by the given user")
	}
	defer stmt.Close()

	var author components.Username
	if err = stmt.QueryRow(Username, PostID).Scan(&author.Value); err != nil {
		return err //fmt.Errorf("error while executing the SQL query to retrieve the username associated to the given Auth token")
	}

	return nil

}

func (db appdbimpl) AddLikeToPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("INSERT INTO Like (PostID, Liker) VALUES (?, ?)")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to add the like")
	}
	defer stmt.Close()

	_, err = stmt.Exec(PostID, Username)
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
	defer stmt.Close()

	if _, err = stmt.Exec(PostID, Username); err != nil {
		return err //fmt.Errorf("error while executing the query to add the like")
	}

	return nil

}

func (db appdbimpl) AddCommentToPost(PostID string, Body string, CreationDatetime string, Author string) error {

	stmt, err := db.c.Prepare("INSERT INTO Comment (PostID, Author, CreationDatetime, Comment) VALUES (?, ?, CONVERT(DATETIME, ?), ?)")
	if err != nil {
		return err //fmt.Errorf("error while preparing the SQL statement to add the comment")
	}
	defer stmt.Close()

	_, err = stmt.Exec(PostID, Author, CreationDatetime, Body)
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
	defer stmt.Close()

	_, err = stmt.Exec(PostID, CommentID)
	if err != nil {
		return err //fmt.Errorf("error while executing the query to remove the comment")
	}

	return nil

}

func (db appdbimpl) GetUserStream(startDatetime string, username string) (*components.Stream, error) {

	stmt, err := db.c.Prepare("SELECT P.PostID, P.Author, P.CreationDatetime, P.Description, P.PhotoPath FROM Post P JOIN Follow F ON P.Author = F.Followed WHERE F.Follower = ? AND P.CreationDatetime <= ? ORDER BY P.CreationDatetime DESC LIMIT 16")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, startDatetime)
	if err != nil {
		return nil, err
	}

	var postStream components.Stream
	for rows.Next() {
		var post components.Post
		if err := rows.Scan(&post.PostID.Value, &post.Author.Value, &post.CreationDatetime, &post.Description, &post.Photo); err != nil {
			return nil, err
		}
		postStream.Posts = append(postStream.Posts, post)
	}

	return &postStream, nil

}

func (db appdbimpl) UploadPost(username string, description string) (error, *components.Post) {

	var id int
	if err := db.c.QueryRow("SELECT PostID FROM Post ORDER BY PostID DESC LIMIT 1").Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			id = 0
		} else {
			return err, nil
		}
	}

	stmt, err := db.c.Prepare("INSERT INTO Post (Author, CreationDatetime, Description, PhotoPath) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err, nil
	}

	t := time.Now()
	creationDatetime := strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + " " + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
	photoPath := "photos/posts/" + username + "_" + strconv.Itoa(id+1) + ".png"
	if _, err := stmt.Exec(username, creationDatetime, description, photoPath); err != nil {
		return err, nil
	}

	return nil, &components.Post{
		PostID:           components.ID{Value: strconv.Itoa(id + 1)},
		Author:           components.Username{Value: username},
		CreationDatetime: creationDatetime,
		Description:      description,
		Photo:            photoPath,
	}

}

func (db appdbimpl) DeletePost(postID string) (*string, error) {

	stmt, err := db.c.Prepare("SELECT PhotoPath FROM Post WHERE PostID = ?")
	if err != nil {
		return nil, err
	}

	var photoPath string
	if err := stmt.QueryRow(postID).Scan(&photoPath); err != nil {
		return nil, err
	}

	if _, err := db.c.Exec("DELETE FROM Post WHERE PostID = ?", postID); err != nil {
		return nil, err
	}

	return &photoPath, nil

}
