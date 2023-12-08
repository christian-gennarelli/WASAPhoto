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
	var Username components.Username
	err := json.NewDecoder(r.Body).Decode(&Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while decoding the body of the request"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while parsing the username from the request body",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}

		return
	}

	// Check if the provided username is valid
	valid, err := Username.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while decoding the body of the request"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}

		return
	}

	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username is not valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "provided username does not satisfy its associated regular expression",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	// HTTP Error 400: Unacceptable
	if err != nil {

		// Write the header for the response code
		w.WriteHeader(http.StatusNotAcceptable)

		// Return the error to the logger
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while parsing the username from the request body"))

		// Return the error in the response body
		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "error while parsing the username from the request body...",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}

		return
	}

	// Get the ID from the database
	ID, err := rt.db.PostUserID(Username.Uname)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while parsing the id for the given user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while parsing the username from the DB...",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}

		return
	}

	response, err := json.Marshal(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while enconding the response body as JSON",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response error as JSON"))
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while writing the response body in the response body",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error as JSON"))
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}
		return
	}

}
