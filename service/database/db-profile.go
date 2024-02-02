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
									P.PhotoPath
							FROM Post P WHERE Author = ? ORDER BY P.CreationDatetime DESC`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	//(SELECT COUNT(*) FROM Like L WHERE L.PostID = P.PostID) as Likes

	rows, err := stmt.Query(Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	defer rows.Close()

	var posts []components.Post
	for rows.Next() {
		var post components.Post
		if err = rows.Scan(&post.PostID, &post.Author, &post.Description, &post.CreationDatetime, &post.Photo); err != nil {
			return nil, err
		}

		likers, err := db.GetPostLikes(post.PostID)
		if err != nil {
			return nil, err
		}
		post.Likes = *likers

		comments, err := db.GetPostComments(post.PostID)
		if err != nil {
			return nil, err
		}
		post.Comments = *comments

		posts = append(posts, post)
	}

	followings, err := db.GetFollowingList(Username)
	if err != nil {
		return nil, err
	}

	followers, err := db.GetFollowersList(Username)
	if err != nil {
		return nil, err
	}

	banned, err := db.GetBanUserList(Username)
	if err != nil {
		return nil, err
	}

	profile := components.Profile{
		User:       user,
		Posts:      posts,
		Followings: *followings,
		Followers:  *followers,
		Banned:     *banned,
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &profile, nil

}
