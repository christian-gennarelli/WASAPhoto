package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
	ownerUsername, postID := helperPost(w, r, ps, ctx, rt, true)

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
	_, postID := helperPost(w, r, ps, ctx, rt, true)

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
	ownerUsername, postID := helperPost(w, r, ps, ctx, rt, true)

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
	var comment components.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while decoding the comment from the request body")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while decoding the comment from the request body").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	if err := comment.CheckIfValid(); err != nil {
		var mess []byte
		if err == components.ErrCommentNotValid {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided comment not valid")
			mess = []byte(fmt.Errorf(components.StatusBadRequest, "provided comment not valid").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the comment is valid")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the comment is valid").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	t := time.Now()
	currentDatetime := strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + "T" + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
	// Add the comment to the post
	err = rt.db.AddCommentToPost(postID.Value, comment.Body, currentDatetime, username.Value)
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
	_, postID := helperPost(w, r, ps, ctx, rt, true)

	// Retrieve the id of the comment from the path and check if it is valid
	commentID := components.ID{Value: ps.ByName("comment_id")}
	err := commentID.CheckIfValid()
	if err != nil {
		var mess []byte
		if err == components.ErrIDNotValid {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided comment ID not valid")
			mess = []byte(fmt.Errorf(components.StatusBadRequest, "provided comment ID not valid").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the comment ID is valid")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the comment ID is valid" /*err*/).Error())
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

func (rt _router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Retrieve the username of the authenticated user
	usernameAuth := helperAuth(w, r, ps, ctx, rt)
	if usernameAuth == nil {
		return
	}

	usernameOwner, _ := helperPost(w, r, ps, ctx, rt, false)

	// Check if the username in the path and the authenticated one are the same
	if usernameOwner.Value != usernameAuth.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("authenticated user cannot post a post on the profile of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot post a post on the profile of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the post from the request body
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while decoding the body of the request")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while decoding the body of the request").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	// Access the request body
	formData := r.MultipartForm

	// Accessing the photo file
	photo := formData.File["photo"][0]
	// Open the file
	fileReader, err := photo.Open()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while opening the photo sent in the request")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while opening the photo sent in the request").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if the file sent by the client is actually an image (NOT WORKING - DON'T KNOW WHY)
	// if _, _, err := image.DecodeConfig(fileReader); err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	ctx.Logger.WithError(err).Error("provided file is not an image")
	// 	if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "provided file is not an image").Error())); err != nil {
	// 		ctx.Logger.WithError(err).Error("error while writing the response")
	// 	}
	// 	return
	// }

	defer fileReader.Close()

	// Accessing the description field
	description := formData.Value["description"][0]
	if err := (components.Comment{Body: description}).CheckIfValid(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("provided comment not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided comment not valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	err, post := rt.db.UploadPost(usernameOwner.Value, description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while posting the photo")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while posting the photo").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	// Save the file locally
	uploadedFile, err := os.Create(post.Photo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while loading the photo locally")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while loading the photo locally").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		if _, err := rt.db.DeletePost(post.PostID.Value); err != nil {
			ctx.Logger.WithError(err).Error("error while deleting the record just uploaded")
		}
		return
	}
	defer uploadedFile.Close()
	// Copy the file content to the local file
	_, err = io.Copy(uploadedFile, fileReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while copying the photo locally")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while copying creating the photo locally").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		if _, err := rt.db.DeletePost(post.PostID.Value); err != nil {
			ctx.Logger.WithError(err).Error("error while deleting the record just uploaded")
		}
		return
	}

}

func (rt _router) deletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Retrieve the username of the authenticated user
	usernameAuth := helperAuth(w, r, ps, ctx, rt)
	if usernameAuth == nil {
		return
	}

	ownerUsername, postID := helperPost(w, r, ps, ctx, rt, true)

	// Check if the username in the path and the authenticated one are the same
	if ownerUsername.Value != usernameAuth.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("authenticated user cannot delete a post on the profile of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot post a post on the profile of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Delete the record representing the post from the database
	photoPath, err := rt.db.DeletePost(postID.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while deleting the post from the database")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while deleting the post from the database").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Delete the file
	if err := os.Remove(*photoPath); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while deleting the photo frrom the server")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while deleting the photo frrom the server").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

}

func (rt _router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	usernameAuth := helperAuth(w, r, ps, ctx, rt)
	if usernameAuth == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	username := components.Username{Value: ps.ByName("username")}
	if err := username.CheckIfValid(); err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided username not valid")
			mess = []byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the username is valid")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the username is valid").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check that the username from the path and the authenticated username is the same
	if usernameAuth.Value != username.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("authenticated user cannot see the stream of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot see the stream of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the starting datetime (the last 16 posts between the provided datetime and the current one will be returned)
	startDatetime := r.URL.Query().Get("start-datetime")
	if len(startDatetime) == 0 {
		t := time.Now()
		startDatetime = strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + "T" + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
	}

	// Retrieve the stream of the user
	postStream, err := rt.db.GetUserStream(startDatetime, username.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while retrieving the stream for the given user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while retrieving the stream for the given user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Encode the response as JSON
	response, err := json.MarshalIndent(postStream.Posts, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while encoding the response as JSON")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while encoding the response as JSON").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Send the response to the client, if not empty
	if len(postStream.Posts) > 0 {
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while writing the response")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while writing the response").Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response")
			}
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}
