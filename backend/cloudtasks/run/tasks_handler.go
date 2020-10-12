package appengine

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/vvakame/sdlog/aelog"
	"google.golang.org/api/idtoken"
)

func (h *Handlers) TasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "request.header", r.Header)
	if err := ValidateJWTFromAppEngine(r, h.projectNumber, h.projectID); err != nil {
		aelog.Errorf(ctx, "failed ValidateJWTFromAppEngine. err=%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		log.InfoKV(ctx, "request.body", r.Body)
	}
}

// validateJWTFromAppEngine validates a JWT found in the
func ValidateJWTFromAppEngine(r *http.Request, projectNumber string, projectID string) error {
	jwt := r.Header.Get("Authorization")
	tokens := strings.Split(jwt, " ")
	if len(tokens) < 1 {
		return fmt.Errorf("invalid token")
	}
	jwt = tokens[1]
	// projectNumber := "123456789"
	// projectID := "your-project-id"
	ctx := context.Background()
	aud := fmt.Sprintf("/projects/%s/apps/%s", projectNumber, projectID)

	payload, err := idtoken.Validate(ctx, jwt, aud)
	if err != nil {
		return fmt.Errorf("idtoken.Validate: %v", err)
	}

	log.InfoKV(ctx, "idtolen.payload", payload)

	return nil
}
