package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	c "simple-video-server/controllers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/uploadVideo", c.UploadVideo).Methods("POST")
	err := http.ListenAndServe(":4000", handlers.LoggingHandler(os.Stdout, router))
	if err != nil {
		log.Fatal(err)
	}
}
