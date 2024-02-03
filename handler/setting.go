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
	"time"

	"github.com/labstack/echo/v4"
)

func GetSetting(c echo.Context) error {
	var v t.SettingValidator
	on := "on"

	if app.Setting.EnableAuthentication {
		v.EnableAuthentication = &on
	} else {
		v.EnableAuthentication = nil
	}

	v.AuthenticationType = app.Setting.AuthenticationType

	if app.Setting.EnableAutomaticScanns {
		v.EnableAutomaticScanns = &on
	} else {
		v.EnableAutomaticScanns = nil
	}

	if app.Setting.AutomaticScannsAtStartup {
		v.AutomaticScannsAtStartup = &on
	} else {
		v.AutomaticScannsAtStartup = nil
	}

	v.AutomaticScannsInterval = int(app.Setting.AutomaticScannsInterval.Minutes())

	return helper.Render(c,
		http.StatusOK,
		views.Setting(
			nil,
			fmt.Sprintf("%s - Settings", app.Name),
			v,
		),
	)
}

func PostSetting(c echo.Context) error {
	var v t.SettingValidator
	if err := c.Bind(&v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setting(
				errors.New("bad request"),
				fmt.Sprintf("%s - Setting", app.Name),
				v,
			),
		)
	}
	if err := app.Validate.Struct(v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setting(
				err,
				fmt.Sprintf("%s - Setting", app.Name),
				v,
			),
		)
	}

	if v.EnableAuthentication != nil && *v.EnableAuthentication == "on" {
		// check if user had been already created
		var existingUsers int64
		if err := app.DB.Model(&m.User{}).Count(&existingUsers).Error; err != nil {
			c.Echo().Logger.Error("Failed to count existing users", err)
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setting(
					errors.New("internal server error"),
					fmt.Sprintf("%s - Setting", app.Name),
					v,
				),
			)
		}

		if existingUsers == 0 {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setting(
					errors.New("before you enable authentication you have to create an user"),
					fmt.Sprintf("%s - Setting", app.Name),
					v,
				),
			)
		}

		app.Setting.EnableAuthentication = true
	} else {
		app.Setting.EnableAuthentication = false
	}

	if app.Setting.EnableAuthentication {
		if v.AuthenticationType == nil {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setting(
					errors.New("AuthenticationType cant be empty when authentication is enabled"),
					fmt.Sprintf("%s - Setting", app.Name),
					v,
				),
			)
		}

		app.Setting.AuthenticationType = v.AuthenticationType
	}

	if v.EnableAutomaticScanns != nil && *v.EnableAutomaticScanns == "on" {
		app.Setting.EnableAutomaticScanns = true
	} else {
		app.Setting.EnableAutomaticScanns = false
	}

	if v.AutomaticScannsAtStartup != nil && *v.AutomaticScannsAtStartup == "on" {
		app.Setting.AutomaticScannsAtStartup = true
	} else {
		app.Setting.AutomaticScannsAtStartup = false
	}

	app.Setting.AutomaticScannsInterval = time.Duration(v.AutomaticScannsInterval) * time.Minute

	var setting m.Setting
	if err := app.DB.First(&setting).Error; err != nil {
		c.Echo().Logger.Error("Failed to get setting", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setting(
				errors.New("internal server error"),
				fmt.Sprintf("%s - Setting", app.Name),
				v,
			),
		)
	}

	setting.Value = *app.Setting

	if err := app.DB.Save(&setting).Error; err != nil {
		c.Echo().Logger.Error("Failed to update setting", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setting(
				errors.New("internal server error"),
				fmt.Sprintf("%s - Setting", app.Name),
				v,
			),
		)
	}

	return c.Redirect(http.StatusFound, "/setting")
}
