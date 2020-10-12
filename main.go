package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/compute/metadata"
	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
	metadatabox "github.com/sinmetalcraft/gcpbox/metadata"
	gcpboxtestCloudtasks "github.com/sinmetalcraft/gcpboxtest/backend/cloudtasks"
	gcpboxtestCloudtasksAppEngine "github.com/sinmetalcraft/gcpboxtest/backend/cloudtasks/appengine"
	gcpboxtestCloudtasksRun "github.com/sinmetalcraft/gcpboxtest/backend/cloudtasks/run"
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
	projectNumber, err := metadata.NumericProjectID()
	if err != nil {
		panic(err)
	}
	serviceAccountEmail, err := metadatabox.ServiceAccountEmail()
	if err != nil {
		panic(err)
	}

	onGAE := true
	gaeService, err := metadatabox.AppEngineService()
	if errors.Is(err, metadatabox.ErrNotFound) {
		onGAE = false
	} else if err != nil {
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

	cloudtasksHandlers, err := gcpboxtestCloudtasks.NewHandlers(ctx, projectID, projectNumber, serviceAccountEmail, gaeService, taskboxService, cloudtasksClient)
	if err != nil {
		panic(err)
	}

	if onGAE {
		handlers, err := gcpboxtestCloudtasksAppEngine.NewHandlers(ctx, projectID, projectNumber, gaeService, taskboxService)
		if err != nil {
			panic(err)
		}

		http.HandleFunc("/cloudtasks/appengine/json-post-task", handlers.TasksHandler)
	} else {
		handlers, err := gcpboxtestCloudtasksRun.NewHandlers(ctx, projectID, projectNumber, cloudtasksClient)
		if err != nil {
			panic(err)
		}

		http.HandleFunc("/cloudtasks/run/json-post-task", handlers.TasksHandler)
	}

	log.Printf("Listening on port %s", port)
	http.HandleFunc("/storage/pubsubnotify", storage.StoragePubSubNotifyHandler)
	http.HandleFunc("/cloudtasks/appengine/add-task", cloudtasksHandlers.AddTask)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.DefaultServeMux); err != nil {
		log.Printf("failed ListenAndServe err=%+v", err)
	}
}
