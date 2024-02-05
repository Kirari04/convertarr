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

func GetSettingUser(c echo.Context) error {
	var v t.SettingUserValidation
	v.Username = app.Setting.Username
	return helper.Render(c,
		http.StatusOK,
		views.SettingUser(
			helper.TCtx(c),
			fmt.Sprintf("%s - Setting User", app.Name),
			v,
		),
	)
}

func PostSettingUser(c echo.Context) error {
	var v t.SettingUserValidation
	if err := c.Bind(&v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingUser(
				helper.TCtxWError(c, errors.New("bad request")),
				fmt.Sprintf("%s - Setting User", app.Name),
				v,
			),
		)
	}
	if err := app.Validate.Struct(v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingUser(
				helper.TCtxWError(c, err),
				fmt.Sprintf("%s - Setting User", app.Name),
				v,
			),
		)
	}

	settingTmp := *app.Setting

	settingTmp.Username = v.Username
	if len(v.Password) > 0 {
		pwdHash, err := bcrypt.GenerateFromPassword([]byte(v.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Echo().Logger.Error("Failed to hash password", err)
			return helper.Render(c,
				http.StatusBadRequest,
				views.SettingUser(
					helper.TCtxWError(c, errors.New("internal server error")),
					fmt.Sprintf("%s - Setting User", app.Name),
					v,
				),
			)
		}
		settingTmp.PwdHash = string(pwdHash)
	}

	var setting m.Setting
	if err := app.DB.First(&setting).Error; err != nil {
		c.Echo().Logger.Error("Failed to get setting", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingUser(
				helper.TCtxWError(c, errors.New("internal server error")),
				fmt.Sprintf("%s - Setting User", app.Name),
				v,
			),
		)
	}

	setting.Value = settingTmp

	if err := app.DB.Save(&setting).Error; err != nil {
		c.Echo().Logger.Error("Failed to update setting", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingUser(
				helper.TCtxWError(c, errors.New("internal server error")),
				fmt.Sprintf("%s - Setting User", app.Name),
				v,
			),
		)
	}

	app.Setting = &settingTmp

	return c.Redirect(http.StatusFound, "/setting/user")
}
