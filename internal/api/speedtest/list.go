package speedtest

import (
	"context"
	"net/http"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/labstack/echo/v5"
)

type listRepo interface {
	List(context.Context, results.ListOpts) (*database.Pagination[results.Result], error)
}

type listHandler struct {
	repo listRepo
}

type ListRequest struct {
	Page  int `query:"page"  default:"1"  description:"The page number of the results"`
	Limit int `query:"limit" default:"25" description:"The number of rows per page"`
}

func (l *listHandler) Handler() handlers.Handler[ListRequest, database.Pagination[Result]] {
	return func(c *echo.Context, req ListRequest) (*database.Pagination[Result], error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "List results")
		defer span.End()

		results, err := l.repo.List(ctx, results.ListOpts{
			Page:  req.Page,
			Limit: req.Limit,
		})
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		out := &database.Pagination[Result]{
			Limit:      results.Limit,
			Page:       results.Page,
			Total:      results.Total,
			TotalPages: results.TotalPages,
		}
		data := []Result{}
		for _, d := range results.Data {
			data = append(data, convert(c, d))
		}
		out.Data = data

		return out, nil
	}
}

func (l *listHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (l *listHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "List results",
		Tag:          "Speedtest",
		Method:       http.MethodGet,
		Path:         "/v1/speedtest",
		Code:         http.StatusOK,
		Description:  "List speedtest results",
		GenerateSpec: true,
	}
}
