package redirect

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type Redirector interface {
	Redirect(ctx context.Context, urlId string) (string, error)
}

func New(log *zerolog.Logger, redirector Redirector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		start := time.Now()
		var statusCode int
		defer func() {
			observeRedirect(time.Since(start), statusCode)
		}()

		id := mux.Vars(r)["id"]
		if id == "" {
			statusCode = http.StatusBadRequest
			http.Error(w, "id is required", statusCode)
			return
		}
		log.Info().
			Str("urlId", id).
			Msg("handle redirect request")

		resUrl, err := redirector.Redirect(r.Context(), id)
		if err != nil {
			log.Warn().Err(err)
			statusCode = http.StatusInternalServerError
			http.Error(w, http.StatusText(statusCode), statusCode)
			return
		}
		log.Info().
			Str("url", resUrl).
			Msg("successfully redirected to url")
		statusCode = http.StatusFound
		http.Redirect(w, r, resUrl, statusCode)
	}
}
