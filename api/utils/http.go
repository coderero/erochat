package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SaveCookie(c echo.Context, key, value string) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}

func DeleteCookie(c echo.Context, key string) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = ""
	cookie.Expires = time.Now()
	c.SetCookie(cookie)
}

func GetCookie(c echo.Context, key string) string {
	cookie, err := c.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.Value
}
