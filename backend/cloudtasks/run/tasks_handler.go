package run

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sinmetalcraft/gcpboxtest/backend/jwt"
	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/vvakame/sdlog/aelog"
	"google.golang.org/api/idtoken"
)

func (h *Handlers) TasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "request.header", r.Header)

	// Cloud Tasksの場合は、audienceはTasksのUriと同じものになる
	//
	// curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" https://... みたいな感じでUserが投げた場合、
	// `{謎のID}.apps.googleusercontent.com` になるので、事前にaudienceを知るのは難しそうなので、JWTPayloadをParseして、
	// audienceを抜き出すことになりそう
	payload, err := ValidateJWTFromCloudRun(ctx, r, "https://gcpboxtest-73zry4yfvq-an.a.run.app/cloudtasks/run/json-post-task")
	if err != nil {
		aelog.Errorf(ctx, "failed ValidateJWTFromCloudRun. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.InfoKV(ctx, "JWT.payload", payload)

	if r.Method == http.MethodPost {
		log.InfoKV(ctx, "request.body", r.Body)
	}
}

// validateJWTFromAppEngine validates a JWT found in the
func ValidateJWTFromCloudRun(ctx context.Context, r *http.Request, audience string) (*jwt.JWTPayload, error) {
	autzHeader := r.Header.Get("Authorization")
	tokens := strings.Split(autzHeader, " ")
	if len(tokens) < 1 {
		return nil, fmt.Errorf("invalid token")
	}
	autzHeader = tokens[1]

	_, err := idtoken.Validate(ctx, autzHeader, audience)
	if err != nil {
		return nil, fmt.Errorf("idtoken.Validate: %v", err)
	}
	payload, err := jwt.ParseJWTPayload(autzHeader)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
