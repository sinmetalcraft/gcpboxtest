package appengine

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sinmetalcraft/gcpboxtest/backend/jwt"
	"google.golang.org/api/idtoken"
)

func IsGCPInternal(r *http.Request) bool {
	return r.Header.Get("X-Google-Internal-Skipadmincheck") == "true"
}

// validateJWTFromAppEngine validates a JWT found in the
// "x-goog-iap-jwt-assertion" header.
func ValidateJWTFromAppEngine(ctx context.Context, r *http.Request, projectNumber string, projectID string) (*jwt.JWTPayload, error) {
	iapJWT := r.Header.Get("X-Goog-IAP-JWT-Assertion")
	aud := fmt.Sprintf("/projects/%s/apps/%s", projectNumber, projectID)

	_, err := idtoken.Validate(ctx, iapJWT, aud)
	if err != nil {
		return nil, fmt.Errorf("idtoken.Validate: %v", err)
	}

	payload, err := jwt.ParseJWTPayload(iapJWT)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
