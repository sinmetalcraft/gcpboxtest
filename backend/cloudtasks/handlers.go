package cloudtasks

import (
	"context"

	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
)

type Handlers struct {
	projectID          string
	projectNumber      string
	targetGAEServiceID string
	taskboxService     *taskbox.Service
}

func NewHandlers(ctx context.Context, projectID string, projectNumber string, targetGAEServiceID string, taskboxService *taskbox.Service) (*Handlers, error) {
	return &Handlers{
		projectID:          projectID,
		projectNumber:      projectNumber,
		targetGAEServiceID: targetGAEServiceID,
		taskboxService:     taskboxService,
	}, nil
}
