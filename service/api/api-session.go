package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Parse the username of the user is trying to login
	if contentType := r.Header.Get("Content-Type"); contentType != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		ctx.Logger.Error("unsupported media type provided")
		if _, err := w.Write([]byte(fmt.Errorf(components.StatusUnsupportedMediaType).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while reading the body of the request")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while reading the body of the request" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	username := string(body)

	// Check if the provided username is valid
	if err := components.CheckIfValid(username, "Username"); err != nil {
		var mess []byte
		if errors.Is(err, components.ErrUsernameNotValid) {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided username not valid: " + username)
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

	// Get the ID from the database
	user, err := rt.db.PostUserID(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while parsing the id for the given user")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while parsing the id for the given user" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	response, err := json.MarshalIndent(&user, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while enconding the response body as JSON")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while enconding the response body as JSON" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Send the response to the client
	w.WriteHeader(http.StatusCreated)

	if _, err = w.Write(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while writing the response body in the response body")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while writing the response body in the response body" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
}
