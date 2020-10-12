package cloudtasks

import (
	"context"

	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
)

type Handlers struct {
	projectID       string
	targetServiceID string
	taskboxService  *taskbox.Service
}

func NewHandlers(ctx context.Context, projectID string, targetServiceID string, taskboxService *taskbox.Service) (*Handlers, error) {
	return &Handlers{
		projectID:       projectID,
		targetServiceID: targetServiceID,
		taskboxService:  taskboxService,
	}, nil
}
