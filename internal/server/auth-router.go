package server

import (
	"net/http"

	auth_handlers "go-chi-sqlite-jwt-starter/internal/auth/handlers"

	"github.com/labstack/echo/v4"
)

func authRouter(g *echo.Group) {
	g.POST("/register", echo.WrapHandler(http.HandlerFunc(auth_handlers.RegisterAccount)))
	g.POST("/login", echo.WrapHandler(http.HandlerFunc(auth_handlers.Login)))
}
