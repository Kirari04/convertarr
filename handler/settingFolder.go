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
	"os"

	"github.com/labstack/echo/v4"
)

func GetSettingFolder(c echo.Context) error {
	var v t.SettingFolderValidation

	var folders []m.Folder
	if err := app.DB.Find(&folders).Error; err != nil {
		c.Logger().Error("failed to list folders", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingFolder(
				helper.TCtxWError(c, errors.New("failed to list folders")),
				fmt.Sprintf("%s - Setting Folder", app.Name),
				v,
				folders,
			),
		)
	}

	return helper.Render(c,
		http.StatusOK,
		views.SettingFolder(
			helper.TCtx(c),
			fmt.Sprintf("%s - Setting Folder", app.Name),
			v,
			folders,
		),
	)
}

func PostSettingFolder(c echo.Context) error {
	var v t.SettingFolderValidation

	var folders []m.Folder
	if err := app.DB.Find(&folders).Error; err != nil {
		c.Logger().Error("failed to list folders", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingFolder(
				helper.TCtxWError(c, errors.New("internal server error")),
				fmt.Sprintf("%s - Setting Folder", app.Name),
				v,
				folders,
			),
		)
	}

	if err := c.Bind(&v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingFolder(
				helper.TCtxWError(c, errors.New("bad request")),
				fmt.Sprintf("%s - Setting Folder", app.Name),
				v,
				folders,
			),
		)
	}
	if err := app.Validate.Struct(v); err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingFolder(
				helper.TCtxWError(c, err),
				fmt.Sprintf("%s - Setting Folder", app.Name),
				v,
				folders,
			),
		)
	}
	folderStat, err := os.Stat(v.Folder)
	if err != nil {
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingFolder(
				helper.TCtxWError(c, err),
				fmt.Sprintf("%s - Setting Folder", app.Name),
				v,
				folders,
			),
		)
	}
	if !folderStat.IsDir() {
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingFolder(
				helper.TCtxWError(c, errors.New("defined path is not a directory")),
				fmt.Sprintf("%s - Setting Folder", app.Name),
				v,
				folders,
			),
		)
	}

	if err := app.DB.Create(&m.Folder{
		Path: v.Folder,
	}).Error; err != nil {
		if err := app.DB.Find(&folders).Error; err != nil {
			c.Logger().Error("failed to create folders", err)
			return helper.Render(c,
				http.StatusBadRequest,
				views.SettingFolder(
					helper.TCtxWError(c, errors.New("internal server error")),
					fmt.Sprintf("%s - Setting Folder", app.Name),
					v,
					folders,
				),
			)
		}
	}

	folders = []m.Folder{}
	if err := app.DB.Find(&folders).Error; err != nil {
		c.Logger().Error("failed to list folders", err)
		return helper.Render(c,
			http.StatusBadRequest,
			views.SettingFolder(
				helper.TCtxWError(c, errors.New("internal server error")),
				fmt.Sprintf("%s - Setting Folder", app.Name),
				v,
				folders,
			),
		)
	}

	return c.Redirect(http.StatusFound, "/setting/folder")
}

func DeleteSettingFolder(c echo.Context) error {
	var v t.SettingFolderDeleteValidation

	if err := c.Bind(&v); err != nil {
		c.Logger().Warn("folder delete bind error", err)
		return c.Redirect(http.StatusFound, "/setting/folder")
	}
	if err := app.Validate.Struct(v); err != nil {
		c.Logger().Warn("folder delete validation error", err)
		return c.Redirect(http.StatusFound, "/setting/folder")
	}

	if err := app.DB.Delete(&m.Folder{}, v.FolderId).Error; err != nil {
		c.Logger().Error("failed to delete folder", err)
		return c.Redirect(http.StatusFound, "/setting/folder")
	}

	return c.Redirect(http.StatusFound, "/setting/folder")
}
