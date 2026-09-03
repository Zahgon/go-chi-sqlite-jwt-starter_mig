package server

import (
	"context"
	"net/http"

	"go-chi-sqlite-jwt-starter/internal/auth"
	category_group_handlers "go-chi-sqlite-jwt-starter/internal/category-group/handlers"
	models "go-chi-sqlite-jwt-starter/internal/models"
	provider "go-chi-sqlite-jwt-starter/internal/provider"
	"go-chi-sqlite-jwt-starter/internal/utils"

	"github.com/labstack/echo/v4"
)

func categoryGroupRouter(g *echo.Group) {
	auth.UseAuthMiddleware(g)

	g.GET("/list", echo.WrapHandler(http.HandlerFunc(category_group_handlers.ListCategoryGroups)))
	g.POST("/create", echo.WrapHandler(http.HandlerFunc(category_group_handlers.CreateCategoryGroup)))

	sub := g.Group("/:categoryGroupID")
	sub.Use(injectURLParam("categoryGroupID"))
	sub.Use(echo.WrapMiddleware(CategoryGroupCtx))
	sub.GET("", echo.WrapHandler(http.HandlerFunc(category_group_handlers.GetCategoryGroup)))
	sub.POST("/rename", echo.WrapHandler(http.HandlerFunc(category_group_handlers.RenameCategoryGroup)))
	sub.DELETE("", echo.WrapHandler(http.HandlerFunc(category_group_handlers.DeleteCategoryGroup)))
}

func CategoryGroupCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categoryGroupID := urlParam(r, "categoryGroupID")
		id, err := utils.StringToInt64(categoryGroupID)
		if err != nil {
			http.Error(w, "Invalid category group ID", http.StatusBadRequest)
			return
		}

		user := utils.GetUserFromContext(w, r.Context())
		catgoryGroup, err := provider.Provider.CategoryGroupService.GetCategoryGroupForUser(id, user.ID)
		if err != nil {
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), models.ContextKeys.CategoryGroup, catgoryGroup)
		ctx = context.WithValue(ctx, models.ContextKeys.CategoryGroupID, categoryGroupID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
