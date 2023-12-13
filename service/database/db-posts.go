package database

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
)

func (db appdbimpl) CheckIfPostExists(PostID string, w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *bool {

	valid := false

	stmt, err := db.c.Prepare("SELECT 1 FROM Post WHERE PostID = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while preparing the query to check if the given post exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while preparing the query to check if the given post exists",
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

	rows, err := stmt.Query(PostID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while executing the query to check if the given post exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while executing the query to check if the given post exists",
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
		defer rows.Close()
	}

	if rows.Next() {
		valid = true
	}

	return &valid

}

func (db appdbimpl) CheckIfCommentExists(CommentID string, w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *bool {

	valid := false

	stmt, err := db.c.Prepare("SELECT 1 FROM Comment WHERE CommentID = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while preparing the SQL statement for checking if the provided comment exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while preparing the SQL statement for checking if the provided comment exists",
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

	rows, err := stmt.Query(CommentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while performing the SQL statement for checking if the provided comment exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while performing the SQL statement for checking if the provided comment exists",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return &valid
	} else {
		defer rows.Close()
	}

	if !rows.Next() {
		valid = true
	}

	if !valid { // Comment does not exist
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided comment does not exist"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided comment does not exist",
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

func (db appdbimpl) CheckIfOwnerPost(Username string, PostID string, w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *bool {

	valid := false

	stmt, err := db.c.Prepare("SELECT 1 FROM User U JOIN Post P ON U.Username = P.Author WHERE U.Username = ? AND P.PostID = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while preparing the SQL statement to check if the given post is owned by the given user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while preparing the SQL statement to check if the given post is owned by the given user",
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

	rows, err := stmt.Query(Username, PostID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while executing the SQL query to retrieve the username associated to the given Auth token"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while executing the SQL query to retrieve the username associated to the given Auth token",
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

	if !rows.Next() {
		valid = true
	}

	if !valid { // Provided username not owner of the given post
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username does not own the given post"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "404",
			Description: "provided username does not own the given post",
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

func (db appdbimpl) AddLikeToPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("INSERT INTO Like (PostID, Liker) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to add the like")
	}

	_, err = stmt.Query(PostID, Username)
	if err != nil {
		return fmt.Errorf("error while executing the query to add the like")
	}

	return nil

}

func (db appdbimpl) RemoveLikeFromPost(Username string, PostID string) error {

	stmt, err := db.c.Prepare("DELETE FROM Like WHERE PostID = ? AND Liker = ?  ")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to add the like")
	}

	_, err = stmt.Query(PostID, Username)
	if err != nil {
		return fmt.Errorf("error while executing the query to add the like")
	}

	return nil

}

func (db appdbimpl) AddCommentToPost(PostID string, Body string, CreationDatetime string, Author string) error {

	stmt, err := db.c.Prepare("INSERT INTO Comment (CommentID, PostID, Author, CreationDatetime, Comment) VALUES (?, ?, ?, CONVERT(DATETIME, ?), ?)")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to add the comment")
	}

	// Generate the comment id
	commentID := uniuri.NewLen(64)

	_, err = stmt.Exec(commentID, PostID, Author, CreationDatetime, Body)
	if err != nil {
		return fmt.Errorf("error while executing the query to add the comment")
	}

	return nil

}

func (db appdbimpl) RemoveCommentFromPost(PostID string, CommentID string) error {

	stmt, err := db.c.Prepare("DELETE FROM Comment WHERE PostID = ? AND CommentID = ?")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to remove the comment")
	}

	_, err = stmt.Exec(PostID, CommentID)
	if err != nil {
		return fmt.Errorf("error while executing the query to remove the comment")
	}

	return nil

}
