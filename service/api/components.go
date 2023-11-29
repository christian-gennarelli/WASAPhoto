// Here (almost) all the schemas defined in the API are defined.

package api

import (
	"time"
)

type ID struct {
	randID string
}

type Username struct {
	uname string
}

type User struct {
	UID       ID
	UName     Username
	Name      string
	BirthDate time.Time
}

/*
	Note: for arrays, a list of IDs is returned, not of objects.
	Their information will be retrieved later one if needed through their IDs.
	This reasoning is applied to UsersList, CommentsList and Stream.
*/

type UserList struct {
	Users []ID
}

type Photo struct {
	PhotoString string
}

type Post struct {
	PostID           ID
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
	ErrorCode   int
	Description string
}
