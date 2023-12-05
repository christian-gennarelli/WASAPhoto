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

	Username := components.Username{Uname: r.Header.Get("username")}

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

	if !*valid { // Username not valid - wrong format
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username is not valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "The provided username does not satisfy its associated regular expression",
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

	// Parse the string we want to match in usernames
	searchedUsername := r.URL.Query().Get("searched-username")

	// Execute the query to the database
	res, err := rt.db.SearchUser(searchedUsername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: err.Error(),
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while formatting the error in JSON"))
			return
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}
		return
	}

	response, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := components.Error{
			ErrorCode:   "500",
			Description: "Error while econding the response body as JSON",
		}
		response, err := json.Marshal(error)
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response error as JSON"))
			return
		}
		_, err = w.Write([]byte(response))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			return
		}
		return
	}

	if len(response) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	_, err = w.Write([]byte(response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "Error while writing the response body in the response body",
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

func (rt _router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	Username := components.Username{Uname: ps.ByName("username")}

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

	if !*valid { // Username not valid - wrong format
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username is not valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "The provided username does not satisfy its associated regular expression",
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

	// Check if the username actually exists or not in WASAPhoto
	valid, err = rt.db.CheckIfUsernameExists(Username.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the username provided exists or not"))

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

	if !*valid { // The provided username does not exist
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username does not exists"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "404",
			Description: "The provided username is not registered on WASAPhoto",
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

	profile, err := rt.db.GetUserProfile(Username.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while getting the profile of the user"))

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

	response, err := json.Marshal(profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "Error while enconding the response as JSON",
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

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "Error while writing the response in the response body",
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
