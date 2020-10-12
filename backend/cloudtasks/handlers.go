package cloudtasks

import (
	"context"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
)

type Handlers struct {
	projectID           string
	projectNumber       string
	serviceAccountEmail string
	targetGAEServiceID  string
	taskboxService      *taskbox.Service
	cloudtasksClient    *cloudtasks.Client
}

func NewHandlers(ctx context.Context, projectID string, projectNumber string, serviceAccountEmail string, targetGAEServiceID string, taskboxService *taskbox.Service, cloudtasksClient *cloudtasks.Client) (*Handlers, error) {
	return &Handlers{
		projectID:           projectID,
		projectNumber:       projectNumber,
		serviceAccountEmail: serviceAccountEmail,
		targetGAEServiceID:  targetGAEServiceID,
		taskboxService:      taskboxService,
		cloudtasksClient:    cloudtasksClient,
	}, nil
}
