package cloudtasks

import (
	"net/http"

	taskbox "github.com/sinmetal/gcpbox/cloudtasks/appengine"
	"github.com/sinmetal/gcpboxtest/backend"
	"github.com/vvakame/sdlog/aelog"
)

func AppEnginePushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	backend.Info(ctx, r.Header)

	if r.Method == http.MethodPost {
		backend.Info(ctx, r.Body)
	}

	th, err := taskbox.GetHeader(r)
	if err != nil {
		aelog.Errorf(ctx, "failed taskbox.GetHeader. err=%+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	backend.Info(ctx, th)
}
