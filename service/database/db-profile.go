package database

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

// Retrieve the profile of the user with the provided username
func (db appdbimpl) GetUserProfile(Username string) (*components.Profile, error) {

	// Retrieve the informations about the user with the provided username
	stmt, err := db.c.Prepare("SELECT Username, COALESCE(Birthdate, ''), COALESCE(Name, ''), ProfilePicPath FROM User WHERE Username = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user components.User
	if err = stmt.QueryRow(Username).Scan(&user.Username, &user.Birthdate, &user.Name, &user.ProfilePic); err != nil {
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

	// Retrieve the photos posted by this user
	stmt, err = db.c.Prepare("SELECT PostID, Author, Description, CreationDatetime, PhotoPath FROM Post WHERE Author = ?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	defer rows.Close()

	var posts []components.Post
	for rows.Next() {
		var post components.Post
		if err = rows.Scan(&post.PostID.Value, &post.Author.Value, &post.Description, &post.CreationDatetime, &post.Photo); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	profile := components.Profile{
		User:  user,
		Posts: posts,
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &profile, nil

}
