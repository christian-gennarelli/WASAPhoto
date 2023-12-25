package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) getFollowingList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Obtain the username from the path and check if it's valid
	username := components.Username{Value: ps.ByName("username")}
	err := username.CheckIfValid()
	if err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
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

	// Check if the authenticated user banned the user of viceversa
	err = rt.db.CheckIfBanned(authUsername.Value, username.Value)
	if err == nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot get the following list of a banned user or that has banned the authenticated user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusForbidden, "cannot get the following list of a banned user or that has banned the authenticated user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the authenticated user banned the other user or viceversa")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while checking if the authenticated user banned the other user or viceversa").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	startDatetime := r.URL.Query().Get("datetime")
	if len(startDatetime) == 0 {
		t := time.Now()
		startDatetime = strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + " " + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
	}

	// Send the request to the database
	users, err := rt.db.GetFollowingList(username.Value, startDatetime)
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

	// Retrieve the username of the authenticated user
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the username from the path and check if it's valid
	followerUsername := components.Username{Value: ps.ByName("username")}
	err := followerUsername.CheckIfValid()
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

	// Check if the username from the path is the same as the authenticated one
	if username.Value != followerUsername.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot follow an user on behalf of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot follow an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the following username from the path and check if it is valid
	followedUsername := components.Username{Value: r.URL.Query().Get("followed_username")}
	if err = followedUsername.CheckIfValid(); err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
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
	if username.Value == followedUsername.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot auto-follow")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusForbidden, "cannot-autofollow").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if the authenticated user banned the user is trying to follow or viceversa
	err = rt.db.CheckIfBanned(username.Value, followedUsername.Value)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("cannot follow a banned user or that has banned the authenticated user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot follow a banned user or that has banned the authenticated user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the authenticated user banned the other user or viceversa")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while checking if the authenticated user banned the other user or viceversa").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Add the authenticated username to the list of users following the username provided in the path
	err = rt.db.FollowUser(username.Value, followedUsername.Value)
	if err != nil {
		if err == components.ErrForeignKeyConstraint {
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
	username := helperAuth(w, r, ps, ctx, rt)
	if username == nil {
		return
	}

	// Retrieve the follower username from the path and check if it's valid
	followerUsername := components.Username{Value: ps.ByName("username")}
	err := followerUsername.CheckIfValid()
	if err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
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

	// Check if the authenticated user and the follower username are the same
	if username.Value != followerUsername.Value {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot unfollow an user on behalf of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "cannot unfollow an user on behalf of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the username of the following user and check if it's valid
	followedUsername := components.Username{Value: ps.ByName("followed_username")}
	err = followedUsername.CheckIfValid()
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

	// No need of ban checks

	// Add the authenticated username to the list of users following the username provided in the path
	err = rt.db.UnfollowUser(followerUsername.Value, followedUsername.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while unfollowing an user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while unfollowing an user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (rt _router) getFollowersList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Obtain the username from the path and check if it's valid
	username := components.Username{Value: ps.ByName("username")}
	err := username.CheckIfValid()
	if err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
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

	// Check if the authenticated user banned the user of viceversa
	err = rt.db.CheckIfBanned(authUsername.Value, username.Value)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("cannot get the followers list of a banned user or that has banned the authenticated user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "cannot get the followers list of an user or that has banned the authenticated user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the authenticated user banned the other user or viceversa")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while checking if the authenticated user banned the other user or viceversa").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	startDatetime := r.URL.Query().Get("datetime")
	if len(startDatetime) == 0 {
		t := time.Now()
		startDatetime = strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + " " + strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()) + ":" + strconv.Itoa(t.Second())
	}

	// Send the request to the database
	users, err := rt.db.GetFollowersList(username.Value, startDatetime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while getting the list of followers")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the list of followers").Error())); err != nil {
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
