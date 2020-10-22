package appengine

import (
	"net/http"

	tasksbox "github.com/sinmetalcraft/gcpbox/cloudtasks"
	gaetasksbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/vvakame/sdlog/aelog"
)

const HttpTargetTasksHandlerUri = "/cloudtasks/appengine/http-target-task"

func (h *Handlers) HttpTargetTasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "HttpTargetTasksHandler.request.header", r.Header)

	// Cloud Tasks Http Target Task を IAP 貫通させると X-Goog-Iap-Jwt-Assertion が付いてる
	payload, err := gaetasksbox.ValidateJWTFromHttpTargetTask(ctx, r, h.projectNumber, h.projectID)
	if err != nil {
		aelog.Errorf(ctx, "failed ValidateJWTFromAppEngine. pn:%s,pID:%s, %v\n", h.projectNumber, h.projectID, err)
	}

	log.InfoKV(ctx, "JWT.payload", payload)

	if r.Method == http.MethodPost {
		log.InfoKV(ctx, "request.body", r.Body)
	}

	th, err := tasksbox.GetHeader(r)
	if err != nil {
		aelog.Errorf(ctx, "failed taskbox.GetHeader. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.InfoKV(ctx, "cloudtasks.Header", th)
}
