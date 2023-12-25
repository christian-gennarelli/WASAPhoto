// Here (almost) all the schemas defined in the API are defined.

package components

import (
	"regexp"
	"time"
)

type ID struct {
	Value string
}

type Username struct {
	Value string `json:"name"`
}

type User struct {
	ID         string
	Username   string
	Birthdate  string
	Name       string
	ProfilePic string // Base64 encoded image
}

type Profile struct {
	User  User
	Posts []Post
}

/*
	Note: for arrays, a list of IDs is returned, not of objects.
	Their information will be retrieved later one if needed through their IDs.
	This reasoning is applied to UsersList, CommentsList and Stream.
*/

type UserList struct {
	Users []User
}

type Post struct {
	PostID           ID
	Author           Username
	Photo            string // URL path to the image, stored server-side
	CreationDatetime string
	Description      string
}

type Stream struct {
	Posts []Post
}

type Comment struct {
	PostID           ID
	Body             string
	CreationDatetime time.Time
	Author           Username
}

type CommentList struct {
	Comments []ID
}

type Error struct {
	ErrorCode   int
	Description string
}

// Check if the provided username is in the correct format
func (Username Username) CheckIfValid() error {

	regex, err := regexp.Compile(USERNAME_REGEXP)
	if err != nil {
		return err
	}

	if !regex.MatchString(Username.Value) {
		return ErrIDNotValid
	}

	return nil
}

func (Id ID) CheckIfValid() error {

	regex, err := regexp.Compile(ID_REGEXP)
	if err != nil {
		return err
	}

	if !regex.MatchString(Id.Value) {
		return ErrUsernameNotValid
	}

	return nil

}

func (comment Comment) CheckIfValid() error {

	regex, err := regexp.Compile(COMMENT_REGEXP)
	if err != nil {
		return err
	}

	if !regex.MatchString(comment.Body) {
		return ErrCommentNotValid
	}

	return nil

}
