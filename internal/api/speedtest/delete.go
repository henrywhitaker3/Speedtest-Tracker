package speedtest

import (
	"context"
	"net/http"

	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/labstack/echo/v5"
)

type deleteRepo interface {
	Delete(context.Context, uuid.UUID) error
}

type deleteHandler struct {
	repo deleteRepo
}

type DeleteRequest struct {
	ID uuid.UUID `path:"id" param:"id"`
}

func (d *deleteHandler) Handler() handlers.Handler[DeleteRequest, any] {
	return func(c *echo.Context, req DeleteRequest) (*any, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Delete result")
		defer span.End()

		if err := d.repo.Delete(ctx, req.ID); err != nil {
			return nil, httperrors.Stack(err)
		}

		return nil, nil
	}
}

func (d *deleteHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (d *deleteHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Delete result",
		Tag:          "Speedtest",
		ID:           new("deleteResult"),
		Path:         "/v1/speedtest/:id",
		Method:       http.MethodDelete,
		Code:         http.StatusOK,
		GenerateSpec: true,
	}
}
