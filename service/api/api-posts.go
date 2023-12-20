package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	_, postID := helperPost(w, r, ps, ctx, rt)

	// Add the username of the authenticated user to the list of likes of the post
	err := rt.db.AddLikeToPost(username.Value, postID.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("error encountered while adding the like to the post")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error encountered while adding the like to the post" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	_, postID := helperPost(w, r, ps, ctx, rt)

	// Check if the username associated with the Auth token and the liker_username provided in the path are the same
	liker_username := components.Username{Value: ps.ByName("liker_username")}
	valid, err := liker_username.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the username is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the username is valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}
	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("provided username not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}

	if liker_username.Value != username.Value {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("authenticated user cannot like another photo on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "authenticated user cannot like another photo on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}

	// Remove the like from the post
	err = rt.db.RemoveLikeFromPost(liker_username.Value, postID.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("error encountered while removing the like to the post")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error encountered while removing the like to the post" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt _router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		ctx.Logger.Error("unsupported media type provided")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "Invalid Content-Type. Only application/json is supported").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	_, postID := helperPost(w, r, ps, ctx, rt)

	// Retrieve the comment from the request body
	var comment components.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error encountered while decoding the comment from the request body")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error encountered while decoding the comment from the request body" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}

	// Add the comment to the post
	err := rt.db.AddCommentToPost(postID.Value, comment.Body, comment.CreationDatetime.Format("2006-01-02T15:04:05"), username.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while adding the comment to the post")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while adding the comment to the post" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the Auth token and check if is valid
	username := helperAuth(w, r, ps, ctx, rt)

	_, postID := helperPost(w, r, ps, ctx, rt)

	// Retrieve the id of the comment from the path and check if it is valid
	commentID := components.ID{Value: ps.ByName("comment_id")}
	err := commentID.CheckIfValid()
	if err != nil {
		if err.Error() == "id not valid" {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided comment not valid")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided comment not valid").Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the comment is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the comment is valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the owner of the comment, and check if the authenticated user is the owner of the comment
	ownerUsernameComment, err := rt.db.GetOwnerUsernameOfComment(commentID.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("provided comment does not exists")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "provided comment does not exists").Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response")
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while getting the owner of the given comment")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the owner of the given comment" /*err*/).Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
		}
		return
	}

	if ownerUsernameComment.Value != username.Value { // Authenticated user not owner of the comment
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("authenticated user cannot uncomment a photo on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "authenticated user cannot uncomment a photo on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}

	// Delete the comment under the given post
	err = rt.db.RemoveCommentFromPost(postID.Value, commentID.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while removing the comment from the post")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while removing the comment from the post" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
