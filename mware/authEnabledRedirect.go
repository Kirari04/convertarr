package mware

import (
	"encoder/app"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func AuthEnabledRedirect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if app.Setting.EnableAuthentication {
			if app.Setting.AuthenticationType == nil {
				c.Logger().Error("AuthenticationType is nil while EnableAuthentication is true")
				return c.NoContent(http.StatusInternalServerError)
			}
			if *app.Setting.AuthenticationType == "form" {
				// if session is correct skip all
				if jwtToken, err := c.Cookie("session"); err == nil {
					token, err := jwt.Parse(jwtToken.Value, func(token *jwt.Token) (interface{}, error) {
						// Don't forget to validate the alg is what you expect:
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
						}
						return []byte(app.JwtSecret), nil
					})
					if err != nil {
						c.SetCookie(&http.Cookie{
							Name:     "session",
							Value:    "",
							Path:     "/",
							HttpOnly: true,
							Expires:  time.Now().Add(-1 * 24 * time.Hour),
						})
						return c.Redirect(http.StatusFound, "/login")
					}

					if !token.Valid {
						c.SetCookie(&http.Cookie{
							Name:     "session",
							Value:    "",
							Path:     "/",
							HttpOnly: true,
							Expires:  time.Now().Add(-1 * 24 * time.Hour),
						})
						return c.Redirect(http.StatusFound, "/login")
					}

					c.Set("IsAuth", true)

					// session is valid
					return next(c)
				}

				if c.Request().URL.Path != "/login" && c.Request().URL.Path != "/favicon.ico" {
					return c.Redirect(http.StatusFound, "/login")
				}
			}
			if *app.Setting.AuthenticationType == "basic" {
				return c.NoContent(http.StatusNotImplemented)
			}
			if c.Request().URL.Path != "/login" && c.Request().URL.Path != "/favicon.ico" {
				c.Logger().Error("unknown AuthenticationType: ", *app.Setting.AuthenticationType)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
		return next(c)
	}
}
