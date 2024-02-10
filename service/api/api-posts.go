package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"image"

	_ "image/jpeg" // Blank import for accepting jpeg images with the image package
	_ "image/png"  // Blank import for accepting png images with the image package
	"io"
	"net/http"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-sqlite3"
)

func (rt _router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Retrieve the owner of the post and its ID
	ownerUsername, postID := helperPost(w, r, ps, ctx, rt, true)
	if ownerUsername == nil || postID == nil {
		return
	}

	// Check if the authenticated user has banned the owner of the post or viceversa
	err := rt.db.CheckIfBanned(*authUsername, *ownerUsername)
	if err == nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot like a photo of a banned user or that has banned the authenticated user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusForbidden, "cannot like a photo of a banned user or that has banned the authenticated user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the authenticated user banned the other user or viceversa")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while checking if the authenticated user banned the other user or viceversa").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Add the username of the authenticated user to the list of likes of the post
	err = rt.db.AddLikeToPost(*authUsername, *postID)
	if err != nil {
		var mess []byte
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr); sqliteErr.Code == sqlite3.ErrConstraint {
			// if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("username or post does NOT exist")
			mess = []byte(fmt.Errorf(components.StatusNotFound, "username or post does NOT exist").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error encountered while adding the like to the post")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error encountered while adding the like to the post").Error())
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
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Retrieve the username of the owner of the post and its ID
	_, postID := helperPost(w, r, ps, ctx, rt, true)
	if postID == nil {
		return
	}

	// No need of ban checks

	// Retrieve the username from the path and check if it is valid
	likerUsername := ps.ByName("liker_username")
	if err := components.CheckIfValid(likerUsername, "Username"); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
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

	// Check if the authenticated user is the same as the liker username provided in the path
	if likerUsername != *authUsername {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("authenticated user cannot like another photo on behalf of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot like another photo on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Remove the like from the post
	if err := rt.db.RemoveLikeFromPost(likerUsername, *postID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("error encountered while removing the like to the post")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error encountered while removing the like to the post").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt _router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "text/plain")

	// Retrieve the username of the authenticated user
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Retrieve the username of the owner of the post and the ID of the post
	ownerUsername, postID := helperPost(w, r, ps, ctx, rt, true)
	if ownerUsername == nil || postID == nil {
		return
	}

	// Check if the authenticated user banned the owner of the post or viceversa
	err := rt.db.CheckIfBanned(*authUsername, *ownerUsername)
	if err == nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot comment a photo of a banned user or that has banned the authenticated user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot comment a photo of an user or that has banned the authenticated user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the authenticated user banned the other user or viceversa")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while checking if the authenticated user banned the other user or viceversa").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the comment from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while decoding the comment from the request body")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while decoding the comment from the request body").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	comment := string(body)

	if err := components.CheckIfValid(comment, "Comment"); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrCommentNotValid) {
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

	// Add the comment to the post
	commentPost, err := rt.db.AddCommentToPost(*postID, comment, *authUsername)
	if err != nil {
		var mess []byte
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr); sqliteErr.Code == sqlite3.ErrConstraint {
			// if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
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

	// Encode the response as JSON
	response, err := json.MarshalIndent(*commentPost, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while encoding the response as JSON")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while encoding the response as JSON").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Send the response to the client, if not empty
	w.WriteHeader(http.StatusAccepted)
	if _, err = w.Write(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while writing the response")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while writing the response").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

}

func (rt _router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Retrieve the username of the owner of the post and its ID
	_, postID := helperPost(w, r, ps, ctx, rt, true)
	if postID == nil {
		return
	}

	// Retrieve the id of the comment from the path and check if it is valid
	commentID := ps.ByName("comment_id")

	// Retrieve the owner of the comment, and check if the authenticated user is the owner of the comment
	ownerUsernameComment, err := rt.db.GetOwnerUsernameOfComment(commentID)
	if err != nil {
		var mess []byte
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("provided comment does not exist")
			mess = []byte(fmt.Errorf(components.StatusNotFound, "provided comment does not exist").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while getting the owner of the given comment")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the owner of the given comment").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if authenticated user is the real owner of the comment
	if *ownerUsernameComment != *authUsername {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("authenticated user cannot uncomment a photo on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot uncomment a photo on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Delete the comment under the given post
	err = rt.db.RemoveCommentFromPost(*postID, commentID)
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
	if usernameOwner == nil {
		return
	}

	// Check if the username in the path and the authenticated one are the same
	if *usernameOwner != *usernameAuth {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("authenticated user cannot post a post on the profile of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot post a post on the profile of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

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

	rawPhoto := formData.File["photo"]
	if len(rawPhoto) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no photo provided")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "no photo provided").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	photo := rawPhoto[0]

	// Access the photo file
	fileReader, err := photo.Open()
	if err != nil {
		http.Error(w, "Unable to open photo file", http.StatusInternalServerError)
		return
	}
	defer fileReader.Close()

	// Check if the provided file is an image (check the first 512 bytes to determine its Content-Type)
	buff := make([]byte, 512)
	if _, err = fileReader.Read(buff); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the provided file is an image")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the provided file is an image").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	ctx.Logger.Info(http.DetectContentType(buff))
	if http.DetectContentType(buff) != "image/png" && http.DetectContentType(buff) != "image/jpeg" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("provided file not an image")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "provided file not an image").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	_, err = fileReader.Seek(0, 0) // Move the byte reader back to the beginning of the file
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking the info about the photo (seeking)")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking the info about the photo (seeking)").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	conf, _, err := image.DecodeConfig(fileReader)
	if err != nil {
		if errors.Is(err, image.ErrFormat) { // err.Error() == "image: unknown format"
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided image not in a valid format")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided image not in a valid format").Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response")
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking the info about the photo")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking the info about the photo").Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response")
			}
		}
		return
	}
	// ctx.Logger.Info("Height:" + strconv.FormatInt(int64(conf.Height), 10))
	// ctx.Logger.Info("Width: " + strconv.FormatInt(int64(conf.Width), 10))

	// Check the size of the image: it must be 16:9
	if conf.Width/conf.Height != 16/9 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("photo does not satisfy size requirements: it must be 16:9")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "photo does not satisfy size requirements: it must be 16:9").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Accessing the description field
	rawDescription := formData.Value["description"]
	if len(rawDescription) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("no comment provided")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "no comment provided").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	description := rawDescription[0]

	if err := components.CheckIfValid(description, "Comment"); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrCommentNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided comment not valid")
			mess = []byte(fmt.Errorf(components.StatusBadRequest, "provided comment not valid").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the description is valid")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the description is valid").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	post, err := rt.db.UploadPost(*usernameOwner, description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while posting the photo")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while posting the photo").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	_, err = fileReader.Seek(0, 0) // Move the byte reader back to the beginning of the file
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error before creating the photo locally")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error before creating the photo locally").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Save the file locally
	uploadedFile, err := os.Create("photos/" + post.Photo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while creating the photo locally")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while creating the photo locally").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		if _, err := rt.db.DeletePost(post.PostID); err != nil {
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
		if _, err := rt.db.DeletePost(post.PostID); err != nil {
			ctx.Logger.WithError(err).Error("error while deleting the record just uploaded")
		}
		return
	}

	response, err := json.MarshalIndent(*post, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while encoding the response")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while encoding the response").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while writing the response")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while writing the response").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

}

func (rt _router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Retrieve the username of the authenticated user
	usernameAuth := helperAuth(w, r, ps, ctx, rt)
	if usernameAuth == nil {
		return
	}

	ownerUsername, postID := helperPost(w, r, ps, ctx, rt, true)
	if ownerUsername == nil || postID == nil {
		return
	}

	// Check if the username in the path and the authenticated one are the same
	if *ownerUsername != *usernameAuth {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("authenticated user cannot delete a post on the profile of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot post a post on the profile of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Delete the record representing the post from the database
	photoPath, err := rt.db.DeletePost(*postID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while deleting the post from the database")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while deleting the post from the database").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Delete the file
	if err := os.Remove("photos/" + *photoPath); err != nil {
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
	username := ps.ByName("username")
	if err := components.CheckIfValid(username, "Username"); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
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
	if *usernameAuth != username {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("authenticated user cannot see the stream of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "authenticated user cannot see the stream of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the stream of the user
	postStream, err := rt.db.GetUserStream(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while retrieving the stream for the given user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while retrieving the stream for the given user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Encode the response as JSON
	response, err := json.MarshalIndent(*postStream, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while encoding the response as JSON")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while encoding the response as JSON").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Send the response to the client, if not empty
	if len(*postStream) > 0 {
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
