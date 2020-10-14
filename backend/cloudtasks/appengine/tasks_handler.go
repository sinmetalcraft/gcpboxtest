package appengine

import (
	"context"
	"fmt"
	"net/http"

	taskbox "github.com/sinmetalcraft/gcpbox/cloudtasks/appengine"
	"github.com/sinmetalcraft/gcpboxtest/backend/log"
	"github.com/vvakame/sdlog/aelog"
	"google.golang.org/api/idtoken"
)

func (h *Handlers) TasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "request.header", r.Header)

	// Cloud Tasks からのRequestだと、 "X-Goog-IAP-JWT-Assertion" は付いてない
	// App Engine Task の場合は “X-Google-Internal-Skipadmincheck”: [“true”] というのが付いている
	if err := ValidateJWTFromAppEngine(r, h.projectNumber, h.projectID); err != nil {
		aelog.Errorf(ctx, "failed ValidateJWTFromAppEngine. pn:%s,pID:%s, %v\n", h.projectNumber, h.projectID, err)
	}

	if r.Method == http.MethodPost {
		log.InfoKV(ctx, "request.body", r.Body)
	}

	th, err := taskbox.GetHeader(r)
	if err != nil {
		aelog.Errorf(ctx, "failed taskbox.GetHeader. err=%+v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.InfoKV(ctx, "tasks.header", th)
}

func ValidateJWTHandlerFunc(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return validateJWTHandler(http.HandlerFunc(handler))
}

func validateJWTHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

// validateJWTFromAppEngine validates a JWT found in the
// "x-goog-iap-jwt-assertion" header.
func ValidateJWTFromAppEngine(r *http.Request, projectNumber string, projectID string) error {
	iapJWT := r.Header.Get("X-Goog-IAP-JWT-Assertion") // req.Header.Get("X-Goog-IAP-JWT-Assertion")
	// projectNumber := "123456789"
	// projectID := "your-project-id"
	ctx := context.Background()
	aud := fmt.Sprintf("/projects/%s/apps/%s", projectNumber, projectID)

	payload, err := idtoken.Validate(ctx, iapJWT, aud)
	if err != nil {
		return fmt.Errorf("idtoken.Validate: %v", err)
	}

	log.InfoKV(ctx, "idtolen.payload", payload)

	return nil
}
