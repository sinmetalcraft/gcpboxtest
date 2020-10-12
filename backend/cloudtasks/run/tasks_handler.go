package appengine

import (
	"net/http"

	"github.com/sinmetalcraft/gcpboxtest/backend/log"
)

func (h *Handlers) TasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.InfoKV(ctx, "request.header", r.Header)

	if r.Method == http.MethodPost {
		log.InfoKV(ctx, "request.body", r.Body)
	}
}
