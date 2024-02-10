package api

import (
	"errors"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-sqlite3"
)

func (rt _router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username from the path and check if it's valid
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	bannerUsername := ps.ByName("username")
	if err := components.CheckIfValid(bannerUsername, "Username"); err != nil {
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

	// Check if the username in the path is the same as the authenticated one
	if bannerUsername != *authUsername {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.Error("cannot ban an user on behalf of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot ban an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the username to be added to the list of banned user of the authenticated user from the query
	bannedUsername := ps.ByName("banned_username")
	if err := components.CheckIfValid(bannedUsername, "Username"); err != nil {
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

	// Ban the user
	if err := rt.db.BanUser(bannerUsername, bannedUsername); err != nil {
		var mess []byte
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr); sqliteErr.Code == sqlite3.ErrConstraint {
			// if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("provided username not found")
			mess = []byte(fmt.Errorf(components.StatusNotFound, "provided username not found").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while banning the user")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while banning the user").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Authenticate the user
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	bannerUsername := ps.ByName("username")
	err := components.CheckIfValid(bannerUsername, "Username")
	if err != nil {
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

	// Check if the username in the path is the same as the authenticated one
	if bannerUsername != *authUsername {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.Error("cannot unban an user on behalf of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot unban an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the username to be remove from the list of banned user of the authenticated user from the path
	bannedUsername := ps.ByName("banned_username")
	if err := components.CheckIfValid(bannedUsername, "Username"); err != nil {
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

	// Unban the user
	if err := rt.db.UnbanUser(bannerUsername, bannedUsername); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while unbanning the user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while unbanning the user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
