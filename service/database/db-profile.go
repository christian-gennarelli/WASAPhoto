package database

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

// Retrieve the profile of the user with the provided username
func (db appdbimpl) GetUserProfile(Username string, w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *components.Profile {

	// Retrieve the photos posted by this user
	stmt, err := db.c.Prepare("SELECT PostID FROM Post WHERE Author = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while preparing the SQL statement to obtain the list of posts posted by the user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while preparing the SQL statement to obtain the list of posts posted by the user",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
		return nil
	}

	postIDs, err := stmt.Query(Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while performing the query to obtain the list of posts posted by the user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while performing the query to obtain the list of posts posted by the user",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
		return nil
	} else {
		defer postIDs.Close()
	}

	var Posts []components.ID
	for postIDs.Next() {
		var postID components.ID
		err = postIDs.Scan(&postID)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while extracting the username from the query"))

			error, err := json.Marshal(components.Error{
				ErrorCode:   "500",
				Description: "error while extracting the username from the query",
			})
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			}
			_, err = w.Write([]byte(error))
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			}
			return nil
		}

		Posts = append(Posts, postID)
	}

	// Retrieve the informations about the user with the provided username
	stmt, err = db.c.Prepare("SELECT * FROM User WHERE Username = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while preparing the SQL statement to obtain the info about the user with the provided username"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while preparing the SQL statement to obtain the info about the user with the provided username",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
		return nil
	}

	users, err := stmt.Query(Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while executing the SQL statement to obtain the info about the user with the provided username"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while executing the SQL statement to obtain the info about the user with the provided username",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
		return nil
	} else {
		defer postIDs.Close()
	}

	var User components.User
	for users.Next() {
		err = postIDs.Scan(&User)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while extracting the username from the query"))

			error, err := json.Marshal(components.Error{
				ErrorCode:   "500",
				Description: "error while extracting the username from the query",
			})
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			}
			_, err = w.Write([]byte(error))
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			}
			return nil
		}
	}

	profile := components.Profile{
		User:  User,
		Posts: Posts,
	}

	return &profile

}
