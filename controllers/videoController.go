package controllers

import (
	"log"
	"net/http"
	u "simple-video-server/utils"
)

// 100 << 20 == 100Mb
const MaxMemory = 100 << 20

var UploadVideo = func(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(MaxMemory); err != nil {
		u.ReturnError(w, err, http.StatusBadRequest)
		log.Println(err)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		u.ReturnError(w, err, http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(MaxMemory); err != nil {
		log.Println(err)
		u.ReturnError(w, err, http.StatusBadRequest)
		return
	}
	err = u.UploadFile(file, handler)
	if err != nil {
		log.Println(err)
		u.ReturnError(w, err, http.StatusBadRequest)
		return
	}
	u.Respond(w, struct {
		Success bool `json:"success"`
	}{Success: true})
}
