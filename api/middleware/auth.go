package middleware

import (
	"strings"

	"github.com/coderero/erochat-server/api/utils"
	"github.com/coderero/erochat-server/interfaces"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware is a middleware that checks if the user is authenticated.
func JWTMiddleware(jwt interfaces.TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the request.
			rToken := c.Request().Header.Get("Authorization")

			if rToken != "" {
				bearer := strings.Split(rToken, " ")
				if len(bearer) != 2 || bearer[0] != "Bearer" || len(bearer[1]) == 0 {
					return echo.ErrUnauthorized
				}
				// Validate the token.
				valid, err := jwt.ValidateToken(bearer[1])
				if err != nil || !valid {
					return echo.ErrUnauthorized
				}

				if !valid {
					return echo.ErrUnauthorized
				} else {
					// Set the user email in the context.
					err := GetAndSetToContext(c, jwt, bearer[1])
					if err != nil {
						return err
					}
					return next(c)
				}
			} else {
				// Cookie token
				var (
					accessToken  string
					refreshToken string
				)

				// Get the access token from the cookie.
				accessToken = utils.GetCookie(c, "__a")
				refreshToken = utils.GetCookie(c, "__r")

				if accessToken == "" && refreshToken == "" {
					return echo.ErrUnauthorized
				}

				accessValid, aErr := jwt.ValidateToken(accessToken)

				if aErr != nil || !accessValid {
					refreshValid, rErr := jwt.ValidateToken(refreshToken)
					if rErr != nil || !refreshValid {
						return echo.ErrUnauthorized
					}

					if refreshValid {
						token, rErr := jwt.RefreshToken(refreshToken)
						if rErr != nil {
							return echo.ErrUnauthorized
						}
						// Set the token in the cookie.
						utils.SaveCookie(c, "__a", token)

						// Set the user email in the context.
						err := GetAndSetToContext(c, jwt, token)
						if err != nil {
							return err
						}
						return next(c)
					}
					return echo.ErrUnauthorized
				}

				if accessValid {
					// Set the user email in the context.
					err := GetAndSetToContext(c, jwt, accessToken)
					if err != nil {
						return err
					}
					return next(c)
				}
			}
			return echo.ErrUnauthorized
		}
	}
}

func GetAndSetToContext(c echo.Context, jwt interfaces.TokenService, token string) error {
	// Get the claims from the token.
	claims, err := jwt.GetClaims(token)
	if err != nil {
		return echo.ErrUnauthorized
	}
	// Set the user email in the context.
	c.Set("user", claims["sub"])
	c.Set("uid", claims["uid"])

	return nil
}
