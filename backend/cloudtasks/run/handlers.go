package run

import (
	"context"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
)

type Handlers struct {
	projectID       string
	projectNumber   string
	targetServiceID string
	taskboxService  *cloudtasks.Client
}

func NewHandlers(ctx context.Context, projectID string, projectNumber string, tasksService *cloudtasks.Client) (*Handlers, error) {
	return &Handlers{
		projectID:      projectID,
		projectNumber:  projectNumber,
		taskboxService: tasksService,
	}, nil
}
