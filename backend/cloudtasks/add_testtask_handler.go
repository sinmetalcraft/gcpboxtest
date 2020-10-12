package cloudtasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
	"github.com/vvakame/sdlog/aelog"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

var (
	gcpboxQueue = &taskbox.Queue{
		ProjectID: "sinmetal-ci",
		Region:    "asia-northeast1",
		Name:      "gcpboxtest",
	}
)

func (h *Handlers) AddTask(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body := struct {
		Name string
	}{
		Name: "Hello World",
	}

	tnAppEngine, err := h.addAppEngineTask(ctx, body)
	if err != nil {
		aelog.Errorf(ctx, "failed addAppEngineTask. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tnHttpTask, err := h.addHttpTask(ctx, body)
	if err != nil {
		aelog.Errorf(ctx, "failed addHttpTask. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := struct {
		AppEngineTaskName string
		HttpTaskName      string
	}{
		AppEngineTaskName: tnAppEngine,
		HttpTaskName:      tnHttpTask,
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		aelog.Errorf(ctx, "failed Response.Write. err=%+v", err)
	}
}

func (h *Handlers) addAppEngineTask(ctx context.Context, body interface{}) (string, error) {
	taskName, err := h.taskboxService.CreateJsonPostTask(ctx, gcpboxQueue, &taskbox.JsonPostTask{
		Routing: &taskbox.Routing{
			Service: h.targetGAEServiceID,
		},
		RelativeUri: "/cloudtasks/appengine/json-post-task",
		Body:        &body,
	})
	return taskName, err
}

func (h *Handlers) addHttpTask(ctx context.Context, body interface{}) (string, error) {
	bb, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	aud := fmt.Sprintf("/projects/%s/run/%s", h.projectNumber, "gcpboxtest") // App Engineのものに似せて、ProjectNumberとCloud Run Service Nameを入れてみた
	task, err := h.cloudtasksClient.CreateTask(ctx, &taskspb.CreateTaskRequest{
		Parent: gcpboxQueue.Parent(),
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					Url:        "https://gcpboxtest-73zry4yfvq-an.a.run.app/cloudtasks/run/json-post-task",
					HttpMethod: taskspb.HttpMethod_POST,
					Body:       bb,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{OidcToken: &taskspb.OidcToken{
						ServiceAccountEmail: h.serviceAccountEmail,
						Audience:            aud, // Audienceを省略した場合、Cloud TasksがTaskのUrlを設定する
					}},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}
	return task.Name, nil
}
