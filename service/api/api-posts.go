package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the owner of the post and its ID
	ownerUsername, postID := helperPost(w, r, ps, ctx, rt)

	// Check if the authenticated user has banned the owner of the post or viceversa
	err := rt.db.CheckIfBanned(username.Value, ownerUsername.Value)
	if err == nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot like a photo of a banned user or that has banned the authenticated user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusForbidden, "cannot like a photo of an user or that has banned the authenticated user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the authenticated user banned the other user or viceversa")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while checking if the authenticated user banned the other user or viceversa").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Add the username of the authenticated user to the list of likes of the post
	err = rt.db.AddLikeToPost(username.Value, postID.Value)
	if err != nil {
		var mess []byte
		if err == components.ErrForeignKeyConstraint {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("username or post does NOT exist")
			mess = []byte(fmt.Errorf(components.StatusBadRequest, "username or post does NOT exist").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error encountered while adding the like to the post")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error encountered while adding the like to the post" /*err*/).Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the username of the owner of the post and its ID
	_, postID := helperPost(w, r, ps, ctx, rt)

	// No need of ban checks

	// Retrieve the username from the path and check if it is valid
	liker_username := components.Username{Value: ps.ByName("liker_username")}
	err := liker_username.CheckIfValid()
	if err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided username not valid")
			mess = []byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the username is valid")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the username is valid" /*err*/).Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if the authenticated user is the same as the liker username provided in the path
	if liker_username.Value != username.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("authenticated user cannot like another photo on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot like another photo on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Remove the like from the post
	err = rt.db.RemoveLikeFromPost(liker_username.Value, postID.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("error encountered while removing the like to the post")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error encountered while removing the like to the post" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt _router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Check if the content of the request body is in a JSON format
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		ctx.Logger.Error("unsupported media type provided")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusUnsupportedMediaType).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the username of the authenticated user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the username of the owner of the post and the ID of the post
	ownerUsername, postID := helperPost(w, r, ps, ctx, rt)

	// Check if the authenticated user banned the owner of the post or viceversa
	err := rt.db.CheckIfBanned(username.Value, ownerUsername.Value)
	if err == nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot comment a photo of a banned user or that has banned the authenticated user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot comment a photo of an user or that has banned the authenticated user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the authenticated user banned the other user or viceversa")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while checking if the authenticated user banned the other user or viceversa").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the comment from the request body
	var commentBody string
	if err := json.NewDecoder(r.Body).Decode(&commentBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while decoding the comment from the request body")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while decoding the comment from the request body" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Add the comment to the post
	err = rt.db.AddCommentToPost(postID.Value, commentBody, time.Now().Format("2006-01-02T15:04:05"), username.Value)
	if err != nil {
		var mess []byte
		if err == components.ErrForeignKeyConstraint {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("username or post does NOT exist")
			mess = []byte(fmt.Errorf(components.StatusNotFound, "username or post does NOT exist").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while adding the comment to the post")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while adding the comment to the post").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the username of the owner of the post and its ID
	_, postID := helperPost(w, r, ps, ctx, rt)

	// Retrieve the id of the comment from the path and check if it is valid
	commentID := components.ID{Value: ps.ByName("comment_id")}
	err := commentID.CheckIfValid()
	if err != nil {
		var mess []byte
		if err == components.ErrIDNotValid {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided comment not valid")
			mess = []byte(fmt.Errorf(components.StatusBadRequest, "provided comment not valid").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the comment is valid")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the comment is valid" /*err*/).Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the owner of the comment, and check if the authenticated user is the owner of the comment
	ownerUsernameComment, err := rt.db.GetOwnerUsernameOfComment(commentID.Value)
	if err != nil {
		var mess []byte
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("provided comment does not exists")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "provided comment does not exists").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while getting the owner of the given comment")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the owner of the given comment" /*err*/).Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if uthenticated user is the real owner of the comment
	if ownerUsernameComment.Value != username.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("authenticated user cannot uncomment a photo on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot uncomment a photo on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Delete the comment under the given post
	err = rt.db.RemoveCommentFromPost(postID.Value, commentID.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while removing the comment from the post")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while removing the comment from the post" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
