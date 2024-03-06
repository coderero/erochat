package middleware

import (
	"github.com/coderero/erochat-server/interfaces"
	"github.com/labstack/echo/v4"
)

type AuthRouteMiddlewareConfig struct {
	// Skip is a list of routes to skip.
	Skip []string

	// TokenService is the token service.
	TokenService interfaces.TokenService
}

// AuthRouteMiddleware is a middleware that checks if the user is already authenticated.
// TODO: Implement the AuthRouteMiddleware function.
func AuthRouteMiddleware(config AuthRouteMiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
