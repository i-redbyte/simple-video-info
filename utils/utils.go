package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func ErrorMessage(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

func ReturnError(w http.ResponseWriter, err error, code int) {
	log.Println(err)
	w.WriteHeader(code)
	Respond(w, ErrorMessage(err.Error()))
}

func UploadFile(file multipart.File, handler *multipart.FileHeader) error {
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	path := filepath.Join(currentDir, "videos")
	fmt.Println("path", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		er := os.MkdirAll(path, os.ModePerm)
		if er != nil {
			return err
		}
	}
	fileName := handler.Filename
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	newFile, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		err := newFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fileBytes, err := ioutil.ReadAll(file)
	_, err = newFile.Write(fileBytes)
	if err != nil {
		return err
	}
	return nil
}
