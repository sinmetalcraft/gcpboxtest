package cloudtasks

import (
	"context"
	"fmt"
	"net/http"

	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
	"github.com/vvakame/sdlog/aelog"
)

func (h *Handlers) AddTask(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body := struct {
		Name string
	}{
		Name: "Hello World",
	}
	taskName, err := h.taskboxService.CreateJsonPostTask(ctx, &taskbox.Queue{
		ProjectID: "sinmetal-ci",
		Region:    "asia-northeast1",
		Name:      "gcpboxtest",
	}, &taskbox.JsonPostTask{
		Routing: &taskbox.Routing{
			Service: h.targetGAEServiceID,
		},
		RelativeUri: "/cloudtasks/appengine/json-post-task",
		Body:        &body,
	})
	if err != nil {
		aelog.Errorf(ctx, "failed taskbox.CreateJsonPostTask. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(fmt.Sprintf("TaskName:%s", taskName)))
	if err != nil {
		aelog.Errorf(ctx, "failed Response.Write. err=%+v", err)
	}
}
