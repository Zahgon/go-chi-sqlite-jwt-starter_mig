package server

import (
	"context"
	"net/http"

	"go-chi-sqlite-jwt-starter/internal/auth"
	category_handlers "go-chi-sqlite-jwt-starter/internal/category/handlers"
	models "go-chi-sqlite-jwt-starter/internal/models"
	provider "go-chi-sqlite-jwt-starter/internal/provider"
	"go-chi-sqlite-jwt-starter/internal/utils"
	"go-chi-sqlite-jwt-starter/internal/validation"

	"github.com/labstack/echo/v4"
)

func categoryRouter(g *echo.Group) {
	auth.UseAuthMiddleware(g)

	g.GET("/list", echo.WrapHandler(http.HandlerFunc(category_handlers.ListCategories)))
	g.POST("/create", echo.WrapHandler(http.HandlerFunc(category_handlers.CreateCategory)))

	sub := g.Group("/:categoryID")
	sub.Use(injectURLParam("categoryID"))
	sub.Use(echo.WrapMiddleware(CategoryCtx))
	sub.GET("", echo.WrapHandler(http.HandlerFunc(category_handlers.GetCategory)))
	sub.PUT("", echo.WrapHandler(http.HandlerFunc(category_handlers.UpdateCategory)))
	sub.DELETE("", echo.WrapHandler(http.HandlerFunc(category_handlers.DeleteCategory)))
}

func CategoryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categoryID := urlParam(r, "categoryID")
		id, err := utils.StringToInt64(categoryID)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		catgory, err := provider.Provider.CategoryService.GetCategory(id)
		if err != nil {
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			return
		}

		user := utils.GetUserFromContext(w, r.Context())
		err = validation.HasAccessToCategoryGroup(w, catgory.CategoryGroupID, user.ID)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), models.ContextKeys.Category, catgory)
		ctx = context.WithValue(ctx, models.ContextKeys.CategoryID, categoryID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
