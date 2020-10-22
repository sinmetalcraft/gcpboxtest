package cloudtasks

import (
	"context"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/sinmetalcraft/gcpbox/cloudrun"
	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
)

type Handlers struct {
	projectID                 string
	projectNumber             string
	serviceAccountEmail       string
	targetGAEServiceID        string
	taskboxService            *taskbox.Service
	cloudtasksClient          *cloudtasks.Client
	gcpboxtestCloudRunService *cloudrun.RunService
}

type HandlersConfig struct {
	ProjectID                 string
	ProjectNumber             string
	ServiceAccountEmail       string
	TargetGAEServiceID        string
	TaskboxService            *taskbox.Service
	CloudtasksClient          *cloudtasks.Client
	GcpboxtestCloudRunService *cloudrun.RunService
}

func NewHandlers(ctx context.Context, config *HandlersConfig) (*Handlers, error) {
	return &Handlers{
		projectID:                 config.ProjectID,
		projectNumber:             config.ProjectNumber,
		serviceAccountEmail:       config.ServiceAccountEmail,
		targetGAEServiceID:        config.TargetGAEServiceID,
		taskboxService:            config.TaskboxService,
		cloudtasksClient:          config.CloudtasksClient,
		gcpboxtestCloudRunService: config.GcpboxtestCloudRunService,
	}, nil
}
