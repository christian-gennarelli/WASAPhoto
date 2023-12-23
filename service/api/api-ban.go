package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) getBanUserList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Retrieve the username of the authenticated user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	bannerUsername := components.Username{Value: ps.ByName("username")}
	err := bannerUsername.CheckIfValid()
	if err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
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

	// Check if the username in the path is the same as the authenticated one
	if bannerUsername.Value != username.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("cannot see the banned list of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot see the banned list of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Get the list of banned users
	bannedUserList, err := rt.db.GetBanUserList(bannerUsername.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while getting the banlist for the user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the banlist for the user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	response, err := json.MarshalIndent(bannedUserList.Users, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while encoding the response as JSON")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while encoding the response as JSON").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	if len(bannedUserList.Users) > 0 {
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

func (rt _router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username from the path and check if it's valid
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	bannerUsername := components.Username{Value: ps.ByName("username")}
	err := bannerUsername.CheckIfValid()
	if err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
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

	// Check if the username in the path is the same as the authenticated one
	if bannerUsername.Value != username.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("cannot ban an user on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot ban an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the username to be added to the list of banned user of the authenticated user from the query
	bannedUsername := components.Username{Value: r.URL.Query().Get("banned_username")}
	if err = bannedUsername.CheckIfValid(); err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
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
	err = rt.db.BanUser(bannerUsername.Value, bannedUsername.Value)
	if err != nil {
		var mess []byte
		if err == components.ErrForeignKeyConstraint {
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
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	bannerUsername := components.Username{Value: ps.ByName("username")}
	err := bannerUsername.CheckIfValid()
	if err != nil {
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

	// Check if the username in the path is the same as the authenticated one
	if bannerUsername.Value != username.Value {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("cannot unban an user on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot unban an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the username to be remove from the list of banned user of the authenticated user from the path
	bannedUsername := components.Username{Value: ps.ByName("banned_username")}
	err = bannedUsername.CheckIfValid()
	if err != nil {
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

	// Unban the user
	err = rt.db.UnbanUser(bannerUsername.Value, bannedUsername.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while unbanning the user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while unbanning the user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
