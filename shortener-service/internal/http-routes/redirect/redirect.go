package redirect

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
)

type Redirector interface {
	Redirect(ctx context.Context, urlId string) (string, error)
}

func New(log *zerolog.Logger, redirector Redirector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		id := mux.Vars(r)["id"]
		if id == "" {
			return
		}
		log.Info().
			Str("urlId", id).
			Msg("handle redirect request")

		resUrl, err := redirector.Redirect(r.Context(), id)
		if err != nil {
			log.Warn().Err(err)
		}
		log.Info().
			Str("url", resUrl).
			Msg("successfully redirected to url")
		http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
