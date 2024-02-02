package api

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) getPhotoFromURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// ADD AUTHENTICATION

	// ADD BAN CHECK

	w.Header().Set("Content-Type", "application/json")

	// Retrieve the path of the photo
	path := r.URL.Query().Get("photo_path")

	// Open the image
	img, err := os.Open("photos/" + path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while opening the file specified by the given path")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while opening the file specified by the given path").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Read the image
	reader := bufio.NewReader(img)
	content, err := io.ReadAll(reader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while reading the content of the file specified by the given path")
		ctx.Logger.Info(path)
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while reading the content of the file specified by the given path").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Send the image to the client
	w.Header().Set("Content-Type", "image/png")
	if _, err = w.Write(content); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(("error while writing the response in the response body"))
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, err).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response in the response body")
		}
		return
	}

}
