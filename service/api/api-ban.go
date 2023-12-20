package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Authenticate the user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	bannerUsername := components.Username{Value: ps.ByName("username")}
	valid, err := bannerUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while retrieving the username from the path")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "provided username not valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("provided username not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if the username in the path is the same as the authenticated one
	if bannerUsername.Value != username.Value {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("cannot ban an user on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot ban an user on behalf of another user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	var bannedUsername components.Username
	// Retrieve the username to be added to the list of banned user of the authenticated user from the query
	err = json.NewDecoder(r.Body).Decode(&bannedUsername) /*r.URL.Query().Get("banned_username")*/
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while decoding the body of the request")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while decoding the body of the request" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	valid, err = bannedUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the username is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the username is valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("provided username not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Execute the request to the database
	err = rt.db.BanUser(bannerUsername.Value, bannedUsername.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while banning the user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while banning the user" /*err*/).Error())); err != nil {
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
	valid, err := bannerUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while retrieving the username from the path")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "provided username not valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("provided username not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if the username in the path is the same as the authenticated one
	if bannerUsername.Value != username.Value {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("cannot ban an user on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot ban an user on behalf of another user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the username to be remove from the list of banned user of the authenticated user from the path
	bannedUsername := components.Username{Value: ps.ByName("banned_username")}
	valid, err = bannedUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while retrieving the username from the path")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while retrieving the username from the path" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("provided username not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	err = rt.db.UnbanUser(bannerUsername.Value, bannedUsername.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while unbanning the user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while unbanning the user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) getBanUserList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Authenticate the user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	bannerUsername := components.Username{Value: ps.ByName("username")}
	valid, err := bannerUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while retrieving the username from the path")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "provided username not valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("provided username not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if the username in the path is the same as the authenticated one
	if bannerUsername.Value != username.Value {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("cannot ban an user on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot ban an user on behalf of another user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Execute the call to the database
	bannedUserList, err := rt.db.GetBanUserList(bannerUsername.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while getting the banlist for the user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the banlist for the user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	response, err := json.MarshalIndent(bannedUserList.Users, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while encoding the response as JSON")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while encoding the response as JSON" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	if len(bannedUserList.Users) > 0 {
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while writing the response")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while writing the response" /*err*/).Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response")
			}
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}
