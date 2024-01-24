package database

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

func (db appdbimpl) CheckIfOwnerPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("SELECT P.Author FROM User U JOIN Post P ON U.Username = P.Author WHERE U.Username = ? AND P.PostID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var author string
	row := stmt.QueryRow(Username, PostID)
	if err = row.Scan(&author); err != nil {
		return err
	}

	if err = row.Err(); err != nil {
		return err
	}

	return nil

}

func (db appdbimpl) AddLikeToPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("INSERT INTO Like (PostID, Liker, CreationDatetime) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	t := time.Now()
	CreationDatetime := strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + " " + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
	_, err = stmt.Exec(PostID, Username, CreationDatetime)
	if err != nil {
		return err
	}

	return nil

}

func (db appdbimpl) RemoveLikeFromPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("DELETE FROM Like WHERE PostID = ? AND Liker = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(PostID, Username); err != nil {
		return err
	}

	return nil

}

func (db appdbimpl) AddCommentToPost(PostID string, Body string, Author string) error {

	stmt, err := db.c.Prepare("INSERT INTO Comment (PostID, Author, CreationDatetime, Comment) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	t := time.Now()
	CreationDatetime := strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + " " + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())

	if _, err = stmt.Exec(PostID, Author, CreationDatetime, Body); err != nil {
		return err
	}

	return nil

}

func (db appdbimpl) RemoveCommentFromPost(PostID string, CommentID string) error {

	stmt, err := db.c.Prepare("DELETE FROM Comment WHERE PostID = ? AND CommentID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(PostID, CommentID); err != nil {
		return err
	}

	return nil

}

func (db appdbimpl) GetUserStream(startDatetime string, username string) (*components.Stream, error) {

	stmt, err := db.c.Prepare(`SELECT 
									P.PostID, 
									P.Author, 
									P.CreationDatetime, 
									P.Description, 
									P.PhotoPath,
									(SELECT COUNT(*) FROM Like L WHERE L.PostID = P.PostID) as Likes 
							FROM Post P JOIN Follow F ON P.Author = F.Followed WHERE F.Follower = ? AND P.CreationDatetime <= ? ORDER BY P.CreationDatetime DESC LIMIT 16`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, startDatetime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postStream components.Stream
	for rows.Next() {
		var post components.Post
		if err := rows.Scan(&post.PostID, &post.Author, &post.CreationDatetime, &post.Description, &post.Photo, &post.Likes); err != nil {
			return nil, err
		}
		postStream.Posts = append(postStream.Posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &postStream, nil

}

func (db appdbimpl) UploadPost(username string, description string) (*components.Post, error) {

	var id int
	if err := db.c.QueryRow("SELECT PostID FROM Post ORDER BY PostID DESC LIMIT 1").Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			id = 0
		} else {
			return nil, err
		}
	}

	stmt, err := db.c.Prepare("INSERT INTO Post (Author, CreationDatetime, Description, PhotoPath) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	t := time.Now()
	creationDatetime := strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + " " + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
	photoPath := "posts/" + username + "_" + strconv.Itoa(id+1) + ".png"
	if _, err := stmt.Exec(username, creationDatetime, description, photoPath); err != nil {
		return nil, err
	}

	return &components.Post{
		PostID:           strconv.Itoa(id + 1),
		Author:           username,
		CreationDatetime: creationDatetime,
		Description:      description,
		Photo:            photoPath,
	}, nil

}

func (db appdbimpl) DeletePost(postID string) (*string, error) {

	stmt, err := db.c.Prepare("SELECT PhotoPath FROM Post WHERE PostID = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var photoPath string
	if err := stmt.QueryRow(postID).Scan(&photoPath); err != nil {
		return nil, err
	}

	if _, err := db.c.Exec("DELETE FROM Post WHERE PostID = ?", postID); err != nil {
		return nil, err
	}

	return &photoPath, nil

}

func (db appdbimpl) GetPostComments(postID string, startDatetime string) (*components.CommentList, error) {

	stmt, err := db.c.Prepare("SELECT * FROM Comment WHERE PostID = ? AND CreationDatetime <= ? ORDER BY CreationDatetime DESC LIMIT 16")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(postID, startDatetime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentList components.CommentList
	for rows.Next() {
		var comment components.Comment
		if err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.Author, &comment.CreationDatetime, &comment.Body); err != nil {
			return nil, err
		}
		commentList.Comments = append(commentList.Comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &commentList, nil

}

func (db appdbimpl) GetPostLikes(postID string, startDatetime string) (*components.UserList, error) {

	stmt, err := db.c.Prepare("SELECT U.Username, U.ProfilePicPath, COALESCE('', U.Birthdate), COALESCE('', U.Name) FROM User U JOIN Like L ON L.Liker = U.Username WHERE L.PostID = ? AND L.CreationDatetime <= ? ORDER BY CreationDatetime DESC LIMIT 16")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(postID, startDatetime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList components.UserList
	for rows.Next() {
		var user components.User
		if err := rows.Scan(&user.Username, &user.ProfilePic, &user.Birthdate, &user.Name); err != nil {
			return nil, err
		}
		userList.Users = append(userList.Users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &userList, nil

}
