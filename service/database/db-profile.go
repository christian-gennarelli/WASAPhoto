package database

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

// Retrieve the profile of the user with the provided username
func (db appdbimpl) GetUserProfile(Username string) (*components.Profile, error) {

	// Retrieve the informations about the user with the provided username
	stmt, err := db.c.Prepare("SELECT Username, COALESCE(Birthdate, ''), COALESCE(Name, ''), ProfilePicPath FROM User WHERE Username = ?")
	if err != nil {
		return nil, err // fmt.Errorf("error while preparing the SQL statement to obtain the info about the user with the provided username")
	}
	defer stmt.Close()

	var user components.User
	if err = stmt.QueryRow(Username).Scan(&user.Username, &user.Birthdate, &user.Name, &user.ProfilePic); err != nil {
		return nil, err //fmt.Errorf("error while executing the SQL statement to obtain the info about the user with the provided username")
	}

	// Open image and turn it into base64
	img, _ := os.Open(user.ProfilePic)
	reader := bufio.NewReader(img)
	content, _ := io.ReadAll(reader)
	user.ProfilePic = base64.StdEncoding.EncodeToString(content)

	// Retrieve the photos posted by this user
	stmt, err = db.c.Prepare("SELECT PostID, Author, Description, CreationDatetime, PhotoPath FROM Post WHERE Author = ?")
	if err != nil {
		return nil, err //fmt.Errorf("error while preparing the SQL statement to obtain the list of posts posted by the user")
	}

	rows, err := stmt.Query(Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err //fmt.Errorf("error while performing the query to obtain the list of posts posted by the user")
	}
	defer rows.Close()

	var posts []components.Post
	for rows.Next() {
		var post components.Post
		if err = rows.Scan(&post.PostID.Value, &post.Author.Value, &post.Description, &post.CreationDatetime, &post.Photo); err != nil {
			return nil, err //fmt.Errorf("error while extracting the posts from the query")
		}
		posts = append(posts, post)
	}

	profile := components.Profile{
		User:  user,
		Posts: posts,
	}

	return &profile, nil

}
