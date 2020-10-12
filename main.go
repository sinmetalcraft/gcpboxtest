package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
	metadatabox "github.com/sinmetalcraft/gcpbox/metadata"
	gcpboxtestCloudtasks "github.com/sinmetalcraft/gcpboxtest/backend/cloudtasks"
	"github.com/sinmetalcraft/gcpboxtest/backend/storage"
)

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	projectID, err := metadatabox.ProjectID()
	if err != nil {
		panic(err)
	}
	gaeService, err := metadatabox.AppEngineService()
	if err != nil {
		panic(err)
	}

	cloudtasksClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	taskboxService, err := taskbox.NewService(ctx, cloudtasksClient)
	if err != nil {
		panic(err)
	}

	cloudtasksHandlers, err := gcpboxtestCloudtasks.NewHandlers(ctx, projectID, gaeService, taskboxService)
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on port %s", port)
	http.HandleFunc("/storage/pubsubnotify", storage.StoragePubSubNotifyHandler)
	http.HandleFunc("/cloudtasks/appengine/json-post-task", cloudtasksHandlers.AppEnginePushHandler)
	http.HandleFunc("/cloudtasks/appengine/get-task", cloudtasksHandlers.AppEnginePushHandler)
	http.HandleFunc("/cloudtasks/appengine/add-task", cloudtasksHandlers.AddTask)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.DefaultServeMux); err != nil {
		log.Printf("failed ListenAndServe err=%+v", err)
	}
}
