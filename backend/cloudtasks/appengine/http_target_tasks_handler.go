package appengine

import (
	"net/http"

	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/vvakame/sdlog/aelog"
)

const HttpTargetTasksHandlerUri = "/cloudtasks/appengine/http-target-task"

func (h *Handlers) HttpTargetTasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "HttpTargetTasksHandler.request.header", r.Header)

	// Cloud Tasks Http Target Task を IAP 貫通させると X-Goog-Iap-Jwt-Assertion が付いてる
	if err := ValidateJWTFromAppEngine(ctx, r, h.projectNumber, h.projectID); err != nil {
		aelog.Errorf(ctx, "failed ValidateJWTFromAppEngine. pn:%s,pID:%s, %v\n", h.projectNumber, h.projectID, err)
	}

	if r.Method == http.MethodPost {
		log.InfoKV(ctx, "request.body", r.Body)
	}

	// TODO Header Check
}
