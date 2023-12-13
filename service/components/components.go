// Here (almost) all the schemas defined in the API are defined.

package components

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
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
func (Username Username) CheckIfValid(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *bool {

	valid := false

	regex, err := regexp.Compile(USERNAME_REGEXP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while compiling the regex for checking the validity of the provided username"))

		error, err := json.Marshal(Error{
			ErrorCode:   "500",
			Description: "error while compiling the regex for checking the validity of the provided username",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return &valid
	}

	valid = regex.MatchString(Username.Uname)

	if !valid { // Username not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username not valid"))

		error, err := json.Marshal(Error{
			ErrorCode:   "400",
			Description: "provided username not valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return &valid
	}

	return &valid

}

func (Id ID) CheckIfValid(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *bool {

	valid := false

	regex, err := regexp.Compile(ID_REGEXP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while compiling the regex for checking the validity of the provided ID"))

		error, err := json.Marshal(Error{
			ErrorCode:   "500",
			Description: "error while compiling the regex for checking the validity of the provided ID",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return &valid
	}

	valid = regex.MatchString(Id.RandID)

	if !valid { // Auth token not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("provided Auth token not valid"))

		error, err := json.Marshal(Error{
			ErrorCode:   "400",
			Description: "provided Auth token not valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return &valid
	}

	return &valid

}
