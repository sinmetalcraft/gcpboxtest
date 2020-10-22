package run

import (
	"fmt"
	"net/http"

	tasksbox "github.com/sinmetalcraft/gcpbox/cloudtasks"
	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/vvakame/sdlog/aelog"
)

func (h *Handlers) TasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "request.header", r.Header)

	// Cloud Tasksの場合は、audienceはTasksのUriと同じものになる
	//
	// curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" https://... みたいな感じでUserが投げた場合、
	// `{謎のID}.apps.googleusercontent.com` になるので、事前にaudienceを知るのは難しそうなので、JWTPayloadをParseして、
	// audienceを抜き出すことになりそう
	payload, err := tasksbox.ValidateJWTFromInvoker(ctx, r, fmt.Sprintf("%s/cloudtasks/run/json-post-task", h.gcpboxtestCloudRunService.URL))
	if err != nil {

		aelog.Errorf(ctx, "failed ValidateJWTFromCloudRun. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.InfoKV(ctx, "JWT.payload", payload)

	if r.Method == http.MethodPost {
		log.InfoKV(ctx, "request.body", r.Body)
	}

	th, err := tasksbox.GetHeader(r)
	if err != nil {
		aelog.Errorf(ctx, "failed tasksbox.GetHeader. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.InfoKV(ctx, "tasks.Header", th)
}
