package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Parse the string we want to match in usernames
	var searchedUsername components.Username

	err := json.NewDecoder(r.Body).Decode(&searchedUsername)
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

	// Check if the provided username is valid
	valid := searchedUsername.CheckIfValid(w, r, ps, ctx)
	if !*valid {
		return
	}

	// Search the users
	res, err := rt.db.SearchUser(searchedUsername.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while searching the users with the provided username as substring"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while searching the users with the provided username as substring",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while formatting the error in JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
		return
	}

	response, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response body as JSON"))

		error := components.Error{
			ErrorCode:   "500",
			Description: "error while encoding the response body as JSON",
		}
		response, err := json.Marshal(error)
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response error as JSON"))
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
		return
	}

	if len(res.Users) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	_, err = w.Write([]byte(response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response body in the response body"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while writing the response body in the response body",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
		return
	}

}
