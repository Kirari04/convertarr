package handler

import (
	"encoder/app"
	"encoder/helper"
	"encoder/t"
	"encoder/views"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetLogin(c echo.Context) error {
	return helper.Render(c,
		http.StatusOK,
		views.Login(
			helper.TCtx(c),
			fmt.Sprintf("%s - Login", app.Name),
			t.LoginValidator{},
		),
	)
}

func PostLogin(c echo.Context) error {
	if !app.Setting.EnableAuthentication {
		return c.NoContent(http.StatusNotFound)
	}

	if app.Setting.AuthenticationType == nil || *app.Setting.AuthenticationType != "form" {
		return c.NoContent(http.StatusNotFound)
	}

	var v t.LoginValidator
	if err := c.Bind(&v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.Login(
				helper.TCtxWError(c, errors.New("bad request")),
				fmt.Sprintf("%s - Login", app.Name),
				v,
			),
		)
	}
	if err := app.Validate.Struct(v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.Login(
				helper.TCtxWError(c, err),
				fmt.Sprintf("%s - Login", app.Name),
				v,
			),
		)
	}

	if app.Setting.Username != v.Username {
		return helper.Render(c,
			http.StatusBadRequest,
			views.Login(
				helper.TCtxWError(c, errors.New("username doesn't match")),
				fmt.Sprintf("%s - Login", app.Name),
				v,
			),
		)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(app.Setting.PwdHash), []byte(v.Password)); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.Login(
				helper.TCtxWError(c, errors.New("password doesn't match")),
				fmt.Sprintf("%s - Login", app.Name),
				v,
			),
		)
	}

	// Set custom claims
	expiresAt := time.Now().Add(time.Hour * 72)
	claims := &t.JwtUserClaims{
		Username: v.Username,
	}
	claims.ExpiresAt = jwt.NewNumericDate(expiresAt)

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(app.JwtSecret))
	if err != nil {
		c.Logger().Error("Failed to generate jwt token", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.Login(
				helper.TCtxWError(c, errors.New("internal server error")),
				fmt.Sprintf("%s - Login", app.Name),
				v,
			),
		)
	}

	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    t,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})

	return c.Redirect(http.StatusFound, "/")
}
