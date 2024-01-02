package api

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/julienschmidt/httprouter"
)

func (rt _router) getPhotoFromURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "image/png")

	// Retrieve the path of the photo
	path := ps.ByName("photo_path")
	regex, err := regexp.Compile("(profile_pics/[0-9]+.png)|(posts/[a-zA-Z][a-zA-Z_]{7,14}[a-zA-Z].png)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while checking if the path is valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, "error while checking if the path is valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Check if the path is valid
	if !regex.Match([]byte(path)) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("provided path not valid")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "provided path not valid").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Read the image
	img, err := os.Open("photos/" + path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while opening the file specified by the given path")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while opening the file specified by the given path").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}
	reader := bufio.NewReader(img)
	content, err := io.ReadAll(reader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error while reading the content of the file specified by the given path")
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusBadRequest, "error while reading the content of the file specified by the given path").Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response")
		}
		return
	}

	// Send the image to the client
	if _, err = w.Write(content); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(("error while writing the response in the response body"))
		if _, err = w.Write([]byte(fmt.Errorf(components.StatusInternalServerError, err).Error())); err != nil {
			ctx.Logger.WithError(err).Error("error while writing the response in the response body")
		}
		return
	}

}
