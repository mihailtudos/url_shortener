package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/mihailtudos/url-shortener/internal/lib/api/response"
	"github.com/mihailtudos/url-shortener/internal/lib/logger/sl"
	"github.com/mihailtudos/url-shortener/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		log.Info("alias", slog.String("alias", alias))

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				log.Info("url not found", slog.String("alias", alias))
				
				render.JSON(w, r, response.Error("not found"))
				return
			}

			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", url))

		// redirect to the URL
		http.Redirect(w, r, url, http.StatusFound)
	}
}
