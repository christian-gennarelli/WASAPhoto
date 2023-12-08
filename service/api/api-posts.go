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

	// Retrieve the Auth token and check if is valid
	token := components.ID{RandID: r.Header.Get("Authorization")}
	valid, err := token.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the provided Auth token is valid or not"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
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

	if *valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided Auth token is not valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "The provided Auth token does not satisfy its associated regular expression",
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

	// Retrieve the username associated to the given Auth token
	username, err := rt.db.GetUsernameByToken(token.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while getting the username associated with the given token from the DB"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "There was an error while trying to fetch the username associated with the given token from the DB",
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

	// Check if there exists an user registered with such token
	if len(username.Uname) == 0 {
		w.WriteHeader(http.StatusBadRequest)

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "no username associated with the provided Auth token.",
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

	// Retrieve the id of the post the user wants to like and the username of the user who should own it
	ownerUsername := components.Username{Uname: ps.ByName("username")}
	postID := components.ID{RandID: ps.ByName("post_id")}

	valid, err = rt.db.CheckIfOwnerPost(ownerUsername.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post is owned by the provided user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
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

	if !*valid {
		w.WriteHeader(http.StatusBadRequest)

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
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

	// Add the username of the authenticated user to the list of likes of the post
	err = rt.db.AddLikeToPost(username.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while adding the like to the post"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
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

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Retrieve the Auth token and check if is valid
	token := components.ID{RandID: r.Header.Get("Authorization")}
	valid, err := token.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the provided Auth token is valid or not"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
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

	if *valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided Auth token is not valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "provided Auth token does not satisfy its associated regular expression",
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

	// Retrieve the username associated to the given Auth token
	username, err := rt.db.GetUsernameByToken(token.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while getting the username associated with the given token from the DB"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while trying to fetch the username associated with the given token from the DB",
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

	// Check if there exists an user registered with such token
	if len(username.Uname) == 0 {
		w.WriteHeader(http.StatusBadRequest)

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "no username associated with the provided Auth token",
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

	// Retrieve the id of the post the user wants to like and the username of the user who should own it
	ownerUsername := components.Username{Uname: ps.ByName("username")}
	postID := components.ID{RandID: ps.ByName("post_id")}

	valid, err = rt.db.CheckIfOwnerPost(ownerUsername.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post is owned by the provided user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
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

	if !*valid {
		w.WriteHeader(http.StatusBadRequest)

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
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
	if liker_username.Uname != username.Uname {
		w.WriteHeader(http.StatusBadRequest)

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided username does not coincide with the liker_username provided in the path",
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

	// Remove the like from the post
	err = rt.db.RemoveLikeFromPost(liker_username.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while removing the like from the post"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
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

}

func (rt _router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Retrieve the Auth token and check if is valid
	token := components.ID{RandID: r.Header.Get("Authorization")}
	valid, err := token.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the provided Auth token is valid or not"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
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

	if *valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided Auth token is not valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "provided Auth token does not satisfy its associated regular expression",
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

	// Retrieve the username associated to the given Auth token
	username, err := rt.db.GetUsernameByToken(token.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while getting the username associated with the given token from the DB"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while trying to fetch the username associated with the given token from the DB",
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

	// Check if there exists an user registered with such token
	if len(username.Uname) == 0 {
		w.WriteHeader(http.StatusBadRequest)

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "no username associated with the provided Auth token.",
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

	// Retrieve the id of the post the user wants to like and the username of the user who should own it
	ownerUsername := components.Username{Uname: ps.ByName("username")}
	postID := components.ID{RandID: ps.ByName("post_id")}

	valid, err = rt.db.CheckIfOwnerPost(ownerUsername.Uname, postID.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while checking if the given post is owned by the provided user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
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

	if !*valid {
		w.WriteHeader(http.StatusBadRequest)

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "The provided username does not own the given post",
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
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
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
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)

}
