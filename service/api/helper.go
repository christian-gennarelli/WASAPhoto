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
	token := components.ID{Value: r.Header.Get("Authorization")}
	err := token.CheckIfValid()
	if err != nil {
		if err.Error() == "id not valid" {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided comment not valid")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided comment not valid").Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
			return nil
		}
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the comment is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the comment is valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return nil
	}

	// Retrieve the username (if valid) associated to the given Auth token and check if there exists an user registered with such token
	var username *components.Username
	username, err = rt.db.GetUsernameByToken(token.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("no user found with such token")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusNotFound, "no user found with such token" /*err*/).Error())); err != nil {
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

	return username // Return nil if not found

}

func helperPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext, rt _router) (*components.Username, *components.ID) {
	// Retrieve the id of the post the user wants to like and check if it exists
	postID := components.ID{Value: ps.ByName("post_id")}
	err := postID.CheckIfValid()
	if err != nil {
		if err.Error() == "id not valid" {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("provided comment not valid")
			if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided comment not valid").Error())); err != nil {
				ctx.Logger.WithError(err).Error("errow while writing the response")
			}
			return nil, nil
		}
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the comment is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the comment is valid" /*err*/).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return nil, nil
	}

	if err = rt.db.CheckIfPostExists(postID.Value); err != nil {
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
	ownerUsername := components.Username{Value: ps.ByName("username")}
	valid, err := ownerUsername.CheckIfValid()
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
	err = rt.db.CheckIfOwnerPost(ownerUsername.Value, postID.Value)
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

// func helperBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext, rt _router) *components.Username {

// 	// Authenticate the user
// 	username := helperAuth(w, r, ps, ctx, rt)
// 	if username == nil {
// 		return nil
// 	}

// 	// Retrieve the username from the path and check if it's valid
// 	bannerUsername := components.Username{Value: ps.ByName("username")}
// 	valid, err := bannerUsername.CheckIfValid()
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		ctx.Logger.WithError(err).Error("error while retrieving the username from the path")
// 		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "provided username not valid" /*err*/).Error())); err != nil {
// 			ctx.Logger.WithError(err).Error("error while writing the response")
// 		}
// 		return nil
// 	}
// 	if !*valid {
// 		w.WriteHeader(http.StatusBadRequest)
// 		ctx.Logger.Error("provided username not valid")
// 		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided username not valid" /*err*/).Error())); err != nil {
// 			ctx.Logger.WithError(err).Error("error while writing the response")
// 		}
// 		return nil
// 	}

// 	// Check if the username in the path is the same as the authenticated one
// 	if bannerUsername.Value != username.Value {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		ctx.Logger.WithError(err).Error("cannot ban an user on behalf of another user")
// 		if _, err = w.Write([]byte(fmt.Errorf(components.StatusUnauthorized, "cannot ban an user on behalf of another user" /*err*/).Error())); err != nil {
// 			ctx.Logger.WithError(err).Error("error while writing the response")
// 		}
// 		return nil
// 	}

// 	return &bannerUsername

// }
