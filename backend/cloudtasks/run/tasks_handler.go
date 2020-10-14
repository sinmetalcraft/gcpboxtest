package appengine

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/vvakame/sdlog/aelog"
	"golang.org/x/xerrors"
	"google.golang.org/api/idtoken"
)

func (h *Handlers) TasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "request.header", r.Header)
	if err := ValidateJWTFromCloudRun(r); err != nil {
		aelog.Errorf(ctx, "failed ValidateJWTFromAppEngine. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		log.InfoKV(ctx, "request.body", r.Body)
	}
	w.WriteHeader(http.StatusInternalServerError) // Tasks の Test のために必ずコケるようにする
}

type JWTPayload struct {
	Audience        string `json:"aud"`            // JWTを受け取る予定の対象者 example:https://gcpboxtest-73zry4yfvq-an.a.run.app/cloudtasks/run/json-post-task
	AuthorizedParty string `json:"azp"`            // 認可された対象者 example:102542703233071533897
	Email           string `json:"email"`          // JWT発行者のEmail example:sinmetal-ci@appspot.gserviceaccount.com
	EmailVerified   bool   `json:"email_verified"` // メールアドレスが検証済みか
	Expires         int    `json:"exp"`            // 有効期限(EpochTime seconds) example:1602514972
	IssuedAt        int    `json:"iat"`            // 発行日時(EpochTime seconds) example:1602511372
	Issuer          string `json:"iss"`            // 発行者 (issuer) example:https://accounts.google.com
	Subject         string `json:"sub"`            // UserID example:102542703233071533897
}

// validateJWTFromAppEngine validates a JWT found in the
func ValidateJWTFromCloudRun(r *http.Request) error {
	jwt := r.Header.Get("Authorization")
	tokens := strings.Split(jwt, " ")
	if len(tokens) < 1 {
		return fmt.Errorf("invalid token")
	}
	jwt = tokens[1]

	ctx := context.Background()

	// Cloud Tasksの場合は、audienceはTasksのUriと同じものになる
	//
	// curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" https://... みたいな感じでUserが投げた場合、
	// `{謎のID}.apps.googleusercontent.com` になるので、事前にaudienceを知るのは難しそうなので、JWTPayloadをParseして、
	// audienceを抜き出すことになりそう
	aud := "https://gcpboxtest-73zry4yfvq-an.a.run.app/cloudtasks/run/json-post-task"

	_, err := idtoken.Validate(ctx, jwt, aud)
	if err != nil {
		return fmt.Errorf("idtoken.Validate: %v", err)
	}
	payload, err := ParseJWTPayload(jwt)
	if err != nil {
		return err
	}

	log.InfoKV(ctx, "jwt.payload", payload)

	return nil
}

func ParseJWTPayload(jwt string) (*JWTPayload, error) {
	list := strings.Split(jwt, ".")
	if len(list) != 3 {
		return nil, fmt.Errorf("invalid JWT")
	}

	var payload *JWTPayload
	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(list[1]))
	if err := json.NewDecoder(r).Decode(payload); err != nil {
		return nil, xerrors.Errorf("invalid JWT :%w", err)
	}
	return payload, nil
}
