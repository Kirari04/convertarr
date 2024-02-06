package mware

import (
	"encoder/app"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func AuthEnabledRedirect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if app.Setting.EnableAuthentication {
			if c.Request().URL.Path == "/favicon.ico" {
				return next(c)
			}
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

				if c.Request().URL.Path != "/login" {
					return c.Redirect(http.StatusFound, "/login")
				}
			}
			if *app.Setting.AuthenticationType == "basic" {
				basicRawCredentials := c.Request().Header.Get("Authorization")
				if basicRawCredentials == "" {
					c.Response().Header().Add("WWW-Authenticate", `Basic realm="Restricted Access"`)
					return c.String(http.StatusUnauthorized, "Unauthorized")
				}

				basicRawCredentials = strings.TrimPrefix(basicRawCredentials, "Basic ")

				basicCredentials, err := base64.StdEncoding.DecodeString(basicRawCredentials)
				if err != nil {
					return c.String(http.StatusUnauthorized, "Unauthorized")
				}
				credentialsSlice := strings.Split(string(basicCredentials), ":")
				if len(credentialsSlice) != 2 {
					return c.String(http.StatusUnauthorized, "Unauthorized")
				}

				if app.Setting.Username != credentialsSlice[0] {
					return c.String(http.StatusUnauthorized, "Unauthorized")
				}

				if err := bcrypt.CompareHashAndPassword([]byte(app.Setting.PwdHash), []byte(credentialsSlice[1])); err != nil {
					return c.String(http.StatusUnauthorized, "Unauthorized")
				}

				c.Set("IsAuth", true)

				return next(c)
			}
			if c.Request().URL.Path != "/login" {
				c.Logger().Error("unknown AuthenticationType: ", *app.Setting.AuthenticationType)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
		return next(c)
	}
}
