package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func helperAuthFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext, rt _router) (*components.Username, *components.Username) {

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

		return nil, nil

	}

	if !*valid { // Auth token not valid
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

		return nil, nil
	}

	// Retrieve the username associated to the given Auth token and check if there exists an user registered with such token
	followerUsername, err := rt.db.GetUsernameByToken(token.RandID)
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

		return nil, nil
	}

	if len(followerUsername.Uname) == 0 { // No username associated with the provided token
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

		return nil, nil
	}

	// Retrieve the username from the path and check if it is valid
	followingUsername := components.Username{Uname: ps.ByName("username")}

	valid, err = followingUsername.CheckIfValid()
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

		return nil, nil
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

		return nil, nil
	}

	return followerUsername, &followingUsername

}

func (rt _router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	followerUsername, followingUsername := helperAuthFollowing(w, r, ps, ctx, rt)
	if followerUsername == nil || followingUsername == nil {
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

	followerUsername, followingUsername := helperAuthFollowing(w, r, ps, ctx, rt)
	if followerUsername == nil || followingUsername == nil {
		return
	}

	// Add the authenticated username to the list of users following the username provided in the path
	err := rt.db.UnfollowUser(followerUsername.Uname, followingUsername.Uname)
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
