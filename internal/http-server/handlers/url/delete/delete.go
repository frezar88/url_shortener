package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url_shortener/internal/lib/api/response"
	"url_shortener/internal/lib/logger/sl"
)

type URLDelete interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.Status(r, 409)

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlDelete.DeleteURL(alias)

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("success delete", slog.String("alias", alias))
		render.JSON(w, r, "success delete")
	}
}
