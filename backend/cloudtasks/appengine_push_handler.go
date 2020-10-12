package cloudtasks

import (
	"net/http"

	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/vvakame/sdlog/aelog"
)

func (h *Handlers) AppEnginePushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Info(ctx, r.Header)

	if r.Method == http.MethodPost {
		log.Info(ctx, r.Body)
	}

	th, err := taskbox.GetHeader(r)
	if err != nil {
		aelog.Errorf(ctx, "failed taskbox.GetHeader. err=%+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Info(ctx, th)
}
