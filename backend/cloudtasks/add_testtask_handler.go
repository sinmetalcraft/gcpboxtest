package cloudtasks

import (
	"context"
	"encoding/json"
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
	tnHttpTask, err := h.addHttpTask(ctx, "https://gcpboxtest-73zry4yfvq-an.a.run.app/cloudtasks/run/json-post-task", body)
	if err != nil {
		aelog.Errorf(ctx, "failed addHttpTask. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tnHttpTaskToGAE, err := h.addHttpTask(ctx, "https://gcpbox-dot-sinmetal-ci.an.r.appspot.com/cloudtasks/appengine/json-post-task", body)
	if err != nil {
		aelog.Errorf(ctx, "failed addHttpTask. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := struct {
		AppEngineTaskName       string
		HttpTaskToCloudRunName  string
		HttpTaskToAppEngineName string
	}{
		AppEngineTaskName:       tnAppEngine,
		HttpTaskToCloudRunName:  tnHttpTask,
		HttpTaskToAppEngineName: tnHttpTaskToGAE,
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

func (h *Handlers) addHttpTask(ctx context.Context, url string, body interface{}) (string, error) {
	bb, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	task, err := h.cloudtasksClient.CreateTask(ctx, &taskspb.CreateTaskRequest{
		Parent: gcpboxQueue.Parent(),
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					Url:        url,
					HttpMethod: taskspb.HttpMethod_POST,
					Body:       bb,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{OidcToken: &taskspb.OidcToken{
						ServiceAccountEmail: h.serviceAccountEmail,
						Audience:            "401580979819-84sh4g7gpk72m6lfum4oildt8pjpvmse.apps.googleusercontent.com", // IAPに向けて投げる時は、IAPのClient IDを指定する https://cloud.google.com/iap/docs/authentication-howto#authenticating_from_a_service_account
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
