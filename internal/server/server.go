package server

import (
	"go-chi-sqlite-jwt-starter/internal/auth"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Initialize() *echo.Echo {
	log.Println("Initializing server...")
	defer log.Println("Server initialized")

	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	useGlobalMiddleware(e)
	auth.InitializeTokenVerifier()

	categoryRouter(e.Group("/category"))
	categoryGroupRouter(e.Group("/category-group"))
	adminRouter(e.Group("/admin"))
	authRouter(e.Group("/auth"))

	return e
}

func useGlobalMiddleware(e *echo.Echo) {
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, ".")
	})
}
