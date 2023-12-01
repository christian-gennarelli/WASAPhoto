package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parse the username of the user is trying to login
	var username components.Username
	err := json.NewDecoder(r.Body).Decode(&username)

	w.Header().Set("Content-Type", "application/json")

	// HTTP Error 400: Unacceptable
	if err != nil {

		// Write the header for the response code
		w.WriteHeader(http.StatusNotAcceptable)

		// Return the error to the logger
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while parsing the username from the request body"))

		// Return the error in the response body
		error := components.Error{
			ErrorCode:   400,
			Description: "Error while parsing the username from the request body...",
		}
		response, err := json.Marshal(error)
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			return
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}

		return
	}

	// Get the ID from the database
	id, err := rt.db.PostUserID(username.Uname)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while parsing the id for the given user"))
		return
	}

	ID := components.ID{
		RandID: id,
	}

	response, err := json.Marshal(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))

}
