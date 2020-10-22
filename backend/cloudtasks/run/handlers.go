package run

import (
	"context"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/sinmetalcraft/gcpbox/cloudrun"
)

type Handlers struct {
	projectID                 string
	projectNumber             string
	targetServiceID           string
	taskboxService            *cloudtasks.Client
	gcpboxtestCloudRunService *cloudrun.RunService
}

func NewHandlers(ctx context.Context, projectID string, projectNumber string, tasksService *cloudtasks.Client, gcpboxtestCloudRunService *cloudrun.RunService) (*Handlers, error) {
	return &Handlers{
		projectID:                 projectID,
		projectNumber:             projectNumber,
		taskboxService:            tasksService,
		gcpboxtestCloudRunService: gcpboxtestCloudRunService,
	}, nil
}
