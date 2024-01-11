package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func helperAuth(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext, rt _router) *components.Username {

	// Retrieve the Auth token and check if is valid
	token := components.ID{Value: r.Header.Get("Authorization")}
	if err := token.CheckIfValid(); err != nil {
		if errors.Is(err, components.ErrIDNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided auth token not valid")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided auth token not valid").Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
			return nil
		}
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the auth token is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the auth token is valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return nil
	}

	// Retrieve the username (if valid) associated to the given Auth token and check if there exists an user registered with such token
	var username *components.Username
	username, err := rt.db.GetUsernameByToken(token.Value)
	if err != nil {
		var mess []byte
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusUnauthorized)
			ctx.Logger.WithError(err).Error("no user found with the provided authenticated token")
			mess = []byte(fmt.Errorf(components.StatusUnauthorized, "no user found with the provided token").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while getting the username associated with the given token")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the username associated with the given token" /*err*/).Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return nil
	}

	return username // Return nil if not found

}

func helperPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext, rt _router, retrievePost bool) (*components.Username, *components.ID) {

	// Retrieve the username from the path and check if it is valid
	ownerUsername := components.Username{Value: ps.ByName("username")}
	if err := ownerUsername.CheckIfValid(); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
			w.WriteHeader(http.StatusUnauthorized)
			mess = []byte(fmt.Errorf(components.StatusUnauthorized, "provided username not valid").Error())
			ctx.Logger.Error("provided username not valid")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the username is valid")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the username is valid").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return nil, nil
	}

	if !retrievePost {
		return &ownerUsername, nil
	}

	// Retrieve the id of the post the user wants to like and check if it exists
	postID := components.ID{Value: ps.ByName("post_id")}
	if err := postID.CheckIfValid(); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrIDNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided post id not valid")
			mess = []byte(fmt.Errorf(components.StatusBadRequest, "provided post id not valid").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the post id is valid")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the post id is valid").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return nil, nil
	}

	// Check if the username in the path is the owner of the given post
	if err := rt.db.CheckIfOwnerPost(ownerUsername.Value, postID.Value); err != nil {
		var mess []byte
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("either the username or the post not found; alternatively, the username does not own the provided post")
			mess = []byte(fmt.Errorf(components.StatusNotFound, "either the username or the post not found; alternatively, the username does not own the provided post").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the given username owns the provided post")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the given username owns the provided post").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return nil, nil
	}

	return &ownerUsername, &postID

}

// func helperBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext, rt _router, username components.Username) *components.Username {

// 	// Retrieve the username from the path and check if it's valid
// 	bannerUsername := components.Username{Value: ps.ByName("username")}
// 	err := bannerUsername.CheckIfValid()
// 	if err != nil {
// 		var mess []byte
// 		if err == components.ErrUsernameNotValid {
// 			w.WriteHeader(http.StatusBadRequest)
// 			ctx.Logger.WithError(err).Error("provided username not valid")
// 			mess = []byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid").Error())
// 		} else {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			ctx.Logger.WithError(err).Error("error while checking if the username is valid")
// 			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the username is valid").Error())
// 		}
// 		if _, err = w.Write(mess); err != nil {
// 			ctx.Logger.WithError(err).Error("error while writing the response")
// 		}
// 		return nil
// 	}

// 	if bannerUsername.Value != username.Value {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		ctx.Logger.Error("cannot unban an user on behalf of another user")
// 		if _, err := w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot unban an user on behalf of another user").Error())); err != nil {
// 			ctx.Logger.WithError(err).Error("error while writing the response")
// 		}
// 		return nil
// 	}

// 	return &bannerUsername

// }
