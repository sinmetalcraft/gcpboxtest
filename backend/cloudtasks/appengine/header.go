package appengine

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"google.golang.org/api/idtoken"
)

func IsGCPInternal(r *http.Request) bool {
	return r.Header.Get("X-Google-Internal-Skipadmincheck") == "true"
}

// validateJWTFromAppEngine validates a JWT found in the
// "x-goog-iap-jwt-assertion" header.
func ValidateJWTFromAppEngine(ctx context.Context, r *http.Request, projectNumber string, projectID string) error {
	iapJWT := r.Header.Get("X-Goog-IAP-JWT-Assertion")
	aud := fmt.Sprintf("/projects/%s/apps/%s", projectNumber, projectID)

	payload, err := idtoken.Validate(ctx, iapJWT, aud)
	if err != nil {
		return fmt.Errorf("idtoken.Validate: %v", err)
	}

	log.InfoKV(ctx, "idtolen.payload", payload)

	return nil
}
