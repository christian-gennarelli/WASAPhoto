package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func helperAuth(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext, rt _router) *components.Username {

	// Retrieve the Auth token and check if is valid
	token := components.ID{RandID: r.Header.Get("Authorization")}

	valid := token.CheckIfValid(w, r, ps, ctx)
	if !*valid {
		return nil
	}

	// Retrieve the username associated to the given Auth token and check if there exists an user registered with such token
	username := rt.db.GetUsernameByToken(token.RandID, w, r, ps, ctx)

	return username // nil if not found

}

func (rt _router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	followerUsername := helperAuth(w, r, ps, ctx, rt)
	if followerUsername == nil {
		return
	}

	// Retrieve the username from the path and check if it is valid
	followingUsername := components.Username{Uname: ps.ByName("username")}
	valid := followingUsername.CheckIfValid(w, r, ps, ctx)
	if !*valid {
		return
	}

	// Add the authenticated username to the list of users following the username provided in the path
	err := rt.db.FollowUser(followerUsername.Uname, followingUsername.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while following the username"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while querying the database to add followerUsername to the followers of followingUsername",
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

func (rt _router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	followerUsername := helperAuth(w, r, ps, ctx, rt)
	if followerUsername == nil {
		return
	}

	// Retrieve the username from the path and check if it is valid
	var followingUsername components.Username
	err := json.NewDecoder(r.Body).Decode(&followingUsername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while decoding the body of the request"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while decoding the body of the request",
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

	valid := followingUsername.CheckIfValid(w, r, ps, ctx)
	if !*valid {
		return
	}

	// Add the authenticated username to the list of users following the username provided in the path
	err = rt.db.UnfollowUser(followerUsername.Uname, followingUsername.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while unfollowing the username"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while querying the database to remove followerUsername frmo the followers of followingUsername",
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
