package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) getFollowersList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Obtain the username from the path and check if it's valid
	username := components.Username{Value: ps.ByName("username")}
	valid, err := username.CheckIfValid()
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

	// Send the request to the database
	users, err := rt.db.GetFollowersList(username.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while getting the list of followers")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the list of followers" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	if len(users.Users) > 0 {
		w.WriteHeader(http.StatusOK)
		response, err := json.MarshalIndent(users.Users, "", " ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(("error while encoding the response as JSON"))
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, err).Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while encoding the response as JSON")
			}
			return
		}
		if _, err = w.Write(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(("error while writing the response in the response body"))
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, err).Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response in the response body")
			}
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (rt _router) getFollowingList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Obtain the username from the path and check if it's valid
	username := components.Username{Value: ps.ByName("username")}
	valid, err := username.CheckIfValid()
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

	// Send the request to the database
	users, err := rt.db.GetFollowingList(username.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while getting the list of followings")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the list of followings" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	if len(users.Users) > 0 {
		w.WriteHeader(http.StatusOK)
		response, err := json.MarshalIndent(users.Users, "", " ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(("error while encoding the response as JSON"))
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, err).Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while encoding the response as JSON")
			}
			return
		}
		if _, err = w.Write(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(("error while writing the response in the response body"))
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, err).Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response in the response body")
			}
			return
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}

func (rt _router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Authenticate the user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	followerUsername := components.Username{Value: ps.ByName("username")}
	valid, err := followerUsername.CheckIfValid()
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

	// Check if the username from the path is the same as the authenticated one
	if username.Value != ps.ByName("username") {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("cannot follow an user on behalf of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot follow an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the username from the path and check if it is valid
	var followingUsername components.Username
	err = json.NewDecoder(r.Body).Decode(&followingUsername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while decoding the body of the request to obtain the following username")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while decoding the body of the request to obtain the following username" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	valid, err = followingUsername.CheckIfValid()
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
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if the user is trying to follow itself
	if username.Value == followingUsername.Value {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("cannot auto-follow")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot-autofollow").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	valid, err = rt.db.CheckIfBanned(username.Value, followingUsername.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the following username has been banned")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while checking if the following username has been banned").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("cannot follow a banned user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot follow a banned user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	// Add the authenticated username to the list of users following the username provided in the path
	err = rt.db.FollowUser(username.Value, followingUsername.Value)
	if err != nil {
		if err.Error() == "FOREIGN KEY constraint failed" { // trying to follow a non-existing user
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("impossible to follow a non-existing user")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "impossible to follow a non-existing user" /*err*/).Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response")
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while following an user")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while following an user" /*err*/).Error())); err != nil {
				ctx.Logger.WithError(err).Error("error while writing the response")
			}
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt _router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Authenticate the user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	followerUsername := components.Username{Value: ps.ByName("username")}
	valid, err := followerUsername.CheckIfValid()
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

	if username.Value != followerUsername.Value {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("cannot unfollow an user on behalf of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot unfollow an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	followingUsername := components.Username{Value: ps.ByName("followed_username")}
	valid, err = followingUsername.CheckIfValid()
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

	// Add the authenticated username to the list of users following the username provided in the path
	err = rt.db.UnfollowUser(followerUsername.Value, followingUsername.Value)
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
