package handler

import (
	"encoder/app"
	"encoder/helper"
	"encoder/m"
	"encoder/t"
	"encoder/views"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetSetup(c echo.Context) error {
	if app.Setting.HasBeenSetup {
		return c.NoContent(http.StatusNotFound)
	}

	return helper.Render(c,
		http.StatusOK,
		views.Setup(
			helper.TCtx(c),
			fmt.Sprintf("%s - Home", app.Name),
			t.SetupValidator{},
		),
	)
}

func PostSetup(c echo.Context) error {
	if app.Setting.HasBeenSetup {
		return c.NoContent(http.StatusNotFound)
	}

	var v t.SetupValidator
	if err := c.Bind(&v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setup(
				helper.TCtxWError(c, errors.New("bad request")),
				fmt.Sprintf("%s - Home", app.Name),
				v,
			),
		)
	}
	if err := app.Validate.Struct(v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setup(
				helper.TCtxWError(c, err),
				fmt.Sprintf("%s - Home", app.Name),
				v,
			),
		)
	}

	settingTmp := *app.Setting

	if v.EnableAuthentication != nil && *v.EnableAuthentication == "on" {
		settingTmp.EnableAuthentication = true
	} else {
		settingTmp.EnableAuthentication = false
	}

	if settingTmp.EnableAuthentication {
		if v.AuthenticationType == nil {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setup(
					helper.TCtxWError(c, errors.New("AuthenticationType cant be empty when authentication is enabled")),
					fmt.Sprintf("%s - Home", app.Name),
					v,
				),
			)
		}
		if v.Username == nil {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setup(
					helper.TCtxWError(c, errors.New("username cant be empty when authentication is enabled")),
					fmt.Sprintf("%s - Home", app.Name),
					v,
				),
			)
		}
		if v.Password == nil {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setup(
					helper.TCtxWError(c, errors.New("password cant be empty when authentication is enabled")),
					fmt.Sprintf("%s - Home", app.Name),
					v,
				),
			)
		}

		if len(helper.PStrToStr(v.Username)) < 4 {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setup(
					helper.TCtxWError(c, errors.New("username has to be more than, or 4 characters")),
					fmt.Sprintf("%s - Home", app.Name),
					v,
				),
			)
		}

		if len(helper.PStrToStr(v.Password)) < 8 {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setup(
					helper.TCtxWError(c, errors.New("password has to be more than, or 8 characters")),
					fmt.Sprintf("%s - Home", app.Name),
					v,
				),
			)
		}

		settingTmp.AuthenticationType = v.AuthenticationType

		pwdHash, err := bcrypt.GenerateFromPassword([]byte(*v.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Echo().Logger.Error("Failed to hash password", err)
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setup(
					helper.TCtxWError(c, errors.New("internal server error")),
					fmt.Sprintf("%s - Home", app.Name),
					v,
				),
			)
		}

		// create user
		settingTmp.Username = *v.Username
		settingTmp.PwdHash = string(pwdHash)
	}

	settingTmp.HasBeenSetup = true

	var setting m.Setting
	if err := app.DB.First(&setting).Error; err != nil {
		c.Echo().Logger.Error("Failed to get setting", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setup(
				helper.TCtxWError(c, errors.New("internal server error")),
				fmt.Sprintf("%s - Home", app.Name),
				v,
			),
		)
	}

	setting.Value = settingTmp

	if err := app.DB.Save(&setting).Error; err != nil {
		c.Echo().Logger.Error("Failed to update setting", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setup(
				helper.TCtxWError(c, errors.New("internal server error")),
				fmt.Sprintf("%s - Home", app.Name),
				v,
			),
		)
	}

	app.Setting = &settingTmp

	return c.Redirect(http.StatusFound, "/")
}
