package post

import (
	resp "go-anonboard/internal/lib/api/response"
	"go-anonboard/internal/lib/sl"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/matoous/go-nanoid/v2"
)

type Request struct {
	Message string `json:"message" validate:"required"`
}

type Response struct {
	resp.Response
	Post *Post `json:"post,omitempty"`
}

type Post struct {
	ID        int64     `json:"id"`
	NanoID    string    `json:"nanoid"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type PostSaver interface {
	SavePost(nanoid, message string) (int64, time.Time, error)
}

func Save(log *slog.Logger, postSaver PostSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.post.Save"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		nanoid := generateNanoid()
		id, created_at, err := postSaver.SavePost(nanoid, req.Message)
		if err != nil {
			log.Error("failed to add post", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add post"))
		}

		log.Info("post added", slog.Int64("id", id))

		var post Post
		post.ID = id
		post.NanoID = nanoid
		post.Message = req.Message
		post.CreatedAt = created_at

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Post:     &post,
		})
	}
}

func generateNanoid() string {
	id, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz1234567890", 10)
	return id
}
