package api

import (
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Authenticate the user
	followerUsername := helperAuth(w, r, ps, ctx, rt)
	if followerUsername == nil {
		return
	}

	// Retrieve the username from the path and check if it is valid
	followingUsername := components.Username{Uname: ps.ByName("username")}
	valid, err := followingUsername.CheckIfValid()
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
		ctx.Logger.WithError(err).Error("provided username not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	if followerUsername.Uname != ps.ByName("username") {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("not authorized to unfollow an user on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "not authorized to follow an user on behalf of another user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
	}

	// Add the authenticated username to the list of users following the username provided in the path
	err = rt.db.FollowUser(followerUsername.Uname, followingUsername.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while following an user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while following an user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt _router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Authenticate the user
	followerUsername := helperAuth(w, r, ps, ctx, rt)
	if followerUsername == nil {
		return
	}

	// Retrieve the username from the path and check if it is valid
	//var followingUsername components.Username
	// err := json.NewDecoder(r.Body).Decode(&followingUsername)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	ctx.Logger.WithError(err).Error("error while decoding the body of the request to obtain the following username")
	// 	if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while decoding the body of the request to obtain the following username" /*err*/).Error())); err != nil {
	// 		ctx.Logger.WithError(err).Error("error while writing the response")
	// 	}
	// 	return
	// }

	followingUsername := components.Username{Uname: ps.ByName("username_followed")}

	valid, err := followingUsername.CheckIfValid()
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
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	if followerUsername.Uname != ps.ByName("username") {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("not authorized to unfollow an user on behalf of another user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "not authorized to unfollow an user on behalf of another user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
	}

	// Add the authenticated username to the list of users following the username provided in the path
	err = rt.db.UnfollowUser(followerUsername.Uname, followingUsername.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while unfollowing an user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while unfollowing an user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
