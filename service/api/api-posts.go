package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the Auth token and check if is valid
	token := components.ID{RandID: r.Header.Get("Authorization")}

	valid, err := token.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the provided Auth token is valid or not"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the provided Auth token is valid or not",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
	}

	if *valid { // Auth token not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("provided Auth token not valid"))

		error, err := json.Marshal(components.Error{
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

		return
	}

	// Retrieve the username associated to the given Auth token and check if there exists an user registered with such token
	username, err := rt.db.GetUsernameByToken(token.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while getting the username associated with the given token from the DB"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while getting the username associated with the given token from the DB",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if len(username.Uname) == 0 { // No username associated with the provided token
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("no username associated with the provided Auth token"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "no username associated with the provided Auth token",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Retrieve the id of the post the user wants to like and check if it exists
	postID := components.ID{RandID: ps.ByName("post_id")}

	valid, err = rt.db.CheckIfPostExists(postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while checking if the given post exists",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Post does not exist
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided post does not exist"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided post does not exist",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Retrieve the username from the path and check if it is valid
	ownerUsername := components.Username{Uname: ps.ByName("username")}

	valid, err = ownerUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the given username is valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the given username is valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Username not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username not valid"))

		error, err := json.Marshal(components.Error{
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

		return
	}

	// Check if the username in the path is the owner of the given post
	valid, err = rt.db.CheckIfOwnerPost(ownerUsername.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post is owned by the provided user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while checking if the given post is owned by the provided user",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Provided username not owner of the given post
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

		return
	}

	// Add the username of the authenticated user to the list of likes of the post
	err = rt.db.AddLikeToPost(username.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while adding the like to the post"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while adding the like to the post",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the Auth token and check if is valid
	token := components.ID{RandID: r.Header.Get("Authorization")}

	valid, err := token.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the provided Auth token is valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the provided Auth token is valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if *valid { // Auth token not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided Auth token not valid"))

		error, err := json.Marshal(components.Error{
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

		return
	}

	// Retrieve the username associated to the given Auth token and check if there exists an user registered with such token
	username, err := rt.db.GetUsernameByToken(token.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while getting the username associated with the given token from the DB"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while getting the username associated with the given token from the DB",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if len(username.Uname) == 0 { // No username associated with the provided Auth token
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("no username associated with the provided Auth token"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "no username associated with the provided Auth token",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Retrieve the id of the post the user wants to like and check if it is valid
	postID := components.ID{RandID: ps.ByName("post_id")}

	valid, err = rt.db.CheckIfPostExists(postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while checking if the given post exists",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Post does not exist
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided post does not exist"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided post does not exist",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Retrieve the username from the path and check if it is valid
	ownerUsername := components.Username{Uname: ps.ByName("username")}

	valid, err = ownerUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the given username is valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the given username is valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Username not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username not valid"))

		error, err := json.Marshal(components.Error{
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

		return
	}

	// Check if the username from the path owns the given post
	valid, err = rt.db.CheckIfOwnerPost(ownerUsername.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post is owned by the provided user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while checking if the given post is owned by the provided user",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid {
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.Error(fmt.Errorf("provided username does not own the given post"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "404",
			Description: "provided username does not own the given post",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}

		return
	}

	// Check if the username associated with the Auth token and the liker_username provided in the path are the same
	liker_username := components.Username{Uname: ps.ByName("liker_username")}

	valid, err = liker_username.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the given username is valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the given username is valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // liker_username not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username not valid"))

		error, err := json.Marshal(components.Error{
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

		return
	}

	if liker_username.Uname != username.Uname {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("provided username does not coincide with the liker_username provided in the path"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided username does not coincide with the liker_username provided in the path",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Remove the like from the post
	err = rt.db.RemoveLikeFromPost(liker_username.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while removing the like from the post"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while removing the like from the post",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the Auth token and check if is valid
	token := components.ID{RandID: r.Header.Get("Authorization")}

	valid, err := token.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the provided Auth token is valid or not"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the provided Auth token is valid or not",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
	}

	if *valid { // Auth token not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("provided Auth token not valid"))

		error, err := json.Marshal(components.Error{
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

		return
	}

	// Retrieve the username associated to the given Auth token and check if there exists an user registered with such token
	username, err := rt.db.GetUsernameByToken(token.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while getting the username associated with the given token from the DB"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while getting the username associated with the given token from the DB",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if len(username.Uname) == 0 { // No username associated with the provided token
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("no username associated with the provided Auth token"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "no username associated with the provided Auth token",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Retrieve the id of the post the user wants to like and check if it exists
	postID := components.ID{RandID: ps.ByName("post_id")}

	valid, err = rt.db.CheckIfPostExists(postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while checking if the given post exists",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Post does not exist
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided post does not exist"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided post does not exist",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Retrieve the username from the path and check if it is valid
	ownerUsername := components.Username{Uname: ps.ByName("username")}

	valid, err = ownerUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the given username is valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the given username is valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Username not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username not valid"))

		error, err := json.Marshal(components.Error{
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

		return
	}

	// Check if the username in the path is the owner of the given post
	valid, err = rt.db.CheckIfOwnerPost(ownerUsername.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post is owned by the provided user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while checking if the given post is owned by the provided user",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Provided username not owner of the given post
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

		return
	}

	// Retrieve the comment from the request body
	var comment components.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while decoding the comment from the request body"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while decoding the comment from the request body",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Add the comment to the post
	err = rt.db.AddCommentToPost(postID.RandID, comment.Body, comment.CreationDatetime.Format("2006-01-02T15:04:05"), username.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while adding the comment to the post"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while adding the comment to the post",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the Auth token and check if is valid
	token := components.ID{RandID: r.Header.Get("Authorization")}

	valid, err := token.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the provided Auth token is valid or not"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the provided Auth token is valid or not",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
	}

	if *valid { // Auth token not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("provided Auth token not valid"))

		error, err := json.Marshal(components.Error{
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

		return
	}

	// Retrieve the username associated to the given Auth token and check if there exists an user registered with such token
	username, err := rt.db.GetUsernameByToken(token.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while getting the username associated with the given token from the DB"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while getting the username associated with the given token from the DB",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if len(username.Uname) == 0 { // No username associated with the provided token
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("no username associated with the provided Auth token"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "no username associated with the provided Auth token",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Retrieve the id of the post the user wants to like and check if it exists
	postID := components.ID{RandID: ps.ByName("post_id")}

	valid, err = rt.db.CheckIfPostExists(postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while checking if the given post exists",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Post does not exist
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided post does not exist"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided post does not exist",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Retrieve the username from the path and check if it is valid
	ownerUsername := components.Username{Uname: ps.ByName("username")}

	valid, err = ownerUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the given username is valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the given username is valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Username not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username not valid"))

		error, err := json.Marshal(components.Error{
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

		return
	}

	// Check if the username in the path is the owner of the given post
	valid, err = rt.db.CheckIfOwnerPost(ownerUsername.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post is owned by the provided user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while checking if the given post is owned by the provided user",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Provided username not owner of the given post
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

		return
	}

	// Retrieve the id of the comment from the path and check if it is valid
	commentID := components.ID{RandID: ps.ByName("comment_id")}

	valid, err = commentID.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the given comment id is valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while checking if the given comment id is valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Comment not valid
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided comment id not valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided comment id not valid",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Check if the given comment exists
	valid, err = rt.db.CheckIfCommentExists(commentID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given comment exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while checking if the given comment exists",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if !*valid { // Comment does not exist
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

		return
	}

	// Retrieve the owner of the comment, and check if the authenticated user is the owner of the comment
	ownerUsernameComment, err := rt.db.GetOwnerUsernameOfComment(commentID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while retrieving the owner of the comment from the DB"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while retrieving the owner of the comment from the DB",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	if ownerUsernameComment.Uname != username.Uname { // Authenticated user not owner of the comment
		w.WriteHeader(http.StatusBadRequest)

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided username does not own the given comment",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	// Delete the comment under the given post
	err = rt.db.RemoveCommentFromPost(postID.RandID, commentID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while removing the comment from the post"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while removing the comment from the post",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)

}
