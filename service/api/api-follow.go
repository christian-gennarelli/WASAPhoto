package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-sqlite3"
)

func (rt _router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	followerUsername := ps.ByName("username")
	err := components.CheckIfValid(followerUsername, "Username")
	if err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided username not valid")
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

	// Check if the username from the path is the same as the authenticated one
	if *authUsername != followerUsername {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot follow an user on behalf of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot follow an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the following username from the path and check if it is valid
	followedUsername := ps.ByName("followed_username")
	if err = components.CheckIfValid(followedUsername, "Username"); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided username not valid")
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

	// Check if the user is trying to follow itself
	if followerUsername == followedUsername {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot auto-follow")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusForbidden, "cannot-autofollow").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if the authenticated user banned the user is trying to follow or viceversa
	err = rt.db.CheckIfBanned(followerUsername, followedUsername)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("cannot follow a banned user or that has banned the authenticated user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot follow a banned user or that has banned the authenticated user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the authenticated user banned the other user or viceversa")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the authenticated user banned the other user or viceversa").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Add the authenticated username to the list of users following the username provided in the path
	err = rt.db.FollowUser(followerUsername, followedUsername)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr); sqliteErr.Code == sqlite3.ErrConstraint {
			// if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("impossible to follow a non-existing user")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "impossible to follow a non-existing user").Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response")
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while following an user")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while following an user").Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response")
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt _router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Retrieve the follower username from the path and check if it's valid
	followerUsername := ps.ByName("username")
	err := components.CheckIfValid(followerUsername, "Username")
	if err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided username not valid")
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

	// Check if the authenticated user and the follower username are the same
	if *authUsername != followerUsername {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot unfollow an user on behalf of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "cannot unfollow an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the username of the following user and check if it's valid
	followedUsername := ps.ByName("followed_username")
	if err = components.CheckIfValid(followedUsername, "Username"); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("provided username not valid")
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

	// No need of ban checks

	// Add the authenticated username to the list of users following the username provided in the path
	if err = rt.db.UnfollowUser(followerUsername, followedUsername); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while unfollowing an user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while unfollowing an user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
