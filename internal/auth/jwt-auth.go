package auth

import (
	"context"
	"go-chi-sqlite-jwt-starter/config"
	"go-chi-sqlite-jwt-starter/internal/models"
	"go-chi-sqlite-jwt-starter/internal/provider"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/labstack/echo/v4"
)

var tokenAuth *jwtauth.JWTAuth

func InitializeTokenVerifier() {
	secret := config.Variables.AUTH_PRIVATE_KEY
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func UseAuthMiddleware(g *echo.Group) {
	// Seek, verify and validate JWT tokens
	g.Use(echo.WrapMiddleware(jwtauth.Verifier(tokenAuth)))

	// Handle valid / invalid tokens. In this example, we use
	// the provided authenticator middleware, but you can write your
	// own very easily, look at the Authenticator method in jwtauth.go
	// and tweak it, its not scary.
	g.Use(echo.WrapMiddleware(jwtauth.Authenticator(tokenAuth)))

	// Check if the user exists in the database & add it to the context
	g.Use(echo.WrapMiddleware(myAuthMiddleware))
}

func myAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())

		userID := int64(claims["user_id"].(float64))
		user, err := provider.Provider.UserService.GetUser(userID)
		if err != nil {
			http.Error(w, "You are authenticated, but we could not find your account.", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), models.ContextKeys.User, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateJWT(user models.User) (string, error) {
	ttl := 7 * 24 * time.Hour // 1 week
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"role":    user.Role,
		"sub":     user.Username,
		"exp":     time.Now().UTC().Add(ttl).Unix(),
	})
	return tokenString, err
}
