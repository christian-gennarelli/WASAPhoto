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

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the username and check if it is valid
	username := components.Username{Uname: ps.ByName("username")}

	valid := username.CheckIfValid(w, r, ps, ctx)
	if !*valid {
		return
	}

	// Check if the username actually exists or not in WASAPhoto
	valid = rt.db.CheckIfUsernameExists(username.Uname, w, r, ps, ctx)
	if !*valid {
		return
	}

	// Retrieve the profile of the user with the given username
	profile := rt.db.GetUserProfile(username.Uname, w, r, ps, ctx)
	if profile == nil {
		return
	}

	// Send the profile to the client
	response, err := json.Marshal(profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while enconding the response as JSON"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while enconding the response as JSON",
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

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(response))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response in the response body"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while writing the response in the response body",
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

func (rt _router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the authorization token and username from the request
	token := components.ID{RandID: r.Header.Get("Authorization")}
	username := components.Username{Uname: ps.ByName("username")}

	// Check if the provided username is valid
	valid := username.CheckIfValid(w, r, ps, ctx)
	if !*valid {
		return
	}

	// Check if the provided Auth token is valid
	valid = token.CheckIfValid(w, r, ps, ctx)
	if !*valid {
		return
	}

	// Retrieve the username associated to the given Auth token
	usernameAuth := rt.db.GetUsernameByToken(token.RandID, w, r, ps, ctx)
	if usernameAuth == nil {
		return
	}

	if username.Uname != usernameAuth.Uname { // Not the same username
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.Error(fmt.Errorf("not authorized to change the username of another user"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "401",
			Description: "not authorized to change the username of another user",
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

	// Retrieve the new username from the request body
	var new_username components.Username
	err := json.NewDecoder(r.Body).Decode(&new_username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while decoding the body of the request"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while parsing the username from the request body",
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

	// Check if the provided username is valid or not
	valid = new_username.CheckIfValid(w, r, ps, ctx)
	if !*valid {
		return
	}

	// Update the username
	err = rt.db.UpdateUsername(new_username.Uname, username.Uname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while updating the username"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while updating the username",
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

	// Send the new username to the client as confirmation of the its new username
	response, err := json.Marshal(new_username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while enconding the response as JSON"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while enconding the response as JSON",
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
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while writing the response",
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
