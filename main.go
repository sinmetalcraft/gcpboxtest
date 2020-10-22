package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/compute/metadata"
	"github.com/sinmetalcraft/gcpbox"
	cloudrunbox "github.com/sinmetalcraft/gcpbox/cloudrun"
	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
	metadatabox "github.com/sinmetalcraft/gcpbox/metadata"
	gaemetadatabox "github.com/sinmetalcraft/gcpbox/metadata/appengine"
	gcpboxtestCloudtasks "github.com/sinmetalcraft/gcpboxtest/backend/cloudtasks"
	gcpboxtestCloudtasksAppEngine "github.com/sinmetalcraft/gcpboxtest/backend/cloudtasks/appengine"
	gcpboxtestCloudtasksRun "github.com/sinmetalcraft/gcpboxtest/backend/cloudtasks/run"
	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/sinmetalcraft/gcpboxtest/backend/storage"
)

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Printf("Defaulting to port %s", port)
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

	runApiContainer, err := cloudrunbox.NewPrimitiveAPIContainer(ctx, gcpbox.TokyoRegion)
	if err != nil {
		panic(err)
	}
	cloudrunboxAdminService, err := cloudrunbox.NewAdminService(ctx, runApiContainer)
	if err != nil {
		panic(err)
	}
	gcpboxtestRunService, err := cloudrunboxAdminService.GetRunService(ctx, projectID, "gcpboxtest")
	if err != nil {
		panic(err)
	}

	onGAE := true
	gaeService, err := gaemetadatabox.Service()
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

	cloudtasksHandlers, err := gcpboxtestCloudtasks.NewHandlers(ctx, &gcpboxtestCloudtasks.HandlersConfig{
		ProjectID:                 projectID,
		ProjectNumber:             projectNumber,
		ServiceAccountEmail:       serviceAccountEmail,
		TargetGAEServiceID:        gaeService,
		TaskboxService:            taskboxService,
		CloudtasksClient:          cloudtasksClient,
		GcpboxtestCloudRunService: gcpboxtestRunService,
	})
	if err != nil {
		panic(err)
	}

	if onGAE {
		handlers, err := gcpboxtestCloudtasksAppEngine.NewHandlers(ctx, projectID, projectNumber, gaeService, taskboxService)
		if err != nil {
			panic(err)
		}

		http.HandleFunc(gcpboxtestCloudtasksAppEngine.AppEngineTasksHandlerUri, handlers.AppEngineTasksHandler)
		http.HandleFunc(gcpboxtestCloudtasksAppEngine.HttpTargetTasksHandlerUri, handlers.HttpTargetTasksHandler)
	} else {
		handlers, err := gcpboxtestCloudtasksRun.NewHandlers(ctx, projectID, projectNumber, cloudtasksClient, gcpboxtestRunService)
		if err != nil {
			panic(err)
		}

		http.HandleFunc("/cloudtasks/run/json-post-task", handlers.TasksHandler)
	}

	fmt.Printf("Listening on port %s", port)
	http.HandleFunc("/storage/pubsubnotify", storage.StoragePubSubNotifyHandler)
	http.HandleFunc("/cloudtasks/appengine/add-task", cloudtasksHandlers.AddTask)
	http.HandleFunc("/", HelloHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.DefaultServeMux); err != nil {
		fmt.Printf("failed ListenAndServe err=%+v", err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "request.header", r.Header)

	_, err := w.Write([]byte(fmt.Sprintf("Hello GCPBOXTEST : %s", time.Now().String())))
	if err != nil {
		fmt.Println(err.Error())
	}
}
