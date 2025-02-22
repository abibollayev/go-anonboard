package post

import (
	resp "go-anonboard/internal/lib/api/response"
	"go-anonboard/internal/lib/sl"
	"go-anonboard/internal/storage/postgres"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type PostGeter interface {
	GetAllPost() ([]postgres.Post, error)
}

type GetResponse struct {
	resp.Response
	Post []postgres.Post `json:"post,omitempty"`
}

func GetAll(log *slog.Logger, postGeter PostGeter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.post.GetAll"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		posts, err := postGeter.GetAllPost()
		if err != nil {
			log.Error("failed to get all post", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get all post"))

			return
		}

		log.Info("get all post")

		render.JSON(w, r, GetResponse{
			Response: resp.OK(),
			Post:     posts,
		})
	}
}
