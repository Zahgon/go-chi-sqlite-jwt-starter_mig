package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// urlParamKey is the context key type used to carry echo path parameters into
// standard net/http middleware. Echo resolves path parameters on echo.Context,
// which the wrapped net/http middleware cannot access directly, so the value is
// stored on the request context before the wrapped middleware runs.
type urlParamKey string

// injectURLParam returns an echo middleware that copies the named echo path
// parameter onto the request context so that downstream net/http middleware can
// read it via urlParam. This mirrors chi.URLParam access for wrapped handlers.
func injectURLParam(name string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), urlParamKey(name), c.Param(name))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

// urlParam retrieves a path parameter previously stored by injectURLParam.
func urlParam(r *http.Request, name string) string {
	if v, ok := r.Context().Value(urlParamKey(name)).(string); ok {
		return v
	}
	return ""
}
