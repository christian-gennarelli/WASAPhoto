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
	searchedUsername := components.Username{Value: r.URL.Query().Get("searched-username")}

	// Check if the provided username is valid
	err := searchedUsername.CheckIfValid()
	if err != nil {
		var mess []byte
		if err == components.ErrUsernameNotValid {
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

	// Search the users
	userList, err := rt.db.SearchUser(searchedUsername.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while searching the users with the provided username as substring")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while searching the users with the provided username as substring" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
	}

	response, err := json.MarshalIndent(userList.Users, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while encoding the response body as JSON")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while encoding the response body as JSON" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	if len(userList.Users) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	if _, err = w.Write(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while writing the response body in the response body")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while writing the response body in the response body" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

}
