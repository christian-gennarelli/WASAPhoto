package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-sqlite3"
)

func (rt _router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username of the authenticated user
	authUsername := helperAuth(w, r, ps, ctx, rt)
	if authUsername == nil {
		return
	}

	// Retrieve the username from the path and check if it is valid
	username := ps.ByName("username")
	if err := components.CheckIfValid(username, "Username"); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided username not valid")
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

	// Check if the authenticated user banned the user provided in the path, or viceversa
	if err := rt.db.CheckIfBanned(*authUsername, username); err == nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("cannot get the profile a banned user or that has banned the authenticated user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusForbidden, "cannot get the profile of an user or that has banned the authenticated user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the authenticated user banned the other user or viceversa")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while checking if the authenticated user banned the other user or viceversa").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the profile of the user with the given username
	profile, err := rt.db.GetUserProfile(username)
	if err != nil {
		var mess []byte
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("provided username does not exist")
			mess = []byte(fmt.Errorf(components.StatusNotFound, "provided username does not exist").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while getting the profile of the user")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the profile of the user").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return
	}

	// Send the profile to the client
	response, err := json.MarshalIndent(profile, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while enconding the response as JSON")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, err).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while enconding the response as JSON")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(("error while writing the response in the response body"))
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, err).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response in the response body")
		}
		return
	}

}

func (rt _router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "text/plain")

	// Retrieve the authorization token and username from the request
	usernameAuth := helperAuth(w, r, ps, ctx, rt)
	if usernameAuth == nil {
		return
	}

	// Retrieve the username from the path and check if it is valid
	username := ps.ByName("username")

	if err := components.CheckIfValid(username, "Username"); err != nil {
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

	// Check if the two usernames coincide
	if username != *usernameAuth { // Not the same username
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Error("not authorized to change the username of another user")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusForbidden, "not authorized to change the username of another user").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Retrieve the new username and check if it's valid
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while decoding the body of the request")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while decoding the body of the request").Error())); err != nil {
			ctx.Logger.WithError(err).Error("erroe while writing the response")
		}
		return
	}
	newUsername := string(body)

	if err = components.CheckIfValid(newUsername, "Username"); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided new username not valid")
			mess = []byte(fmt.Errorf(components.StatusBadRequest, "provided new username not valid").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the new username is valid")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the new username is valid").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Update the username
	err = rt.db.UpdateUsername(newUsername, username)
	if err != nil {
		var mess []byte
		if errors.Is(err, sql.ErrNoRows) { // Old username not found
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("provided username does not exist")
			mess = []byte(fmt.Errorf(components.StatusNotFound, "provided username does not exists").Error())
			// else if errors.Is(err, sqlite3.ErrConstraintUnique) { // Username already exists
		} else if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			w.WriteHeader(http.StatusNotAcceptable)
			ctx.Logger.WithError(err).Error("provided username already exists")
			mess = []byte(fmt.Errorf(components.StatusNotAcceptable, "provided username already exists").Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while updating the username")
			mess = []byte(fmt.Errorf(components.StatusInternalServerError, "error while updating the username").Error())
		}
		if _, err = w.Write(mess); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Change photo profile name
	srcFolder := "photos/profile_pics/" + username + ".png"
	destFolder := "photos/profile_pics/" + newUsername + ".png"
	mvCmd := exec.Command("mv", srcFolder, destFolder)
	if err := mvCmd.Run(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while updating the user's profile pic")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while updating the user's profile pic").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Send the new username to the client as confirmation of the its new username
	response, err := json.MarshalIndent(newUsername, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while enconding the response as JSON")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while enconding the response as JSON").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	if _, err = w.Write(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while writing the response")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while writing the response").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

}
