package appengine

import (
	"net/http"

	authbox "github.com/sinmetalcraft/gcpbox/auth"
	tasksbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/vvakame/sdlog/aelog"
)

const AppEngineTasksHandlerUri = "/cloudtasks/appengine/json-post-task"

func (h *Handlers) AppEngineTasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "AppEngineTasksHandler.request.header", r.Header)

	// App Engine Task の場合は “X-Google-Internal-Skipadmincheck”: [“true”] が付いてる
	if !authbox.IsGCPInternal(r) {
		aelog.Errorf(ctx, "IsGCPInternal is false...")
		http.Error(w, "IsGCPInternal is false...", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		log.InfoKV(ctx, "request.body", r.Body)
	}

	th, err := tasksbox.GetHeader(r)
	if err != nil {
		aelog.Errorf(ctx, "failed taskbox.GetHeader. err=%+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.InfoKV(ctx, "tasks.header", th)
}
