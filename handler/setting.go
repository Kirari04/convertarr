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

	if app.Setting.EnableHevcEncoding {
		v.EnableHevcEncoding = &on
	} else {
		v.EnableHevcEncoding = nil
	}
	if app.Setting.EnableNvidiaGpuEncoding {
		v.EnableNvidiaGpuEncoding = &on
	} else {
		v.EnableNvidiaGpuEncoding = nil
	}
	if app.Setting.EnableAmdGpuEncoding {
		v.EnableAmdGpuEncoding = &on
	} else {
		v.EnableAmdGpuEncoding = nil
	}

	if app.Setting.EncodingCrf <= 0 {
		v.EncodingCrf = 25
	} else {
		v.EncodingCrf = app.Setting.EncodingCrf
	}

	if app.Setting.EncodingResolution <= 100 {
		v.EncodingResolution = 1920
	} else {
		v.EncodingResolution = app.Setting.EncodingResolution
	}

	v.EncodingThreads = app.Setting.EncodingThreads
	v.EncodingMaxRetry = app.Setting.EncodingMaxRetry

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

	settingTmp := *app.Setting

	if v.EnableAuthentication != nil && *v.EnableAuthentication == "on" {
		// check if user had been already created
		if settingTmp.Username == "" {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setting(
					helper.TCtxWError(c, errors.New("before you enable authentication you have to create an user")),
					fmt.Sprintf("%s - Setting", app.Name),
					v,
				),
			)
		}

		settingTmp.EnableAuthentication = true
	} else {
		settingTmp.EnableAuthentication = false
	}

	if settingTmp.EnableAuthentication {
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

		settingTmp.AuthenticationType = v.AuthenticationType
	}

	if v.EnableAutomaticScanns != nil && *v.EnableAutomaticScanns == "on" {
		settingTmp.EnableAutomaticScanns = true
	} else {
		settingTmp.EnableAutomaticScanns = false
	}

	if v.AutomaticScannsAtStartup != nil && *v.AutomaticScannsAtStartup == "on" {
		settingTmp.AutomaticScannsAtStartup = true
	} else {
		settingTmp.AutomaticScannsAtStartup = false
	}

	settingTmp.AutomaticScannsInterval = time.Duration(v.AutomaticScannsInterval) * time.Minute

	if v.EnableEncoding != nil && *v.EnableEncoding == "on" {
		settingTmp.EnableEncoding = true
	} else {
		settingTmp.EnableEncoding = false
	}

	if v.EnableHevcEncoding != nil && *v.EnableHevcEncoding == "on" {
		settingTmp.EnableHevcEncoding = true
	} else {
		settingTmp.EnableHevcEncoding = false
	}

	if settingTmp.EnableEncoding {
		settingTmp.EncodingCrf = v.EncodingCrf
		if v.EncodingResolution%2 != 0 {
			return helper.Render(c,
				http.StatusBadRequest,
				views.Setting(
					helper.TCtxWError(c, errors.New("EncodingResolution is not even")),
					fmt.Sprintf("%s - Setting", app.Name),
					v,
				),
			)
		}
		settingTmp.EncodingResolution = v.EncodingResolution
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
		settingTmp.EncodingThreads = v.EncodingThreads
		settingTmp.EncodingMaxRetry = v.EncodingMaxRetry
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

	setting.Value = settingTmp

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

	app.Setting = &settingTmp

	return c.Redirect(http.StatusFound, "/setting")
}
