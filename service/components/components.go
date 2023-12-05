// Here (almost) all the schemas defined in the API are defined.

package components

import (
	"fmt"
	"regexp"
	"time"
)

type ID struct {
	RandID string
}

type Username struct {
	Uname string `json:"name"`
}

type User struct {
	UID       ID
	UName     Username
	Name      string
	BirthDate time.Time
}

type Profile struct {
	User  User
	Posts []ID
}

/*
	Note: for arrays, a list of IDs is returned, not of objects.
	Their information will be retrieved later one if needed through their IDs.
	This reasoning is applied to UsersList, CommentsList and Stream.
*/

type UserList struct {
	Users []Username
}

type Photo struct {
	PhotoString string
}

type Post struct {
	PostID           ID
	Author           ID
	Photo            Photo
	CreationDatetime time.Time
	Description      string
}

type Stream struct {
	Posts []ID
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
	ErrorCode   string
	Description string
}

// Check if the provided username is in the correct format
func (Username Username) CheckIfValid() (*bool, error) {

	regex, err := regexp.Compile(USERNAME_REGEXP)
	if err != nil {
		return nil, fmt.Errorf("error encountered while compiling the regexp")
	}

	valid := true
	if !regex.MatchString(Username.Uname) {
		valid = false
	}

	return &valid, nil

}
