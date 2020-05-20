package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.DefaultServeMux); err != nil {
		log.Printf("failed ListenAndServe err=%+v", err)
	}
}
