package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func helperAuth(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext, rt _router) *components.Username {

	// Retrieve the Auth token and check if is valid
	token := components.ID{RandID: r.Header.Get("Authorization")}
	valid, err := token.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the username is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the username is valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return nil
	}

	// Retrieve the username (if valid) associated to the given Auth token and check if there exists an user registered with such token
	var username *components.Username
	if *valid {
		username, err = rt.db.GetUsernameByToken(token.RandID)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				ctx.Logger.WithError(err).Error("error while checking if the username is valid")
				if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the username is valid" /*err*/).Error())); err != nil {
					ctx.Logger.WithError(err).Error("errow while writing the response")
				}
				return nil
			}
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while getting the username associated with the given token")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while getting the username associated with the given token" /*err*/).Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
			return nil
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("provided token not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided token not valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return nil
	}

	return username // Return nil if not found

}

func helperPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext, rt _router) (*components.Username, *components.ID) {
	// Retrieve the id of the post the user wants to like and check if it exists
	postID := components.ID{RandID: ps.ByName("post_id")}
	valid, err := postID.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the post is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the post is valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return nil, nil
	}
	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("provided post not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided post not valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return nil, nil
	}

	if err = rt.db.CheckIfPostExists(postID.RandID); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("provided post does not exist")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusNotFound, "provided post does not exist").Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("error while checking if the post exists")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the post exists" /*err*/).Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
		}
		return nil, nil
	}

	// Retrieve the username from the path and check if it is valid
	ownerUsername := components.Username{Uname: ps.ByName("username")}
	valid, err = ownerUsername.CheckIfValid()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the username is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the username is valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return nil, nil
	}
	if !*valid {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("provided username not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("errow while writing the response")
		}
		return nil, nil
	}

	// Check if the username in the path is the owner of the given post
	err = rt.db.CheckIfOwnerPost(ownerUsername.Uname, postID.RandID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("provided username does not own the given post")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusNotFound, "provided username does not own the given post").Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided username not valid")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid").Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
		}
		return nil, nil
	}

	return &ownerUsername, &postID

}
