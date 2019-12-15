package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	u "simple-video-info/utils"
)

// 100 << 20 == 100Mb
const MaxMemory = 100 << 20

var UploadVideo = func(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(MaxMemory); err != nil {
		u.ReturnError(w, err, http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		u.ReturnError(w, err, http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(MaxMemory); err != nil {
		u.ReturnError(w, err, http.StatusBadRequest)
		return
	}
	path, err := u.UploadFileAndGetPath(file, handler)
	if err != nil {
		u.ReturnError(w, err, http.StatusInternalServerError)
		return
	}
	ctx, err := u.NewContext()
	if err != nil {
		errString := fmt.Sprintf("Failed to create context: %v\n", err)
		u.ReturnError(w, errors.New(errString), http.StatusInternalServerError)
		return
	}
	defer ctx.Free()

	u.OpenInput(ctx, filepath.Join(path, handler.Filename))
	respond := u.GetVideoInfo(ctx)
	u.Respond(w, respond)
}
