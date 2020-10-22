package appengine

import (
	"context"

	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
)

type Handlers struct {
	projectID       string
	projectNumber   string
	targetServiceID string
	taskboxService  *taskbox.Service
}

func NewHandlers(ctx context.Context, projectID string, projectNumber string, targetServiceID string, taskboxService *taskbox.Service) (*Handlers, error) {
	return &Handlers{
		projectID:       projectID,
		projectNumber:   projectNumber,
		targetServiceID: targetServiceID,
		taskboxService:  taskboxService,
	}, nil
}
