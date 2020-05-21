package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sinmetal/gcpboxtest/backend/cloudtasks"
	"github.com/sinmetal/gcpboxtest/backend/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	http.HandleFunc("/storage/pubsubnotify", storage.StoragePubSubNotifyHandler)
	http.HandleFunc("/cloudtasks/appengine/json-post-task", cloudtasks.AppEnginePushHandler)
	http.HandleFunc("/cloudtasks/appengine/get-task", cloudtasks.AppEnginePushHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.DefaultServeMux); err != nil {
		log.Printf("failed ListenAndServe err=%+v", err)
	}
}
