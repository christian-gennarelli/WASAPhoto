// Here (almost) all the schemas defined in the API are defined.

package components

import (
	"regexp"
)

type User struct {
	ID         string
	Username   string
	Birthdate  string
	Name       string
	ProfilePic string // Base64 encoded image
}

type Profile struct {
	User       User
	Posts      []Post
	Followings []User
	Followers  []User
	Banned     []User
}

type UserList struct {
	Users []User
}

type Post struct {
	PostID           string
	Author           string
	Photo            string // URL path to the image, stored server-side
	CreationDatetime string
	Description      string
	Likes            []User
}

type Stream struct {
	Posts []Post
}

type Comment struct {
	CommentID        string
	PostID           string
	Body             string
	CreationDatetime string
	Author           string
}

type CommentList struct {
	Comments []Comment
}

type Error struct {
	ErrorCode   int
	Description string
}

// Check if the provided username is in the correct format
func CheckIfValid(content string, contentType string) error {

	var REGEXP string
	var regexpErr error
	if contentType == "ID" {
		REGEXP = ID_REGEXP
		regexpErr = ErrIDNotValid
	} else if contentType == "Username" {
		REGEXP = USERNAME_REGEXP
		regexpErr = ErrUsernameNotValid
	} else if contentType == "Comment" {
		REGEXP = COMMENT_REGEXP
		regexpErr = ErrCommentNotValid
	} else if contentType == "Datetime" {
		REGEXP = DATETIME_REGEXP
		regexpErr = ErrDatetimeNotValid
	} else {
		REGEXP = DATE_REGEXP
		regexpErr = ErrDateNotValid
	}

	regex, err := regexp.Compile(REGEXP)
	if err != nil {
		return err
	}

	if !regex.MatchString(content) {
		return regexpErr
	}

	return nil
}
