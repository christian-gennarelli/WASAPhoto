package database

import (
	"database/sql"
	"errors"

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

	// Retrieve the photos posted by this user
	stmt, err = db.c.Prepare(`SELECT 
									P.PostID, 
									P.Author, 
									P.Description, 
									P.CreationDatetime, 
									P.PhotoPath,
									(SELECT COUNT(*) FROM Like L WHERE L.PostID = P.PostID) as Likes 
							FROM Post P WHERE Author = ? ORDER BY P.CreationDatetime DESC`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	defer rows.Close()

	var posts []components.Post
	for rows.Next() {
		var post components.Post
		if err = rows.Scan(&post.PostID, &post.Author, &post.Description, &post.CreationDatetime, &post.Photo, &post.Likes); err != nil {
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
