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
	"runtime"
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

	if app.Setting.EnableEncoding {
		v.EnableEncoding = &on
	} else {
		v.EnableEncoding = nil
	}

	if app.Setting.EncodingCrf <= 0 {
		v.EncodingCrf = 28
	} else {
		v.EncodingCrf = app.Setting.EncodingCrf
	}

	v.EncodingThreads = app.Setting.EncodingThreads

	return helper.Render(c,
		http.StatusOK,
		views.Setting(
			helper.TCtx(c),
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
				helper.TCtxWError(c, errors.New("bad request")),
				fmt.Sprintf("%s - Setting", app.Name),
				v,
			),
		)
	}
	if err := app.Validate.Struct(v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setting(
				helper.TCtxWError(c, err),
				fmt.Sprintf("%s - Setting", app.Name),
				v,
			),
		)
	}

	if v.EnableAuthentication != nil && *v.EnableAuthentication == "on" {
		// check if user had been already created
		if app.Setting.Username == "" {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setting(
					helper.TCtxWError(c, errors.New("before you enable authentication you have to create an user")),
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
					helper.TCtxWError(c, errors.New("AuthenticationType cant be empty when authentication is enabled")),
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

	if v.EnableEncoding != nil && *v.EnableEncoding == "on" {
		app.Setting.EnableEncoding = true
	} else {
		app.Setting.EnableEncoding = false
	}

	if app.Setting.EnableEncoding {
		app.Setting.EncodingCrf = v.EncodingCrf
		if v.EncodingThreads > runtime.NumCPU() {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setting(
					helper.TCtxWError(c, errors.New("you can use more threads than available")),
					fmt.Sprintf("%s - Setting", app.Name),
					v,
				),
			)
		}
		app.Setting.EncodingThreads = v.EncodingThreads
	}

	var setting m.Setting
	if err := app.DB.First(&setting).Error; err != nil {
		c.Echo().Logger.Error("Failed to get setting", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.Setting(
				helper.TCtxWError(c, errors.New("internal server error")),
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
				helper.TCtxWError(c, errors.New("internal server error")),
				fmt.Sprintf("%s - Setting", app.Name),
				v,
			),
		)
	}

	return c.Redirect(http.StatusFound, "/setting")
}
