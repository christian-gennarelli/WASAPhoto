package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	username := components.Username{Uname: ps.ByName("username")}

	// Check if the provided username is valid
	valid, err := username.CheckIfValid()
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
	valid, err = rt.db.CheckIfUsernameExists(username.Uname)
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

	profile, err := rt.db.GetUserProfile(username.Uname)
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

func (rt _router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Retrieve the authorization token and username from the request
	token := components.ID{RandID: r.Header.Get("Authorization")}
	username := components.Username{Uname: ps.ByName("username")}

	// Check if the provided username is valid
	valid, err := username.CheckIfValid()
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
	valid, err = rt.db.CheckIfUsernameExists(username.Uname)
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

	// Check if the user sending the request is authenticated
	valid, err = rt.db.CheckCombinationIsValid(username.Uname, token.RandID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the user is authenticated or not"))

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

	if !*valid { // The user is not authenticated
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username is not authenticated"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "401",
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

	// Retrieve the new username from the request body
	var new_username components.Username
	err = json.NewDecoder(r.Body).Decode(&new_username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while decoding the body of the request"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "Error encountered while parsing the username from the request body",
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

	// Check if the provided username is valid or not
	valid, err = new_username.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while checking if the new username provided is valid or not"))

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

	if !*valid { // The provided username is not valid
		w.WriteHeader(http.StatusNotFound)
		ctx.Logger.WithError(err).Error(fmt.Errorf("provided username is not valid"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "The provided username is not in a valid format.",
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

	// Call to the relative DB function
	err = rt.db.UpdateUsername(new_username.Uname, username.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while updating the username"))

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

	response, err := json.Marshal(new_username)
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

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "Error while writing the response",
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
