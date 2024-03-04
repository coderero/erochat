package middleware

type AuthRouteMiddlewareConfig struct {
	// Skip is a list of routes to skip.
	Skip []string
}

// AuthRouteMiddleware is a middleware that checks if the user is already authenticated.
// TODO: Implement the AuthRouteMiddleware function.
